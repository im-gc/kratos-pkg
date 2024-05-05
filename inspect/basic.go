package inspect

import (
	"fmt"
	"runtime"
	"time"
)

const (

	// restart
	restartTag  string = "user=%s_restart"
	restartFlag string = "process_restart_flag=%d,service_name=%s"

	//程序心跳
	heartBeatTag  string = "user=%s_heartbeat"
	heartBeatFlag string = "process_heartbeat_flag=%d,service_name=%s"

	//goroutine
	goroutineNumTag  string = "user=%s_goroutine"
	goroutineNumFlag string = "gorountine_num_flag=%d,service_name=%s"

	// mem use
	memSysAllocTag    string = "user=%s_mem_sys_alloc" //系统使用内存
	memSysAllocFlag   string = "mem_sys_alloc_flag=%d,service_name=%s"
	memHeapAllocTag   string = "user=%s_mem_heap_alloc" //堆分配内存
	memHeapAllocFlag  string = "mem_heap_alloc_flag=%d,service_name=%s"
	memHeapInuseTag   string = "user=%s_mem_heap_inuse" //堆分配使用内存
	memHeapInuseFlag  string = "mem_heap_inuse_flag=%d,service_name=%s"
	memHeapIdleTag    string = "user=%s_mem_heap_idle" //堆分配空闲内存
	memHeapIdleFlag   string = "mem_heap_idle_flag=%d,service_name=%s"
	memStackAllocTag  string = "user=%s_mem_stack_alloc" //栈分配内存
	memStackAllocFlag string = "mem_stack_alloc_flag=%d,service_name=%s"
	memStackInuseTag  string = "user=%s_mem_stack_inuse" //栈使用内存
	memStackInuseFlag string = "mem_stack_inuse_flag=%d,service_name=%s"

	// restart
	reloadCfgTag  string = "user=%s_cfg_reload"
	reloadCfgFlag string = "reload_result=%d,service_name=%s,msg=%s"
	reloadOK             = 1
	reloadFaild          = -1
)

func initBasic() {
	reportProcessRestart()
	reportHeartbeat()
	reportGoroutineNum()
	reportMemUse()
}

// reportProcessRestart process restart
func reportProcessRestart() {
	go func() {
		triggerReport(fmt.Sprintf(restartTag, serviceName), fmt.Sprintf(restartFlag, 1, serviceName))
		time.Sleep(time.Second * 60)
		triggerReport(fmt.Sprintf(restartTag, serviceName), fmt.Sprintf(restartFlag, 0, serviceName))
	}()
}

// reportHeartbeat heartbeat
func reportHeartbeat() {
	AsyncTicker(60, func() {
		triggerReport(fmt.Sprintf(heartBeatTag, serviceName), fmt.Sprintf(heartBeatFlag, 1, serviceName))
	})
}

func reportMemUse() {
	AsyncTicker(60, func() {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		triggerReport(fmt.Sprintf(memSysAllocTag, serviceName), fmt.Sprintf(memSysAllocFlag, ms.Sys/1000/1000, serviceName))
		triggerReport(fmt.Sprintf(memHeapAllocTag, serviceName), fmt.Sprintf(memHeapAllocFlag, ms.Alloc/1000/1000, serviceName))
		triggerReport(fmt.Sprintf(memHeapInuseTag, serviceName), fmt.Sprintf(memHeapInuseFlag, ms.HeapInuse/1000/1000, serviceName))
		triggerReport(fmt.Sprintf(memHeapIdleTag, serviceName), fmt.Sprintf(memHeapIdleFlag, ms.HeapIdle/1000/1000, serviceName))
		triggerReport(fmt.Sprintf(memStackAllocTag, serviceName), fmt.Sprintf(memStackAllocFlag, ms.StackSys/1000/1000, serviceName))
		triggerReport(fmt.Sprintf(memStackInuseTag, serviceName), fmt.Sprintf(memStackInuseFlag, ms.StackInuse/1000/1000, serviceName))
	})
}

// report goroutine num
func reportGoroutineNum() {
	AsyncTicker(60, func() {
		triggerReport(fmt.Sprintf(goroutineNumTag, serviceName), fmt.Sprintf(goroutineNumFlag, runtime.NumGoroutine(), serviceName))
	})
}

func ReloadCfgTriggerReportWithError(err error) {

	var tag, field string

	tag = fmt.Sprintf(reloadCfgTag, serviceName)
	if field = fmt.Sprintf(reloadCfgFlag, reloadOK, serviceName, "OK"); nil != err {
		field = fmt.Sprintf(reloadCfgFlag, reloadFaild, serviceName, err.Error())
	}
	triggerReport(tag, field)
}
