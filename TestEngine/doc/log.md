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
tmp mysqldump -h 127.0.0.1 --port 13316 -u root -p testengine interactives user_like_bizs user_collection_bizs collections > intr_4.7.sql
# 将 docker 里的数据 cp 到本地目录文件
docker cp <container_id_or_name>:/tmp/intr_4.7.sql ./intr_4.7.sql

# 新建一个数据库 比如 testengine_intr
# 登录 mysql 数据库
mysql -h 127.0.0.1 --port 13316 -uroot -proot

# 切换到新的数据库
create database if not exists testengine_intr;
use testengine_intr;
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
go build -o testengine .

```


nginx 

vim /etc/nginx/sites-available/default
- 
- nginx 日志
/var/log/nginx/error.log
```shell
server {
        listen 80;
        listen [::]:80;

        server_name 47.239.187.141;

        root /root/TestCopilot/TestPilot/dist;
        index index.html;

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /api/ {
       proxy_pass http://localhost:3002/;
       proxy_set_header Host $http_host;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
       proxy_set_header X-Forwarded-Proto $scheme;

       # 添加 CORS 头部
       add_header 'Access-Control-Allow-Origin' '*' always;
       add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
       add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization' always;
       add_header 'Access-Control-Expose-Headers' 'Content-Length,Content-Range' always;

       if ($request_method = OPTIONS) {
           add_header 'Access-Control-Allow-Origin' '*' always;
           add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
           add_header 'Access-Control-Allow-Headers' 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization' always;
           add_header 'Access-Control-Max-Age' 3600 always;
           add_header 'Content-Type' 'text/plain charset=UTF-8' always;
           add_header 'Content-Length' 0 always;
           return 204;
       }
   }

}
sudo nginx -t
sudo systemctl restart nginx

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
    root /root/TestCopilot/TestPilot/dist;
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
        add_header 'Access-Control-Allow-Origin' '*' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
        add_header 'Access-Control-Allow-Headers' 'Content-Type, Authorization' always;

        # 处理 OPTIONS 预检请求
        if ($request_method = OPTIONS) {
            add_header 'Access-Control-Allow-Origin' '*' always;
            add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS, PUT, DELETE' always;
            add_header 'Access-Control-Allow-Headers' 'Content-Type, Authorization' always;
            add_header 'Access-Control-Max-Age' 1728000 always;
            add_header 'Content-Type' 'text/plain charset=UTF-8' always;
            add_header 'Content-Length' 0 always;
            return 204;
        }
    }
}

```