package kv

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type (
	OptionFn[K comparable, V any] func(c *simpleKV[K, V])
	simpleKV[K comparable, V any] struct {
		ename string
		name  string

		sync.RWMutex
		dataCheck DataCheck[K, V]
		data      map[K]V
		dataKey   map[K]int64

		metricReport MetricReport

		expired int64 // key 过期时间

		interval4Key int            // 0表示不刷新缓存
		fetchKey     FetchKey[K, V] // 按key数据查询器

		interval4Global int               // 0表示不刷新缓存
		fetchGlobal     FetchGlobal[K, V] // 全局数据查询器
	}
)

func NewKV[K comparable, V any](opts ...OptionFn[K, V]) (KV[K, V], error) {

	kv := &simpleKV[K, V]{
		data:            make(map[K]V, 100),
		dataKey:         make(map[K]int64, 100),
		interval4Key:    0,
		interval4Global: 0,
		expired:         NeverExpired,
	}

	for _, opt := range opts {
		opt(kv)
	}
	if nil == kv.fetchKey && nil == kv.fetchGlobal {
		return nil, ErrDataQueryerNotProvider
	}

	ticker(kv.interval4Global, kv.refresh4Global, nil)
	ticker(kv.interval4Key, kv.refresh4Key, nil)
	ticker(60, kv.expiredHandle, nil)
	ticker(30*60, kv.reset, nil)

	return kv, nil

}

func (kv *simpleKV[K, V]) reset() {

	kv.Lock()
	defer kv.Unlock()

	data := make(map[K]V, len(kv.data))
	dataKey := make(map[K]int64, len(kv.dataKey))

	for k, v := range kv.data {
		data[k] = v
	}
	for k, v := range kv.dataKey {
		dataKey[k] = v
	}

	kv.data = data
	kv.dataKey = dataKey
}

func (kv *simpleKV[K, V]) flush4Global(datas map[K]V) {

	kv.Lock()
	defer kv.Unlock()

	for k, v := range datas {
		kv.data[k] = v
		kv.flushKey(k)
	}
}

func (kv *simpleKV[K, V]) flush4Key(k K, v V) {

	kv.Lock()
	defer kv.Unlock()

	kv.data[k] = v
	kv.flushKey(k)
}

func (kv *simpleKV[K, V]) flushKey(k K) {

	if kv.expired == NeverExpired {
		return
	}
	kv.dataKey[k] = time.Now().Local().Unix() + kv.expired
}

func (kv *simpleKV[K, V]) refresh4Global() {

	if nil == kv.fetchGlobal {
		return
	}
	data, err := kv.fetchGlobal()
	if nil != err {
		log.Errorf("KV[%s]系统错误，全局刷新缓存失败[%s]", kv.name, err.Error())
		return
	}
	if nil == kv.doDataCheck(data) {
		kv.flush4Global(data)
	}
}

func (kv *simpleKV[K, V]) refresh4Key() {

	if nil == kv.fetchKey {
		return
	}
	keys := kv.keys()
	for _, k := range keys {
		data, err := kv.fetchKey(k)
		if nil != err {
			log.Errorf("KV[%s]系统错误，按key[%+v]刷新缓存失败[%s]", kv.name, k, err.Error())
			continue
		}
		if nil == kv.doDataCheck(data) {
			for k, v := range data {
				kv.flush4Key(k, v)
			}
		}
	}
}

func (kv *simpleKV[K, V]) doMetricReport(err error) {
	if nil != kv.metricReport {
		kv.metricReport(kv.ename, kv.name, err)
	}
}

func (kv *simpleKV[K, V]) doDataCheck(datas map[K]V) KVError {

	var err error

	defer func() {
		kv.doMetricReport(err)
	}()

	if nil != kv.dataCheck {
		if err = kv.dataCheck(datas); nil != err {
			log.Errorf("KV[%s]系统错误，数据检测失败。失败原因：%s", kv.name, err.Error())
			return NewErrDefault(kv.ename, err.Error())
		}
	}
	return nil
}

func (kv *simpleKV[K, V]) first(key K) (V, KVError) {

	kv.RLock()
	defer kv.RUnlock()

	v, e := kv.data[key]
	if !e {
		return v, NewErrKeyNotFound(kv.name, key)
	}
	return v, nil
}

func (kv *simpleKV[K, V]) all() (map[K]V, KVError) {

	kv.RLock()
	defer kv.RUnlock()

	clone := make(map[K]V, len(kv.data))
	for k, v := range kv.data {
		clone[k] = v
	}
	return clone, nil
}

func (kv *simpleKV[K, V]) keys() []K {

	kv.RLock()
	defer kv.RUnlock()

	var keys []K
	for k := range kv.data {
		keys = append(keys, k)
	}
	return keys
}

func (kv *simpleKV[K, V]) expiredHandle() {

	keys := kv.scanKey()
	if len(keys) <= 0 {
		return
	}
	kv.removeExpireKeys(keys)
}

// scanKey 扫描出已过期的key
func (kv *simpleKV[K, V]) scanKey() []K {

	kv.RLock()
	defer kv.RUnlock()

	var (
		expiredKeys []K
	)

	now := time.Now().Local().Unix()
	for k, v := range kv.dataKey {
		if v < now && v != NeverExpired {
			expiredKeys = append(expiredKeys, k)
		}
	}

	return expiredKeys
}

// scanKey 扫描出已过期的key
func (kv *simpleKV[K, V]) removeExpireKeys(keys []K) {

	kv.Lock()
	defer kv.Unlock()

	for _, key := range keys {
		delete(kv.data, key)
		delete(kv.dataKey, key)
	}
}

