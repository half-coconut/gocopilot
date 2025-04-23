- ...
- 增加 Aicopilot
- 增加 user,profile,api 等

2025-03-14

- 修复了接口返回格式问题，`MiddlewareBuilder` 的 `responseWriter` 里对 `len(data) < 1024` 做了限制；
- 增加 debug 到 edit 接口，实现 http接口运行, 前端加相应调整

2025-03-30

- 增加了 debug history，以及 JOSN格式的展示
- 增加 Task 模块，新增 task，下一步完成 批量运行接口测试...报告展示等

2025-3-31

- 完成 Task 模块，debug 和 run 接口功能

```shell
{"Content-Type": "application/json",
"User-Agent": "PostmanRuntime/7.43.0"}
http://127.0.0.1:3002/users/login
{"email": "test@123.com","password":"Cc12345!"}


https://api.infstones.com/core/mainnet/6e97213d22994a2fae3917c0e00715d6
{"jsonrpc": "2.0", "method": "eth_accounts", "params": [], "id": 1}
{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 0}

```

```shell
https://nvbtdjgdbhgsgccuwkap.supabase.co/storage/v1/object/public/avatars//avatar-e71dbd5f-bb33-4443-ade3-b6379c11555f-0.1510105863854192
```

2025-04-03 关于为何使用基于MySQL抢占式分布式定时任务框架
1.redis 分布式锁
2.根据节点动态的调整-负载均衡

- go 里没有分布式调度平台
- 满足更多自定义的功能需要，比如开始，结束，
- 任务编排：自定义编排顺序，子任务 A -> 子任务 B，有向顺序的执行图
    - a1,a2,a3 任务, 执行成功，-> 任务 B
    - a1,a2 成功了,a3 成功无所谓 任务, 执行成功，-> 任务 B
- 负载均衡

2025-04-07

- 完成了 taskService 的 api,task debug, execute task，和前端页面展示;
- 计划下一步完成 report 和 task 服务的解耦，完成更清晰的 implement;


- 数据迁移方案

```shell
# 8.4 版本之后，客户端登录会有问题，记得安装指定版本的mysql
brew install mysql-client@8.4
brew unlink mysql
brew link mysql-client@8.4

docker exec -it <container_id_or_name> bash

# 将原表的数据 dump 下来
tmp mysqldump -h 127.0.0.1 --port 13316 -u root -p coreengine interactives user_like_bizs user_collection_bizs collections > intr_4.7.sql
# 将 docker 里的数据 cp 到本地目录文件
docker cp <container_id_or_name>:/tmp/intr_4.7.sql ./intr_4.7.sql

# 新建一个数据库 比如 coreengine_intr
# 登录 mysql 数据库
mysql -h 127.0.0.1 --port 13316 -uroot -proot

# 切换到新的数据库
create database if not exists coreengine_intr;
use coreengine_intr;
source intr_4.7.sql
```

- 在ecs中操作：
- ssh root@47.239.187.141

```shell
# 安装 apt, git, docker, npm, 下载 github 仓库
# 安装 docker 
ping google.com
sudo rm /usr/share/keyrings/docker-archive-keyring.gpg
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo apt update
sudo apt install docker.io
# 启动 docker 服务
sudo systemctl status docker
sudo systemctl start docker
sudo usermod -aG docker ubuntu


sudo nano /etc/docker/daemon.json
{
  "registry-mirrors": ["https://i44jb9ta.mirror.aliyuncs.com"]
}

docker-compose up -d
# 安装 golang-go
apt  install golang-go
go build -o coreengine .

```

nginx

```shell
vim /etc/nginx/sites-available/default
sudo nginx -t
sudo systemctl restart nginx

# nginx 日志
/var/log/nginx/error.log

npm install
npm run build
```

前端: http://47.239.187.141/login
后端: http://47.239.187.141:3002/users/login

