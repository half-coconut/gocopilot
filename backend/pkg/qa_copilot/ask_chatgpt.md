安装 cobra，初始化项目

```shell
1.go mod init k8scopilot

2.go install github.com/spf13/cobra-cli@latest

3.配置命令名称
vim ~/.zshrc
alias cobra-cli="~/go/bin/cobra-cli"
source ~/.zshrc

4.cobra-cli init --author "aiops" --license mit 

```

构建 k8scopilot，运行

```shell
go build -o k8scopilot
./k8scopilot
```

创建命令

```shell
cobra-cli add hello

chenchen@chenchendeMacBook-Pro-2 k8scopilot % cobra-cli add hello
hello created at /Users/chenchen/Downloads/plgo-main/book/aiops/module_7/k8scopilot

```

添加 多级自命令，subcommand

```shell
cobra-cli add world -p 'helloCmd'

```

kubectl ask chatgpt

```shell
cobra-cli add ask
cobra-cli add chatgpt -p 'askCmd'

```

k8s

```shell
帮我生成一个deployment，镜像是nginx
```

对话实例

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot  
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot ask chatgpt
I'm k8s copilot, what can I do for you?
> 帮我生成一个deployment，镜像是nginx
OpenAI hopes to request the function: generateAndDeploymentResource, parameters: {"user_input":"生成一个deployment，镜像是nginx"}
> 查询 kube-system NS 的 Pod
OpenAI hopes to request the function: queryResource, parameters: {"namespace":"kube-system","resource_type":"Pod"}
> 帮我删除 default NS 的 Pod，名称是 pod-xxx
OpenAI hopes to request the function: deleteResource, parameters: {"namespace":"default","resource_name":"pod-xxx","resource_type":"Pod"}

```

生成 yaml 示例

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot  
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot ask chatgpt
I'm k8s copilot, what can I do for you?
> 帮我生成一个deployment，镜像是nginx
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
> exit
Bye!


```

创建一个 deployment 和 pod

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot  
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot ask chatgpt
I'm k8s copilot, what can I do for you?
> 帮我创建一个deployment, 镜像是nginx
YAML content:
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80

Deployment successful.
> 帮我创建一个nginx pod，镜像是nginx，port是8080
YAML content:
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 8080

Deployment successful.


```

查询 pod, deploy, service示例

```shell
chenchen@chenchendeMacBook-Pro-2 k8scopilot % go build -o k8scopilot
chenchen@chenchendeMacBook-Pro-2 k8scopilot % ./k8scopilot ask chatgpt
I'm k8s copilot, what can I do for you?
> 查询一下 default NS 的 pod 
Found pod: demo-1-6d9fbb798-khm5p
Found pod: nginx-deployment-96b9d695-8gqsm
Found pod: nginx-deployment-96b9d695-g8sts
Found pod: nginx-deployment-96b9d695-klw8v
Found pod: nginx-pod

> 查询一下 default ns 的 deploy
Found deployment: demo-1
Found deployment: nginx-deployment
> 查询 kube-system NS 的 service
Found service: kube-dnsFound service: metrics-server
> 查询 default ns 的 service
Found service: db
Found service: hello-egg
Found service: kubernetes
Found service: redis
Found service: result
Found service: vote
> exit
Bye!

```