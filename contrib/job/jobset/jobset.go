package job

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/imkouga/kratos-pkg/inspect"
)

type Option func(j *JobSet)

type JobSet struct{}

func NewJobSet(opts ...Option) (*JobSet, error) {
	return nil, nil
}

// AsyncDo 异步执行
func (j *JobSet) AsyncDo(ctx context.Context, jobEName, jobName string, fn func() (string, error)) (string, error) {
	go j.do(ctx, jobEName, jobName, fn)
	return "", nil
}

// SyncDo 同步执行
func (j *JobSet) SyncDo(ctx context.Context, jobEName, jobName string, fn func() (string, error)) (string, error) {
	return j.do(ctx, jobEName, jobName, fn)
}

func (j *JobSet) do(ctx context.Context, jobEName, jobName string, fn func() (string, error)) (string, error) {

	var (
		result string
		err    error
	)

	defer func() {
		if e := recover(); nil != e {
			err = fmt.Errorf("%+v", e)
		}
		inspect.TaskTriggerReportWithError(jobEName, jobName, err)
		if nil != err {
			log.Infof("系统异常，任务名：%s，错误：%s", jobName, err.Error())
		}
	}()

	result, err = fn()
	return result, err
}