```shell
server {
    listen 80;
    server_name 47.239.187.141;
    
    # 静态文件路径
    root /root/TestCopilot/react-pilot/dist;
    index index.html;

    # 处理前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 代理 API 请求到 Go 后端
     location /api/ {
      proxy_pass http://localhost:3002/;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
  
      # 添加 CORS 头部
      add_header 'Access-Control-Allow-Origin' 'http://47.239.187.141' always;
      add_header 'Access-Control-Allow-Credentials' 'true' always;  
      add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
      add_header 'Access-Control-Allow-Headers' 'Content-Type, Authorization' always;
  
      # 处理 OPTIONS 预检请求
      if ($request_method = OPTIONS) {
          add_header 'Access-Control-Allow-Origin' 'http://47.239.187.141' always;
          add_header 'Access-Control-Allow-Credentials' 'true' always;
          add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
          add_header 'Access-Control-Allow-Headers' 'Content-Type, Authorization' always;
          add_header 'Access-Control-Max-Age' 1728000 always;
          add_header 'Content-Type' 'text/plain charset=UTF-8' always;
          add_header 'Content-Length' 0 always;
          return 204;
      }
}
}

scp root@47.239.187.141:/etc/nginx/sites-available/default /Users/chenchen/Desktop/default
```

2025-04-09

- 设计任务调度服务，web，创建 CronJob
- Job 分类，http 请求任务，内部方法任务(使用svc.func()内部调用)，定时任务，需要提取出来 schedule 来执行

2025-04-10

- Job 完成运行一个定时任务的接口，

2025-04-11

- 完成开启、关闭的功能
- 开启时释放任务，设置下一次运行时间，抢占式获取执行锁，运行任务后，释放任务，设置下一次运行时间
- 关闭功能，调用关闭后，状态为 cronjobStatusPaused，此时运行任务的 goroutine 在间隔后执行任务暂停

2025-04-14

```shell
# 列出当前使用端口 8090 的所有进程
lsof -i:8090

brew install etcd
etcdctl --version
etcdctl --endpoints=localhost:12379 get service/user/127.0.0.1:8090

➜  ~ etcdctl --endpoints=localhost:12379 get service/user/127.0.0.1:8090
service/user/127.0.0.1:8090
{"Op":0,"Addr":"127.0.0.1:8090","Metadata":"2025-04-14 13:58:11.399592417 +0800 CST m=+464.031308876"}

➜  ~ etcdctl --endpoints=localhost:12379 get service/user --prefix      
service/user/198.18.0.1:8090
{"Op":0,"Addr":"198.18.0.1:8090","Metadata":"2025-04-14 15:31:02.734662959 +0800 CST m=+57.030341418"}
➜  ~ 
```

2025-04-17

- 新建前端 job 页面

2025-04-18

- 新增日志输出配置项，标准输出，日志级别等

2025-04-19

- docker 镜像上传 Docker Hub

```shell
➜  core-engine git:(main) ✗ docker tag core-engine:v0.0.1 halfcoconut/gocopilot:core-engine
➜  core-engine git:(main) ✗ docker push  halfcoconut/gocopilot:core-engine

```

2025-04-23

- 分离出 report service
- 存入 mongoDB，集合 debug_logs，summary
    - 使用雪花算法生成 id
    - 雪花算法之所以被广泛使用，是因为它能够提供全局唯一、高性能、准有序、分布式友好和可配置的 ID 生成方案，满足了分布式系统中对
      ID 生成的各种需求

```shell
# 查询某个任务的 debug_logs 有多少条
$ db.getCollection("debug_logs").find({"task_id" : 2}).count()
```

计划明天完成：

- summary 存入数据库
- 完成使用消息队列，将 task 和 report 解耦，
- 消息队列解决的问题：异步，削峰，解耦
- 场景：秒杀，支付(支付失败，考虑延迟队列)，
- 关于消息队列的应该设置多少个分区，生产者和消费者的计算：
    - max(发送者总速率/单一分区写入速率，发送者总速率/单一消费者总速率) + buffer
- 消息积压的问题：同一个分区，只能有一个消费者
- 消息有序执行的问题：同一个分区，保证消息的有序