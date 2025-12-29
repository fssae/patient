package main

import (
	"log"

	_ "classroom-analysis/docs" // swagger docs
)

// @title           东软顾养中心 API
// @version         1.0
// @description     东软顾养中心管理系统API文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	//初始化
	setupEnvironment()
	//启动定时任务
	startCronJobs()
	//启动服务
	app := InitWebServer()
	if err := app.Start(); err != nil {
		log.Fatalf("应用启动失败: %v", err)
	}
}
