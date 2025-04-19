### Egg-3-week

1. 查看端口占用情况，并终止

```shell
lsof -i :3002
kill -9 <PID>
```

2. 使用 JWT，重写 Login，声明 UserClaims 结构体，定义过期时间

- 修改 middleware 中跨域部分，添加 header。

```shell
# 为了使用 jwt
AllowedHeaders:   []string{"Content-Type", "Authorization"},
osedHeaders:   []string{"x-jwt-token"}
  ```

3. 增加 init.sql
4. k8s 部署，JWT 接口安全等
5. wrk 压测：Login, Signup, Profile等接口

- 命令行发压
    - t12: 使用 12 个线程进行压测。
    - c400: 保持 400 个并发连接。
    - d30s: 持续压测 30 秒。
    - url: 目标 API 地址。

```shell
wrk -t12 -c400 -d30s http://localhost:3002/hello

Running 30s test @ http://localhost:3002/hello
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    44.89ms   57.51ms 716.95ms   91.00%
    Req/Sec   542.41    428.77     3.46k    83.15%
  188358 requests in 30.10s, 16.98MB read
  Socket errors: connect 157, read 162, write 0, timeout 0
  Non-2xx or 3xx responses: 165970
Requests/sec:   6257.92
Transfer/sec:    577.62KB
```

- 使用 lua 脚本
    - ` wrk -t12 -c400 -d30s -s ./script/wrk/profile.lua 'http://localhost:3002'`
      `
    - -s test.lua: 指定要使用的 Lua 脚本。

```lua
-- 设置请求头
wrk.headers["Content-Type"] = "application/json"

-- 定义请求函数
request = function()
    -- 构造请求体
    local body = '{"name": "test", "age": 30}'
    -- 发送 POST 请求
    return wrk.format("POST", "/users", wrk.headers, body)
end
```

- profile.lua 压测结果：

```shell
Running 30s test @ http://localhost:3002
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    94.61ms   88.66ms 764.54ms   87.96%
    Req/Sec   435.00    358.60     2.33k    85.55%
  148210 requests in 30.03s, 26.71MB read
  Socket errors: connect 0, read 468, write 1, timeout 0
  Non-2xx or 3xx responses: 148210
Requests/sec:   4935.36
Transfer/sec:      0.89MB

chenchen@chenchendeMacBook-Pro-2 egg-yolk % wrk -t12 -c400 -d30s -s ./script/wrk/profile.lua 'http://localhost:3002'
Running 30s test @ http://localhost:3002
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   919.26ms  473.05ms   2.00s    65.09%
    Req/Sec    24.10     23.44   207.00     84.94%
  6776 requests in 30.10s, 1.53MB read
  Socket errors: connect 0, read 389, write 0, timeout 1227
  Non-2xx or 3xx responses: 6776
Requests/sec:    225.13
Transfer/sec:     52.11KB

chenchen@chenchendeMacBook-Pro-2 egg-yolk % wrk -t12 -c400 -d30s -s ./script/wrk/profile.lua 'http://localhost:3002'
Running 30s test @ http://localhost:3002
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   940.97ms  430.56ms   2.00s    65.76%
    Req/Sec    26.28     32.40   262.00     88.45%
  7110 requests in 30.10s, 1.61MB read
  Socket errors: connect 0, read 399, write 3, timeout 1759
  Non-2xx or 3xx responses: 7110
Requests/sec:    236.22
Transfer/sec:     54.67KB
```

6. 压测(TODO: 扩展):

- 切换不同的加密算法，测试注册和登录接口。
- 在我准备的这些脚本的基础上，你在数据库中插入 100W 条用户数据，然后再测试登录接口。
- 在数据库中插入 1000W 用户数据，然后再测试登录接口。

7.
    8. val, err := json.Marshal(user) 是 Go 语言中用于将一个结构体或变量 user 编码为 JSON 格式字符串的代码。
8. 性能优化：缓存，给 Profile 接口加缓存，性能测试对比前后变化，增加方法 FindById。
9. https://console.cloud.tencent.com/smsv2
10. redis

```shell
redis-cli
redis-cli -p 6380

127.0.0.1:6379> get user:info:3
"{\"Id\":3,\"Email\":\"1234@qq.com\",\"Password\":\"$2a$10$59ylwNoq12DoldPb57zvsen1ZD5gEUu5bhXPtDCplRYgYLWRAr4gu\",\"Phone\":\"\",\"NickName\":\"\",\"Department\":\"\",\"Role\":\"\",\"Description\":\"\",\"Ctime\":\"2024-08-19T13:54:49.318+08:00\",\"Utime\":\"2024-08-19T13:54:49.318+08:00\"}"

127.0.0.1:6379> ttl user:info:3
(integer) 782
127.0.0.1:6379> ttl user:info:3
(integer) 777
127.0.0.1:6379> ttl user:info:3
(integer) 775
```

11. 依赖注入，

- 先安装，然后执行 go mod init 和 go mod tidy
- 创建 wire.go 文件
- 执行 wire

```bash
go install github.com/google/wire/cmd/wire@latest

go mod init egg-yolk
chenchen@chenchendeMacBook-Pro-2 egg-yolk % go mod tidy
go: finding module for package github.com/google/wire
go: found github.com/google/wire in github.com/google/wire v0.6.0

wire
```

- 当使用  `go build -tags wireinject`  命令编译项目时，
- Wire 会扫描 wire.go 文件，并根据 wire.Build 函数生成依赖注入代码，
- 最终生成一个 InitializeApp 函数，该函数负责创建和返回 App 实例。

#### 总结：

- `//go:build wireinject` 是一个构建约束标签，
- 用于控制代码在特定条件下是否参与编译。它通常与 Google Wire 依赖注入工具一起使用，
- 用于标记包含 Wire 配置的代码文件。