// 直接获取缓存
func (kv *simpleKV[K, V]) GetCacheByKey(key K) (V, KVError) {

	var (
		v   V
		err KVError
	)
	if nil == kv {
		return v, ErrNotInit
	}

	if v, err = kv.first(key); nil != err {
		return v, err
	}
	return v, nil
}

// 获取缓存，如果没有命中则回源
func (kv *simpleKV[K, V]) GetCacheOrElseSourceByKey(key K) (V, KVError) {

	var (
		v   V
		err KVError
	)

	if v, err = kv.GetCacheByKey(key); nil == err {
		return v, nil
	}

	if v, err = kv.GetSourceByKey(key); nil == err {
		kv.flush4Key(key, v)
		return v, nil
	}

	if _, ok := err.(*ErrFn4KeyNotProvider); ok {
		if _, err = kv.GetSource(); nil != err {
			return v, NewErrDefault(kv.ename, err.Error())
		}
		return kv.GetCacheByKey(key)
	}
	return v, err
}

// 直接获取缓存
func (kv *simpleKV[K, V]) GetCache() (map[K]V, KVError) {

	if nil == kv {
		return nil, ErrNotInit
	}
	d, e := kv.all()
	if nil != e {
		return nil, NewErrDefault(kv.name, e.Error())
	}
	return d, nil
}

// 获取缓存，如果没有命中则回源
func (kv *simpleKV[K, V]) GetCacheOrElseSource() (map[K]V, KVError) {

	var (
		d    map[K]V
		err  KVError
		err1 KVError
	)

	if d, err = kv.GetCache(); nil == err {
		return d, nil
	}
	if d, err1 = kv.GetSource(); nil == err1 {
		kv.flush4Global(d)
		return d, nil
	}
	return nil, err1
}

// 直接回源
func (kv *simpleKV[K, V]) GetSourceByKey(key K) (V, KVError) {

	var (
		v V
	)
	if nil == kv {
		return v, ErrNotInit
	}
	if nil != kv.fetchKey {
		kvs, err := kv.fetchKey(key)
		if nil != err {
			return v, NewErrDefault(kv.name, err.Error())
		}
		if err := kv.doDataCheck(kvs); nil != err {
			return v, err
		}
		if v, exist := kvs[key]; exist {
			return v, nil
		}
		return v, NewErrKeyNotFound(kv.name, key)
	}

	return v, NewErrFn4KeyNotProvider(kv.name)
}

// 直接回源
func (kv *simpleKV[K, V]) GetSource() (map[K]V, KVError) {

	if nil == kv {
		return nil, ErrNotInit
	}
	if nil != kv.fetchGlobal {
		d, e := kv.fetchGlobal()
		if nil != e {
			return nil, NewErrDefault(kv.name, e.Error())
		}
		if err := kv.doDataCheck(d); nil != err {
			return nil, err
		}
		return d, nil
	}

	return nil, NewErrFn4GlobalNotProvider(kv.name)
}

// 查询数据源，如果查询失败 则通过本地缓存兜底
// 如果查询源数据失败，则走本地缓存。 error将返回ErrCacheDefense
func (kv *simpleKV[K, V]) GetSourceOrElseCache() (map[K]V, KVError) {

	var (
		d    map[K]V
		err  KVError
		err1 KVError
	)

	if d, err = kv.GetSource(); nil == err {
		return d, nil
	}
	if d, err1 = kv.GetCache(); nil == err1 {
		return d, nil
	}
	return nil, err
}

// 查询数据源，如果查询失败 则通过本地缓存兜底
// 如果查询源数据失败，则走本地缓存。 error将返回ErrCacheDefense
func (kv *simpleKV[K, V]) GetSourceOrElseCacheByKey(key K) (V, KVError) {

	var (
		v   V
		err KVError
	)

	if v, err = kv.GetSourceByKey(key); nil != err {
		cache, e := kv.GetCacheByKey(key)
		if nil != e {
			return cache, err
		}
		return cache, NewErrCacheDefense(kv.name, err.Error())
	}

	kv.flush4Key(key, v)

	return v, nil
}

// WithEname KV英文名
func WithEname[K comparable, V any](ename string) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.ename = ename
	}
}

// WithName KV中文名
func WithName[K comparable, V any](name string) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.name = name
	}
}

// WithExpired key的过期时间 默认表示不过期
func WithExpired[K comparable, V any](expired int64) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.expired = expired
	}
}

// WithInterval4FetchKey 指定周期更新 0表示不更新
func WithInterval4FetchKey[K comparable, V any](interval int) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.interval4Key = interval
	}
}

// WithInterval4FetchGlobal 指定周期更新 0表示不更新
func WithInterval4FetchGlobal[K comparable, V any](interval int) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.interval4Global = interval
	}
}

// WithFetchKeyQuery 按key数据查询器
func WithFetchKeyQuery[K comparable, V any](fn FetchKey[K, V]) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.fetchKey = fn
	}
}

// WithFetchGlobalQuery 全局数据查询器
func WithFetchGlobalQuery[K comparable, V any](fn FetchGlobal[K, V]) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.fetchGlobal = fn
	}
}

// WithDataCheck 数据检查器
func WithDataCheck[K comparable, V any](fn DataCheck[K, V]) OptionFn[K, V] {
	return func(kv *simpleKV[K, V]) {
		kv.dataCheck = fn
	}
}

func WithMetricReport[K comparable, V any](fn MetricReport) OptionFn[K, V] {
	return func(c *simpleKV[K, V]) {
		c.metricReport = fn
	}
}
