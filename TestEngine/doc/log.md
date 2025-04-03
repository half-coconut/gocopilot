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