package main

import (
	_ "go.uber.org/zap"
	_ "gorm.io/driver/mysql"
)

func main() {
	app := InitWebServer()
	for _, c := range app.Consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.Server
	server.Run(":3002")
}
