package main

import (
	"fmt"

	"github.com/davidddw2017/panzer/proj/ginMvc/config"
	"github.com/davidddw2017/panzer/proj/ginMvc/driver"
	"github.com/davidddw2017/panzer/proj/ginMvc/routes"
	"github.com/gin-gonic/gin"
)

var httpServer *gin.Engine

func main() {
	defer driver.MysqlDb.Close()

	// 启动服务
	Run(httpServer)
}

// 配置并启动服务
func Run(httpServer *gin.Engine) {
	// 服务配置
	serverConfig := config.SystemConfig.Server

	// gin 运行时 release debug test
	gin.SetMode(serverConfig.Env)

	httpServer = gin.Default()

	// 配置视图
	if "" != serverConfig.ViewPattern {
		httpServer.LoadHTMLGlob(serverConfig.ViewPattern)
	}

	if "" != serverConfig.StaticPattern {
		httpServer.Static("/vender/static", serverConfig.StaticPattern)
	}

	// 注册路由
	routes.RegisterRoutes(httpServer)

	serverAddr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	// 启动服务
	err := httpServer.Run(serverAddr)
	if nil != err {
		panic("server run error: " + err.Error())
	}
}
