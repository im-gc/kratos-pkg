package inspect

import (
	"errors"
	"fmt"

	"github.com/imkouga/kratos-pkg/inspect/alarmer"
)

const (
	defaultEndpointInfo = ""
	ok                  = 1
	failed              = -1
)

var serviceName string //组件名
var onOff bool         //开关
var ErrInitParamMustBeProvided = errors.New("sorry, metric or serviceName must be provided")

func setServiceName(name string) {
	serviceName = name
}

func setOnOff(value bool) {
	onOff = value
}

func isOnOff() bool {
	return onOff
}

func triggerReport(tags, fields string) {
	if !isOnOff() {
		return
	}
	alarmer.TriggerAlarm(defaultEndpointInfo, tags, fields)
}

func triggerReportSpecifiedTime(tags, fields string, timestamp int64) {
	if !isOnOff() {
		return
	}
	alarmer.TriggerAlarmSpecifiedTime(defaultEndpointInfo, tags, fields, timestamp)
}

func TriggerReportWithError(tags string, err error) {
	tags = fmt.Sprintf("user=%s_%s", serviceName, tags)
	fields := fmt.Sprintf("result=%d,service_name=%s,msg=%s", failed, serviceName, err.Error())
	triggerReport(tags, fields)
}

// TriggerReport is a function of report to monitor agent
func TriggerReport(tags, fields string) {
	tags = fmt.Sprintf("user=%s_%s", serviceName, tags)
	fields = fmt.Sprintf("%s,service_name=%s", fields, serviceName)
	triggerReport(tags, fields)
}

// TriggerReportSpecifiedTime is a function of report to monitor agent
func TriggerReportSpecifiedTime(tags, fields string, timestamp int64) {
	tags = fmt.Sprintf("user=%s_%s", serviceName, tags)
	fields = fmt.Sprintf("%s,service_name=%s", fields, serviceName)
	triggerReportSpecifiedTime(tags, fields, timestamp)
}
