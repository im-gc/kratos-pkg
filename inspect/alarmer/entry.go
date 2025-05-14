// Author houguofa
// Copyright @2018 houguofa. All Rights Reserved.

package alarmer

import (
	"errors"
	"time"

	center "github.com/im-gc/kratos-pkg/inspect/alarmer/workCenter"
)

func Init(cMetric string) error {
	return CustomMetricFlag(cMetric)
}

func CustomMetricFlag(cMetric string) error {
	if len(cMetric) <= 0 {
		return nil
	}

	center.SetDefaultMetric(cMetric)
	return nil
}

func FixAlarmReportApi(url string) error {
	center.SetDefaultReportApi(url)
	return nil
}

func TriggerAlarmSpecifiedTime(endpoint, tags, fields string, timestamp int64) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMonitorInfo([]string{endpoint}, []string{tags}, []string{fields}, timestamp)
	return nil
}

func TriggerAlarm(endpoint, tags, fields string) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMonitorInfo([]string{endpoint}, []string{tags}, []string{fields}, time.Now().Local().Unix())
	return nil
}

func TriggerAlarmSetSpecifiedTime(endpoint, tags, fields []string, timestamp int64) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMonitorInfo(endpoint, tags, fields, timestamp)
	return nil
}

func TriggerAlarmSet(endpoint, tags, fields []string) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMonitorInfo(endpoint, tags, fields, time.Now().Local().Unix())
	return nil
}

func TriggerMetricAlarmSpecifiedTime(metric string, endpoint, tags, fields string, timestamp int64) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMetricMonitorInfo([]string{metric}, []string{endpoint}, []string{tags}, []string{fields}, timestamp)
	return nil
}

func TriggerMetricAlarm(metric string, endpoint, tags, fields string) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMetricMonitorInfo([]string{metric}, []string{endpoint}, []string{tags}, []string{fields}, time.Now().Local().Unix())
	return nil
}

func TriggerMetricAlarmSetSpecifiedTime(metric, endpoint, tags, fields []string, timestamp int64) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMetricMonitorInfo(metric, endpoint, tags, fields, timestamp)
	return nil
}

func TriggerMetricAlarmSet(metric, endpoint, tags, fields []string) error {
	if len(tags) == 0 || len(fields) == 0 {
		return errors.New("args error")
	}

	go center.DisposeMetricMonitorInfo(metric, endpoint, tags, fields, time.Now().Local().Unix())
	return nil
}
