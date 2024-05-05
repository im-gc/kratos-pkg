package inspect

import (
	"fmt"
	"strings"
)

var (
	businessTag  string = "user=%s_business_%s"
	businessFlag string = "business_name=%s,business_result=%d,service_name=%s,msg=%s"
)

// BusinessTriggerReportOk 上报业务成功
func BusinessTriggerReportOk(businessName string) {
	tag := fmt.Sprintf(businessTag, serviceName, businessName)
	field := fmt.Sprintf(businessFlag, businessName, ok, serviceName, "ok")
	triggerReport(tag, field)
}

// BusinessTriggerReportFailed 上报业务失败
func BusinessTriggerReportFailed(businessName string, msg string) {
	tag := fmt.Sprintf(businessTag, serviceName, businessName)
	field := fmt.Sprintf(businessFlag, businessName, failed, serviceName, strings.ReplaceAll(msg, ",", separate))
	triggerReport(tag, field)
}

// BusinessTriggerReportWithError if err is nil, report ok
func BusinessTriggerReportWithError(businessName string, err error) {
	if nil == err {
		BusinessTriggerReportOk(businessName)
		return
	}
	BusinessTriggerReportFailed(businessName, err.Error())
}
