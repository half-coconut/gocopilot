# Docker 环境配置

### Docker 安装

1 卸载旧版本
```shell
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done

```

2 添加 Docker 官方 GPG 密钥
```shell
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg

sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

```

3 将仓库添加到 APT 源中
```shell
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update

```

4 安装 Docker
```shell
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 查看 Docker 运行状态
systemctl status docker

[chen@dev workspace]$systemctl status docker
● docker.service - Docker Application Container Engine
     Loaded: loaded (/lib/systemd/system/docker.service; enabled; preset: enabled)
     Active: active (running) since Wed 2025-04-16 10:44:56 CST; 35s ago
TriggeredBy: ● docker.socket
       Docs: https://docs.docker.com
   Main PID: 397272 (dockerd)
      Tasks: 9
     Memory: 22.7M
        CPU: 320ms
     CGroup: /system.slice/docker.service
             └─397272 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
[chen@dev workspace]$
```
5 Docker 配置文件
 /etc/docker/daemon.json

```shell
sudo tee /etc/docker/daemon.json << EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "insecure-registries": [],
  "live-restore": true,
  "bip": "172.16.0.1/24",
  "storage-driver": "overlay2",
  "registry-mirrors": ["https://registry.docker-cn.com"],
  "data-root": "/data/lib/docker",
  "log-driver": "json-file",
  "dns": [],
  "default-runtime": "runc",
  "log-opts": {
    "max-size": "100m",
    "max-file": "10"
  }
}
EOF

# 重启 docker
sudo systemctl restart docker

# 测试 docker 是否安装成功
sudo docker run hello-world
```
### Docker 安装后配置

1 使用 non-root 用户操作 docker
```shell
sudo -i groupadd docker  # 创建 `docker` 用户组
sudo -i usermod -aG docker $USER # 将当前用户添加到 `docker` 用户组下
newgrp docker # 重新加载组成员身份
docker run hello-world # 确认能够以普通用户实际使用 docker

```
2 配置 docker 开启启动
```shell
sudo systemctl enable docker.service # 设置 docker 开机启动
sudo systemctl enable containerd.service # 设置 containerd 开机启动

# 安装 docker-compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

```

3 排障
- 1.通过 kill -9 xxx 杀死所有的 containerd 进程；
- 2.重启 containerd 服务，systemctl restart containerd;
- 3.重启 docker 服务：systemctl restart docker.

### OneX 项目容器化部署
1 通过 github.com host 来解决网速慢的问题
```shell
sudo tee -a /etc/hosts << EOF
140.82.114.4 github.com
EOF
```

2 下载源码 onex
```shell
mkdir -p $WORKSPACE/golang/src/github.com/onexstack
cd $WORKSPACE/golang/src/github.com/onexstack
git clone https://github.com/onexstack/onex

```
3 初始化工作区
```shell
cd $WORKSPACE/golang/src/github.com/onexstack
go work init
go env GOWORK
go work use ./onex

```

4 一键部署 OneX 服务
```shell
cd $WORKSPACE/golang/src/github.com/onexstack/onex
make docker-install
```
5 启动 Swagger API 文档
```shell
make tools.install.swagger
make serve-swagger

http://119.28.193.119:65534/docs
```

6 简单测试 声明式 API
```shell
cd ${ONEX_ROOT}
source manifests/env.local
# 查看 onex-apiserver 中 onex 自定义资源
kubectl --kubeconfig=${ONEX_ADMIN_KUBECONFIG} api-resources | egrep apps.onex.io
# 查看 minerset 资源列表
kubectl --kubeconfig=${ONEX_ADMIN_KUBECONFIG} get minerset

```



