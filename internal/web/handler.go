package web

import (
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/middleware"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	svc *service.PatientService
}

func NewPatientHandler(svc *service.PatientService) *PatientHandler {
	return &PatientHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册教师相关路由
func (h *PatientHandler) RegisterRoutes(server *gin.Engine) {
	// 教师登录  注册 - 不需要JWT验证
	server.POST("/login", h.Login)
	server.POST("/register", h.Register)

	// 需要JWT验证的接口
	PatientGroup := server.Group("/Patient")
	PatientGroup.Use(
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/Patient/login").
			IgnorePaths("/Patient/register").
			//ReWritPaths("/Patient/login").
			Build(),
	)
}
