# gocopilot

### Architecture

![coreengine drawio](https://github.com/user-attachments/assets/b12b1ac1-a31a-418b-8095-ef353373f883)

- 微服务架构，新增缓存方案，消息队列，分布式任务调度等模块

### Features

- 接口测试 interface testing
    - 录入接口和参数，复制接口和参数，调试和执行测试
    - 支持 DeepSeek 生成接口，转化为接口测试用例集
- 性能测试任务 performance testing tasks
    - 复用接口测试，批量选择执行的接口
    - 设置 rate limit，并发量，最大并发量，持续时间
    - 生成性能测试报告，包含性能指标：
      ```shell
      +++ Requests +++
      [total 总请求数: 100]
      [rate 请求速率: 9.64]
      [throughput 吞吐量: 9.64]
        
      +++ Duration +++
      [total 总持续时间: 10.378s]
        
      +++ Latencies +++
      [min 最小响应时间: 249.259ms]
      [mean 平均响应时间: 103.779ms]
      [max 最大响应时间: 2.008s]
      [P50 百分之50 响应时间 (中位数): 270.165ms]
      [P90 百分之90 响应时间: 515.427ms]
      [P95 百分之95 响应时间: 942.023ms]
      [P99 百分之99 响应时间: 2.008s]
        
      +++ Success +++
      [ratio 成功率: 100.00%]
      [status codes:  200:100]
      [passed: 100]
      [failed: 0]
      ```
- 任务调度 cron jobs
  - 创建多种任务类型，外部调用http任务，内部调用 tasks 任务调度
  - 基于 MySql 的分布式任务调度来提高任务调度能力、并发执行任务调度
  - 区分长短任务模式，按照超时时间或 cron 表达式执行
  - 支持手动触发任务，允许用户根据需要随时启动，关闭测试
  - 提供设置重试机制提高成功率
- 用户模块 users
    - 身份鉴权，jwt 校验等
- 工作日志 notes
    - 编辑日志，发布等
- 监控和告警
  - 系统运行状态的监控和告警
  - 接入 Prometheus 监控
    - 统计 GIN 的 HTTP 接口
    - 统计 GORM 的 执行时间
    - 统计 HTTP 业务错误码
- 数据存储优化：
  - 根据数据特点选择合适的存储方案
  - 考虑使用NoSQL数据库存储非结构化数据
  - MySQL 用于存储结构化数据，例如用户信息、API 配置信息等
  - MongoDB 用于存储半结构化数据，例如测试报告、测试执行日志等

### Installation

#### core-engine

快速部署&&本地调试

```shell
$ git clone https://github.com/half-coconut/gocopilot.git
$ cd gocopilot/core-engine
$ docker-compose up -d # 安装依赖，注意 MYSQL 没有数据持久化
$ make docker # Golang 1.24, Ubuntu 24.04
```

使用 k8s 部署应用

- 部署 mysql，并持久化数据

```shell
$ kubectl apply -f k8s-mysql-pv.yaml
$ kubectl apply -f k8s-mysql-pvc.yaml
$ kubectl apply -f k8s-mysql-deployment.yaml
$ kubectl apply -f k8s-mysql-service.yaml
```

- 部署 core-engine

```shell
$ kubectl apply -f k8s-coreengine-deployment.yaml
$ kubectl apply -f k8s-coreengine-service.yaml
$ kubectl logs <pod-name> 
```

- 部署 redis

```shell
$ kubectl apply -f k8s-redis-deployment.yaml
$ kubectl apply -f k8s-redis-service.yaml
$ redis-cli -h <ip> -p 30003
```

安装 helm 和 ingress-nginx

```shell
# 安装 helm
$ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm- 3
$ chmod 700 get_helm.sh
$ ./get_helm.sh
# 使用 helm 安装 ingress-nginx
$ helm upgrade --install ingress-nginx ingress-nginx \
--repo https://kubernetes.github.io/ingress-nginx \
--namespace ingress-nginx --create-namespace
```

- 部署 nginx
- 部署 ingress

```shell
$ kubectl apply -f k8s-ingress-nginx.yaml
```


使用 k8s 集群部署应用

- 本地安装 K8s 集群的办法 Kind
- 安装 Kind：https://kind.sigs.k8s.io/docs/user/quick-start/#installation
- 执行 kind create 命令，创建 K8s 集群

使用 kubectl create service 命令创建 Service：

```shell
kubectl create service clusterip <docker-image> --tcp=5000:5000
service/<docker-image> created
```

使用 kubectl create ingress 命令创建 Ingress：

```shell
kubectl create ingress <docker-image> 
ingress.networking.k8s.io/<docker-image> created
```

部署 Ingress-nginx：

```shell
kubectl create -f ingress-nginx.yaml
namespace/ingress-nginx created
serviceaccount/ingress-nginx created
serviceaccount/ingress-nginx-admission created
......
```

#### react-pilot

本地安装，启动：

```shell
$ npm install # Node.js v20.10.0 以上
$ npm run dev
```
