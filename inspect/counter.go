package inspect

import (
	"fmt"
)

var (
	counterTag  string = "user=%s_counter_%s,counter_name=%s"
	counterFlag string = "counter_ename=%s,count=%d,service_name=%s"
)

// CounterReport 计数器指标上报
func CounterReport(counterEname, counterName string, count int64) {
	tag := fmt.Sprintf(counterTag, serviceName, counterEname, counterName)
	field := fmt.Sprintf(counterFlag, counterEname, count, serviceName)
	triggerReport(tag, field)
}
