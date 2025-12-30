package ioc

import (
	"classroom-analysis/internal/web"
	"classroom-analysis/internal/ws"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitGin 初始化Gin引擎
func InitGin(
	loggers []gin.HandlerFunc,
	patientHandler *web.PatientHandler,
	fileHandler *web.FileHandler,
	userHandler *web.UserHandler,
	healthManagerHandler *web.HealthManagerHandler,
	roomHandler *web.RoomHandler,
	bedHandler *web.BedHandler,
	careLevelHandler *web.CareLevelHandler,
	dietPlanHandler *web.DietPlanHandler,
	customerHandler *web.CustomerHandler,
	recordHandler *web.RecordHandler,
	serviceHandler *web.ServiceHandler,
	careRecordHandler *web.CareRecordHandler,
	wsMgr *ws.WebSocketManager,
) *gin.Engine {
	engine := gin.Default()
	// 1. 注册基础全局中间件（不包含 JWT）
	engine.Use(gin.Recovery(), gin.Logger())
	for _, m := range loggers {
		engine.Use(m)
	}
	// 2. 公共路由分组（不需要 JWT）
	// 假设你的登录注册在 userHandler 里
	userHandler.RegisterRoutes(engine)
	patientHandler.RegisterRoutes(engine)
	authGroup := engine.Group("/")
	authGroup.Use(GetJWTMiddleware()) // 只在这个组里应用 JWT
	{
		// 注册文件上传相关路由
		fileHandler.RegisterRoutes(authGroup)

		// 注册健康管家相关路由
		healthManagerHandler.RegisterRoutes(authGroup)
		// 注册房间相关路由
		roomHandler.RegisterRoutes(authGroup)
		// 注册床位相关路由
		bedHandler.RegisterRoutes(authGroup)
		// 注册护理级别相关路由
		careLevelHandler.RegisterRoutes(authGroup)
		// 注册膳食计划相关路由
		dietPlanHandler.RegisterRoutes(authGroup)
		// 注册客户相关路由
		customerHandler.RegisterRoutes(authGroup)
		// 注册登记记录相关路由
		recordHandler.RegisterRoutes(authGroup)
		// 注册服务相关路由
		serviceHandler.RegisterRoutes(authGroup)
		// 注册护理记录相关路由
		careRecordHandler.RegisterRoutes(authGroup)
	}
	// Swagger文档路由
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket 路由
	engine.GET("/api/ws/alerts", wsMgr.Handler)

	return engine
}
