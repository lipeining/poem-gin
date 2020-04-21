package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"lipeining.com/poem-gin/config"
	"lipeining.com/poem-gin/database"
)

func main() {
	fmt.Println(config.Config)
	database.New(config.Config.Xorm.User, config.Config.Xorm.Passwd, config.Config.Xorm.Database, config.Config.Xorm.SecurePivFile)

	// 1.创建路由
	r := setupRouter()
	// fmt.Printf("config %v", config.Config)
	// fmt.Printf("config.DB %v", config.Config.DB)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
