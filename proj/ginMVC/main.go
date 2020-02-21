package main

import (
	"github.com/davidddw2017/panzer/proj/ginMvc/drivers"
	"github.com/davidddw2017/panzer/proj/ginMvc/server"
	"github.com/gin-gonic/gin"
)

var httpServer *gin.Engine

func main() {

	defer drivers.MySQLDB.Close()
	// 启动服务
	server.Run(httpServer)
}