12. 关于 ioc

- IOC (Inversion of Control)，即控制反转，
- 是一种软件设计原则，它将对象的创建和依赖关系的管理从应用程序代码中转移到外部容器或框架中。

- 传统的依赖关系管理：
- 在传统的编程方式中，对象通常会直接创建和管理其依赖项。
- 例如，一个 UserService 对象可能需要创建一个 UserRepository 对象来访问用户数据：

```go
type UserService struct {
userRepository *UserRepository
}

func NewUserService() *UserService {
return &UserService{
userRepository: NewUserRepository(), // 直接创建依赖项
}
}
```

- 这种方式会导致代码耦合度高，难以维护和测试。

- IOC 的优势：

  降低耦合度： 对象不再直接依赖于具体的依赖项，而是依赖于抽象接口，从而降低了代码的耦合度，提高了代码的可维护性和可测试性。
  提高代码复用性： 依赖项可以被多个对象共享，提高了代码的复用性。
  简化配置： 依赖关系的管理由外部容器或框架负责，简化了应用程序的配置。

- IOC 的实现方式：

  依赖注入 (Dependency Injection, DI): 通过构造函数、Setter 方法或接口注入的方式将依赖项传递给对象。
  服务定位器 (Service Locator): 对象通过服务定位器来获取其依赖项。

- IOC 容器的职责：

  创建对象： 根据配置信息创建对象实例。
  管理依赖关系： 将依赖项注入到对象中。
  管理对象生命周期： 控制对象的创建、销毁和生命周期。

- 常见的 IOC 框架：

  Spring (Java): 一个功能强大的 IOC 容器和应用程序框架。
  Google Guice (Java): 一个轻量级的依赖注入框架。
  Dagger (Java/Android): 一个编译时依赖注入框架。
  Wire (Go): 一个 Go 语言的依赖注入工具。

#### 总结：

- IOC (控制反转) 是一种重要的软件设计原则，它可以帮助您降低代码耦合度，
- 提高代码的可维护性和可测试性。通过使用 IOC 容器或框架，
- 您可以将对象的创建和依赖关系的管理从应用程序代码中分离出来，从而简化应用程序的开发和维护。

13. 使用 vegeta 压测

```zsh
vegeta attack -targets=targets.txt -rate=50 -duration=30s | vegeta report

vegeta attack -targets=targets.txt -duration=10s -rate=10/s | tee results.bin
base results.bin | vegeta report > results.txt 

#-targets=targets.txt：
#这个参数指定了一个包含 HTTP 请求定义的文件。每行代表一个请求，可以包含请求方法（如 GET 或 POST）、URL、headers 和 body 等信息。

#-duration=10s：
#这个参数设置了压力测试的持续时间，这里是 10 秒。Vegeta 会在这个时间段内持续发送请求。

#-rate=10/s：
#这个参数定义了每秒发送的请求数量，这里是每秒 10 个请求。Vegeta 会尝试以指定的速率发送请求，但实际速率可能会受到系统性能和网络条件的限制。

#| tee results.bin：
#这是一个管道命令，tee 命令用于将输出同时发送到标准输出（控制台）和文件。在这里，results.bin 是一个文件，它将存储压力测试的结果。这样，你可以在测试结束后查看或分析这些结果。

# 留个问题：如何可视化
```

- 测试结果

```shell
chenchen@chenchendeMacBook-Pro-2 vegeta % vegeta attack -targets=targets.txt -rate=50 -duration=30s | vegeta report
Requests      [total, rate, throughput]         1500, 50.03, 50.03
Duration      [total, attack, wait]             29.984s, 29.981s, 3.503ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.308ms, 4.492ms, 3.483ms, 8.208ms, 10.187ms, 14.562ms, 30.093ms
Bytes In      [total, mean]                     523500, 349.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:1500  
Error Set:
```

14 参考 runnergo
`https://github.com/Runner-Go-Team/runnerGo?tab=readme-ov-file`

15 要查看 vegeta 的源码，特别是发送 HTTP 请求的部分，你应该重点关注以下几个文件和目录：

    `lib/vegeta.go`: 这是 vegeta 的核心文件，包含了主要的逻辑，包括构建 HTTP 请求、发送请求、处理响应等。 你可以在这里找到 Attack 结构体的定义以及 attack 函数，它们是发起攻击的核心部分。
    `lib/requester.go`: 这个文件包含了实际发送 HTTP 请求的代码。 你可以在这里找到 Requester 接口的定义以及它的不同实现，例如 DefaultRequester 和 TCPRequester。 这些实现使用了 Go 标准库的 net/http 包来发送 HTTP 请求。
    `lib/targets/targets.go`: 这个文件处理目标的解析和迭代。 你可以在这里找到如何从不同的输入源（例如文件、命令行参数）读取目标信息，并将其转换为 Target 结构体。

建议你按照以下步骤进行阅读：

    从 `lib/vegeta.go` 的 `attack` 函数开始，了解 vegeta 如何发起攻击的整体流程。
    然后深入到 `lib/requester.go` 中，查看 `Requester` 接口的不同实现，了解如何使用 `net/http` 包发送 HTTP 请求。
    最后，可以查看 `lib/targets/targets.go`，了解如何解析和迭代目标。

一些额外的提示：

    你可以使用代码编辑器或 IDE 来浏览代码，例如 VS Code、GoLand 等，它们可以提供代码跳转、函数调用关系图等功能，方便你理解代码结构。
    你可以尝试在本地运行 vegeta 的测试用例，并使用调试器单步调试代码，这可以帮助你更好地理解代码的执行流程。
