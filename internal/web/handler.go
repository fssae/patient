package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Login 患者登录
func (h *PatientHandler) Login(c *gin.Context) {
	var req domain.PatientLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "登录成功",
		"data":    resp.Patient,
		"success": true,
		"token":   resp.Token,
	})
}

// Register 患者注册
func (h *PatientHandler) Register(c *gin.Context) {
	var req domain.PatientRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}

	err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		print("err")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "注册成功",
		"success": true,
	})
}
