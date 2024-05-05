# xxl-job

基于 XXL-Job 封装的 go-sdk, 提供了 proto 定义，建议直接引入

在 kratos 配置文件 conf.proto 中引入 `api/job/job.proto` 并写入到 Bootstrap 结构体中
重新生成配置文件，即可加载 xxl-job 的配置文件内容

支持传递 go-kratos 的 http.Server 为其提供任务注册端口复用