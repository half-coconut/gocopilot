# testcopilot


![testengine drawio](https://github.com/user-attachments/assets/b12b1ac1-a31a-418b-8095-ef353373f883)



- 旧版架构图如上，缺少 Job 模块和消息队列的实现；
- 后期更新为微服务架构，新增缓存方案，消息队列，分布式任务调度等模块；

#### 后端 TestEngine 部署
快速部署&&本地调试
```shell
$ git clone https://github.com/half-coconut/TestCopilot.git
$ cd TestCopilot/core-engine
$ docker-compose up -d # 安装依赖
$ make docker 
```
将容器镜像部署到k8s

Makefile 编译，Dockerfile 生成镜像，使用 k8s 集群部署应用

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

- Pod 会被 Deployment 工作负载管理起来，例如创建和销毁等；
- Service 相当于弹性伸缩组的负载均衡器，它能以**加权轮训**的方式将流量转发到多个 Pod 副本上；
- Ingress 相当于集群的外网访问入口；


#### 前端 TestPilot 部署
安装和启动：
```shell
npm install
npm run dev
```
