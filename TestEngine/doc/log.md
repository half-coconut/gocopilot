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

{"email": "test@123.com","password":"Cc12345!"}

http://127.0.0.1:3002/users/login
```

```shell
https://nvbtdjgdbhgsgccuwkap.supabase.co/storage/v1/object/public/avatars//avatar-e71dbd5f-bb33-4443-ade3-b6379c11555f-0.1510105863854192
```