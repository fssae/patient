package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/middleware"
	"net/http"

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

// Login 患者登录
// @Summary      患者登录
// @Description  患者登录接口，使用手机号和密码登录
// @Tags         患者端
// @Accept       json
// @Produce      json
// @Param        request  body      domain.PatientLoginRequest  true  "登录信息"
// @Success      200      {object}  map[string]interface{}   "登录成功"
// @Failure      400      {object}  map[string]interface{}   "请求参数错误"
// @Failure      401      {object}  map[string]interface{}   "用户名或密码错误"
// @Router       /login [post]
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
// @Summary      患者注册
// @Description  患者注册接口
// @Tags         患者端
// @Accept       json
// @Produce      json
// @Param        request  body      domain.PatientRegisterRequest  true  "注册信息"
// @Success      200      {object}  map[string]interface{}      "注册成功"
// @Failure      400      {object}  map[string]interface{}      "请求参数错误"
// @Router       /register [post]
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
