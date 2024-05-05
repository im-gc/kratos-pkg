package inspect

import (
	"fmt"
)

var (
	taskQueueTag  string = "user=%s_taskqueue_%s,taskqueue_name=%s"
	taskQueueFlag string = "taskqueue_ename=%s,taskqueue_length=%d,service_name=%s"
)

// TaskQueueLengthReport 任务队列长度上报
func TaskQueueLengthReport(taskQueueEname, taskQueueName string, length int64) {
	tag := fmt.Sprintf(taskQueueTag, serviceName, taskQueueEname, taskQueueName)
	field := fmt.Sprintf(taskQueueFlag, taskQueueEname, length, serviceName)
	triggerReport(tag, field)
}
