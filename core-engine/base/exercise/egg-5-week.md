### Egg-5-week

- 接入配置模块
- 接入日志模块

##### 封装日志接口
```
egg_yolk/pkg/logger/type.go
三种封装日志接口的风格比较：
Logger 兼容性最好
LoggerV1 认同参数要有名字
LoggerV2 有完善的代码流程，否则不建议使用
```
##### 使用 middleware 封装 logger builder
```
egg_yolk/pkg/ginx/middlewares/logger/builder.go
用于特殊处理的日志打印，比如打印 http 请求
```
##### channel
在 Go 语言中，通道 (chan) 可以用于 goroutine 之间的通信。通道可以是双向的，也可以是单向的：

    双向通道 (`chan T`): 既可以发送数据，也可以接收数据。
    只发送通道 (`chan<- T`): 只能发送数据，不能接收数据。
    只接收通道 (`<-chan T`): 只能接收数据，不能发送数据。

##### kube
```shell
mkdir /home/zhihao/.aws

vim credentials

[dev-cloud-iam-eks]
region = us-east-1
role_arn=arn:aws:iam::221091472662:role/dev-cloud-iam-eks-user-infpools-io
credential_source=Ec2InstanceMetadata

export AWS_PROFILE=dev-cloud-iam-eks
aws eks update-kubeconfig --region us-east-1 --name dev-cloud-eks-cluster-infpools-io

kubectl get pods -n dev-cloud-qa-work
kubectl get pods -n dev-cloud-dos-work

kubectl exec -it websitemonitor-7db7dc96f7-n4ch8 /bin/bash -n dev-cloud-qa-work
```

##### runtime.NumCPU()
获取逻辑 CPU 核心数： runtime.NumCPU() 返回一个整数，表示当前机器可用的逻辑 CPU 核心数。
辅助并发控制： 这个函数的返回值可以用来辅助设置 Go 程序的并发度，例如，可以将 runtime.GOMAXPROCS 设置为 runtime.NumCPU()，以便让 Go 程序充分利用所有可用的 CPU 核心。

##### 原子函数
```shell
atomic.AddInt64(&counter, 1)
// 安全的读
atomic.LoadInt64(&counter)
// 安全的写
atomic.StoreInt64(&counter,1)
```
##### 生成安全的随机数
在 Go 1.20 及以后的版本中，建议使用 rand.NewSource(seed) 和 rand.New(source)  来创建新的随机数生成器，以替代已弃用的 rand.Seed() 函数。如果需要更安全的随机数，可以使用 crypto/rand 包。
