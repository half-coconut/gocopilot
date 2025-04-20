# gocopilot

### Architecture

![coreengine drawio](https://github.com/user-attachments/assets/b12b1ac1-a31a-418b-8095-ef353373f883)

- 旧版架构图如上，缺少 Job 模块和消息队列的实现
- 后期更新为微服务架构，新增缓存方案，消息队列，分布式任务调度等模块

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
    - 长短任务模式，按照超时时间或 cron 表达式执行
    - 支持分布式、并发执行任务调度
    - 开启，关闭单个定时任务
- 用户模块 users
    - 身份鉴权，jwt 校验等
- 工作日志 notes
    - 编辑日志，发布等

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

$ kubectl delete -f k8s-coreengine-deployment.yaml
$ kubectl delete -f k8s-coreengine-service.yaml
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

Deployment, Pod, Service, Ingress 的关系

- Pod 会被 Deployment 工作负载管理起来，例如创建和销毁等
- Service 相当于弹性伸缩组的负载均衡器，它能以**加权轮训**的方式将流量转发到多个 Pod 副本上
- Ingress 相当于集群的外网访问入口

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
