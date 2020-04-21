package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"lipeining.com/poem-gin/controller"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", controller.Ping)
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", controller.Home)
	r.GET("/tbl_tmp_id", controller.TblTmpID)
	r.GET("/tmp_id", controller.TmpID)
	r.GET("new_tbl", controller.NewTbl)
	r.GET("build_tmp_table", controller.BuildTmpTable)
	r.GET("build_table", controller.BuildTable)
	r.GET("infile", controller.Infile)
	r.GET("import_tang_author", func(c *gin.Context) {
		// http://localhost:8000/import
		filename := c.DefaultQuery("dynasty", "tang")
		fmt.Println(filename)
		c.String(http.StatusOK, "ok")
	})
	return r
}
