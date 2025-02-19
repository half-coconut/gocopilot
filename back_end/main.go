package main

import (
	_ "go.uber.org/zap"
	_ "gorm.io/driver/mysql"
)

func main() {
	server := InitWebServer()
	server.Run(":3002")
}
