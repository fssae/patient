package ioc

import (
	"classroom-analysis/internal/web"

	"github.com/gin-gonic/gin"
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
) *gin.Engine {
	engine := gin.Default()

	middlewares := InitMiddlewares(loggers)
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}
	// 注册教师相关路由
	patientHandler.RegisterRoutes(engine)
	// 注册文件上传相关路由
	fileHandler.RegisterRoutes(engine)
	// 注册用户相关路由（注册、登录）
	userHandler.RegisterRoutes(engine)
	// 注册健康管家相关路由
	healthManagerHandler.RegisterRoutes(engine)
	// 注册房间相关路由
	roomHandler.RegisterRoutes(engine)
	// 注册床位相关路由
	bedHandler.RegisterRoutes(engine)
	// 注册护理级别相关路由
	careLevelHandler.RegisterRoutes(engine)
	// 注册膳食计划相关路由
	dietPlanHandler.RegisterRoutes(engine)
	// 注册客户相关路由
	customerHandler.RegisterRoutes(engine)
	// 注册登记记录相关路由
	recordHandler.RegisterRoutes(engine)
	// 注册服务相关路由
	serviceHandler.RegisterRoutes(engine)

	return engine
}
