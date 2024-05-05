# job 公共基类

## 提供以下功能

- 支持同步执行任务、异步执行任务
- 收集任务执行结果，推送至监控
  
## 接入指南

```
在main.go中初始化监控包

import "github.com/imkouga/kratos-pkg/inspect"

if err := inspect.Init1("组件名", "推送至监控的metric"); nil != err {
    // do something
}
```

```
// job示例

import "github.com/imkouga/kratos-pkg/contrib/job/jobset"

NodeJob struct {
    *BaseSet
}

func NewNodeJob()(*NodeJob,error) {

    js,err := jobset.NewJobSet()
    if nil != err {
        return nil, err
    }

    return &NodeJob{
        js,
    }, nil
}

func (j *NodeJob) CalNodeMinPlanLine(ctx context.Context, payload string) (string, error) {

    // 异步执行
	j.AsyncDo(ctx, "task_ename", "任务中文名", func() (string, error) {
        return "", nil
    }

    // 同步执行
    j.SyncDo(ctx, "task_ename", "任务中文名", func() (string, error) {
        return "", nil
    } 
}
```