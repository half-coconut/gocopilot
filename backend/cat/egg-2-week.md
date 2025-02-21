### Egg-2-week
- Gin, Gorm 的使用
- 分层的项目结构
- 熟悉 interface 和调用结构体的使用
- 异常处理
- Service 层，做密码加密

1. 这里使用了 gin，gorm

2. 不同的分层的使用：

```
web         -> service          -> repository   -> dao
跟http打交道 -> 主要的业务逻辑在这里 -> 数据存储的抽象 -> 数据库操作
```

3. 注意interface 和调用的结构体的区别：

```go
type UserDAO interface {
	Insert(ctx context.Context, user dao.User) error
}

type GORMUserDAO struct {
	// 这里是 gorm的db 实现的结构体
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	// 在这里传入的 db，是这个调用的结构体，GORMUserDAO
	// 返回的是接口 UserDAO
	// 注意区分这2者的区别
	return &GORMUserDAO{
		db: db,
	}
}
```
4. 已处理的问题：password 数据库里没有存上，repository 没传 Password。
5. 在 Service 层处理 password 加密。
6. 异常处理，层层传出。体会 web 层和 service 层的报错的差别。
   - `邮箱冲突`
   - `邮箱不存在`
   - `邮箱/用户或者密码不正确`
7. 响应 JSON 处理。定义结构体 Result。
8. Repository 层处理 domainToEntity, entityToDomain 的数据转换。
9. Session 的处理，注意使用的方式:
   ```shell
   // 使用 cookie 存储session 数据
   store := cookie.NewStore([]byte("secret"))
   // 注册 sessions 中间件，并指定 session 名称和存储方式
   server.Use(sessions.Sessions("ssid", store))
   ```
10. Edit 编辑功能，修改需要注意：
   ```
   domain.User 的其它字段，尤其是密码、邮箱和手机，
   修改都要通过别的手段，邮箱和手机都要验证，密码更加不用多说了
   ```
11. 记录日志
```go
// log 方式一：使用 gin.log 文件
logFIle, _ := os.Create("gin.log")
log.SetOutput(logFIle)
// 方式二：使用 go 标准输出到控制台
log.SetOutput(os.Stdout)
```
12. 使用 session 作为身份鉴权
```go 
// 使用 cookie 存储session 数据
store := cookie.NewStore([]byte("secret"))

// 基于 Redis 的实现：
store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	[]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
if err != nil {
	panic(err)
}

// 注册 sessions 中间件，并指定 session 名称和存储方式
server.Use(sessions.Sessions("ssid", store))
```
13. 几个传参的方式
```go
// 静态路由
server.GET("/hello", func (context *gin.Context) {
context.String(http.StatusOK, "hello world")
})
// 参数路由
server.GET("/users/:name", func (context *gin.Context) {
name := context.Param("name")
context.String(http.StatusOK, "这是你传过来的名字：%s", name)
})
// 查询参数
server.GET("/info", func (context *gin.Context) {
id := context.Query("id")
context.String(http.StatusOK, "这是你传过来的 ID 是：%v", id)
})
// 通配符匹配
server.GET("/views/*.html", func (context *gin.Context) {
path := context.Param(".html")
context.String(http.StatusOK, "匹配上的值是：%s", path)
})

```
14. 使用 middleware 的简单示例
```go
server.Use(func (context *gin.Context) {
log.Println("第一个middleware")
}, func (context *gin.Context) {
log.Println("第二个middleware")
})

```