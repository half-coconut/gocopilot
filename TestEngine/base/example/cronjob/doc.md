使用 k8s 作为任务调度

除了 cron，之外还有 gocron，其他：

- 基于 k8s 的任务调度
- 配置没有秒，从分钟开始
- 完成配置文件的编写

```shell
GOOS=linux GOARCH=arm go build -o demojob .
docker build -t halfcoconut/hello-egg/demojob-live:v0.0.1 .
docker build -t demojob-live:v0.0.1 .
```

打包编译，生成镜像，然后生成 k8s resource 为CronJob 的 yaml 文件；

- 注意：使用 k8s 完成任务调度，没有亮点可言，缺少解决定制化的需求，
- 比如任务的开始，结束，**任务编排**

```shell
 kubectl apply -f demojob.yaml
 kubectl get jobs --watch

➜  cronjob kubectl get cronjobs              
NAME      SCHEDULE    TIMEZONE   SUSPEND   ACTIVE   LAST SCHEDULE   AGE
demojob   * * * * *   <none>     False     3        16s             2m55s

 
```

事实上，基于 Mysql 的实现是一个简单的分布式任务调度平台。在这个基础上，进一步提供管理功能，做成一个分布式任务调度平台。

- 加入部门管理和权限控制功能（比如执行 1000 分钟，cpu 资源，在 k8s 基础上等，做好资源管理）
- 加入 HTTP 和 GRPC 任务支持(也就是调度一个任务，就是调用一个 HTTP 接口，或者调用一个 GRPC 接口)
- 加入任务执行历史的功能(也就是记录任务的每一次执行情况)
    - 手动触发功能
    - 历史查询功能
    - 监控告警功能

性能优化：缓存方案，SQL 优化
可用性优化：服务治理相关的，熔断，限流(滑动窗口等之外)，降级(返回默认响应)...

微服务网关...



blibli 架构设计
明确问题
1. 内部系统
2. mvp
3. jobs 支持全部类型

功能需求：
1. 用户提交任务，长任务，短任务
2. dashboard
3. a.预定在未来某刻执行，定时任务
b.有向环形图的 job

non functional requirements
1.延迟 10s, dashboard 60s
2.可扩展
3.可靠性

job data model
1 repo(code,config,excutable binary)
e metadata

job 
id
owner_id
binary url
input path
output path
created_time
status
num_of_retry

status: ready->waiting -> running-> success
                  |         |
                retry <-3次 failed ->3次 final failure

