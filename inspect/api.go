// 调用API巡检包
package inspect

import (
	"fmt"
	"strings"
)

var (
	apiCallTag  string = "user=%s_apicall_%s"
	apiCallFlag string = "api_name=%s,api_result=%d,service_name=%s,msg=%s"
)

// APITriggerReportOk 上报API成功调用API
func APITriggerReportOk(apiName string) {
	tag := fmt.Sprintf(apiCallTag, serviceName, apiName)
	field := fmt.Sprintf(apiCallFlag, apiName, ok, serviceName, "ok")
	triggerReport(tag, field)
}

// APITriggerReportFailed 上报API失败调用
func APITriggerReportFailed(apiName string, msg string) {
	tag := fmt.Sprintf(apiCallTag, serviceName, apiName)
	field := fmt.Sprintf(apiCallFlag, apiName, failed, serviceName, strings.ReplaceAll(msg, ",", separate))
	triggerReport(tag, field)
}

// APITriggerReportWithError if err is nil, report ok
func APITriggerReportWithError(apiName string, err error) {
	if nil == err {
		APITriggerReportOk(apiName)
		return
	}
	APITriggerReportFailed(apiName, err.Error())
}
