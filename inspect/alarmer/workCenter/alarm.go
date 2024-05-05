package workCenter

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func DisposeMonitorInfo(endpoint, tags, fileds []string, timestamp int64) {
	ri := buildReportInfo(endpoint, tags, fileds, timestamp)
	disposeMonitorInfo(ri)
}

func DisposeMetricMonitorInfo(metric, endpoint, tags, fileds []string, timestamp int64) {
	headers := buildSingleReportHeaders(metric, endpoint, tags, fileds, timestamp)
	r := &ReportInfo{Data: headers}
	disposeMonitorInfo(r)
}

func disposeMonitorInfo(ri *ReportInfo) {

	if len(ri.Data) <= 0 {
		return
	}

	var ai AlarmInfos
	ai = make([]*AlarmHeader, 0, len(ri.Data))
	for _, rh := range ri.Data {
		ah := newAlarmHeader(rh.EndPoint, rh.Metric, rh.Tags, rh.Fields, rh.Timestamp)
		ai = append(ai, ah)
	}

	reportAlarm(ai)
}

func reportAlarm(ai AlarmInfos) {

	data, err := encode2Alarm(&ai)
	if nil != err {
		fmt.Println(err)
		return
	}

	report(reportAlarmDataUrl, data)
}

func report(url string, data []byte) {

	var (
		body   *bytes.Buffer
		req    *http.Request
		resp   *http.Response
		client *http.Client
		err    error
	)

	body = bytes.NewBuffer(data)
	if req, err = http.NewRequest(http.MethodPost, url, body); nil != err {
		fmt.Println("major issue : local workCenter api request failed ...", err)
		return
	}

	client = &http.Client{
		Timeout: time.Second * 60,
	}
	if resp, err = client.Do(req); nil != err {
		fmt.Println("major issue : local workCenter api response failed ...", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("major issue : local workCenter api response failed ...", err)
		return
	}
}
