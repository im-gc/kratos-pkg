package inspect

import (
	"fmt"
	"strings"
)

var (
	taskTag  string = "user=%s_task_%s,task_name=%s"
	taskFlag string = "task_ename=%s,task_result=%d,service_name=%s,msg=%s"
)

// TaskTriggerReportOk 上报业务成功
func TaskTriggerReportOk(taskEname, taskName string) {
	tag := fmt.Sprintf(taskTag, serviceName, taskEname, taskName)
	field := fmt.Sprintf(taskFlag, taskEname, ok, serviceName, "success")
	triggerReport(tag, field)
}

// TaskTriggerReportFailed 上报业务失败
func TaskTriggerReportFailed(taskEname, taskName string, msg string) {
	tag := fmt.Sprintf(taskTag, serviceName, taskEname, taskName)
	field := fmt.Sprintf(taskFlag, taskEname, failed, serviceName, strings.ReplaceAll(msg, ",", separate))
	triggerReport(tag, field)
}

// TaskTriggerReportWithError if err is nil, report ok
func TaskTriggerReportWithError(taskEname, taskName string, err error) {
	if nil == err {
		TaskTriggerReportOk(taskEname, taskName)
		return
	}
	TaskTriggerReportFailed(taskEname, taskName, err.Error())
}
