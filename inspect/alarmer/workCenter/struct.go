// Author houguofa
// Copyright @2018 houguofa. All Rights Reserved.

package workCenter

import (
	"encoding/json"
	"strings"
	"time"
)

// [{"counterType": "GAUGE","step": 60,"value":0,"metric": "gslb_center_data","endpoint": "hostname", "tags":"user=ping",
// "fields":"size=","timestamp":1473327905}]
const (
	COUNTER_TYPE  = "GAUGE"
	HOSTNAME_FILE = "/allconf/hostname.conf"
)

var (
	localEndPoint      string
	metricFlag         = "default_alarm_metric"
	reportAlarmDataUrl = "http://127.0.0.1:10699/v1/push"
)

type AlarmHeader struct {
	CounterType string `json:"counterType"`
	Step        int    `json:"step"`
	Value       int    `json:"value"`
	Metric      string `json:"metric"`
	EndPoint    string `json:"endpoint"`
	Tags        string `json:"tags"`
	Fields      string `json:"fields"`
	Timestamp   int64  `json:"timestamp"` // 当前秒数
}

type AlarmInfos []*AlarmHeader

type ReportHeader struct {
	Metric    string `json:"metric"`
	Tags      string `json:"tags"`
	Fields    string `json:"fields"`
	EndPoint  string `json:"endpoint"`
	Timestamp int64  `json:"timestamp"`
}

func buildSingleReportStructer(endpoint, tags, fields string, timestamp int64) *ReportHeader {
	if len(endpoint) <= 0 {
		endpoint = localEndPoint
	}

	return &ReportHeader{
		Metric:    metricFlag,
		EndPoint:  endpoint,
		Tags:      tags,
		Fields:    fields,
		Timestamp: timestamp,
	}
}

func buildReportStructers(endpoint, tags, fields []string, timestamp int64) []*ReportHeader {

	sts := make([]*ReportHeader, 0, len(tags))

	for i := 0; i < len(tags); i++ {
		st := buildSingleReportStructer(endpoint[i], tags[i], fields[i], timestamp)
		sts = append(sts, st)
	}

	return sts
}

func buildReportInfo(endpoint, tags, fields []string, timestamp int64) *ReportInfo {
	info := buildReportStructers(endpoint, tags, fields, timestamp)
	return &ReportInfo{Data: info}
}

type ReportInfo struct {
	Data []*ReportHeader `json:"data"`
}

func getLocalTimeBySecond() int64 {
	return int64(time.Now().Unix()/60) * 60
}

func newAlarmHeader(endpoint, metric, tags, fields string, timestamp int64) *AlarmHeader {

	if timestamp == 0 {
		timestamp = getLocalTimeBySecond()
	}

	ai := &AlarmHeader{CounterType: COUNTER_TYPE,
		Step:      60,
		Metric:    metric,
		EndPoint:  endpoint,
		Tags:      tags,
		Fields:    fields,
		Timestamp: timestamp}

	return ai
}

func DecodeReport(data []byte) (*ReportInfo, error) {
	ri := new(ReportInfo)
	err := json.Unmarshal(data, ri)
	return ri, err
}

func encode2Alarm(info *AlarmInfos) ([]byte, error) {
	return json.Marshal(info)
}

func SetDefaultMetric(metric string) {

	if len(strings.TrimSpace(metric)) <= 0 {
		return
	}

	metricFlag = metric
}

func SetDefaultReportApi(url string) {
	if len(strings.TrimSpace(url)) <= 0 {
		return
	}

	reportAlarmDataUrl = url
}

func buildSingleReportHeader(metric, endpoint, tags, fields string, timestamp int64) *ReportHeader {
	if len(endpoint) <= 0 {
		endpoint = localEndPoint
	}

	return &ReportHeader{
		Metric:    metric,
		EndPoint:  endpoint,
		Tags:      tags,
		Fields:    fields,
		Timestamp: timestamp,
	}
}

func buildSingleReportHeaders(metric, endpoint, tags, fields []string, timestamp int64) []*ReportHeader {

	sts := make([]*ReportHeader, 0, len(tags))

	for i := 0; i < len(tags); i++ {
		st := buildSingleReportHeader(metric[i], endpoint[i], tags[i], fields[i], timestamp)
		sts = append(sts, st)
	}

	return sts
}
