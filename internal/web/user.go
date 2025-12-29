package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册用户相关路由
func (h *UserHandler) RegisterRoutes(server gin.IRouter) {
	// 注册和登录不需要JWT验证
	server.POST("/api/user/register", h.Register)
	server.POST("/api/user/login", h.Login)
}

// Register 用户注册
// @Summary      用户注册
// @Description  用户注册接口，需要提供手机号、密码、姓名、年龄、性别
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        request  body      domain.UserRegisterRequest  true  "注册信息"
// @Success      200      {object}  map[string]interface{}      "注册成功"
// @Failure      400      {object}  map[string]interface{}      "请求参数错误"
// @Router       /user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
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

// Login 用户登录
// @Summary      用户登录
// @Description  用户登录接口，使用手机号和密码登录
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Param        request  body      domain.UserLoginRequest  true  "登录信息"
// @Success      200      {object}  map[string]interface{}   "登录成功，返回token和用户信息"
// @Failure      400      {object}  map[string]interface{}   "请求参数错误"
// @Failure      401      {object}  map[string]interface{}   "手机号或密码错误"
// @Router       /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
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
		"data":    resp.User,
		"success": true,
		"token":   resp.Token,
	})
}
