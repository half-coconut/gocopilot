//go:build k8s

// 使用 k8s 这个编译标签
package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(coreengine-mysql:11309)/coreengine",
	},
	Redis: RedisConfig{
		Addr: "coreengine-redis:11479",
	},
}
