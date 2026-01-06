package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/ginx"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(server gin.IRouter) {
	server.POST("/api/user/register", ginx.WrapBody[domain.UserRegisterRequest](h.Register))
	server.POST("/api/user/login", ginx.WrapBody[domain.UserLoginRequest](h.Login))
}

func (h *UserHandler) Register(c *gin.Context, req domain.UserRegisterRequest) (ginx.Result, error) {
	if err := h.svc.Register(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("注册成功"), nil
}

func (h *UserHandler) Login(c *gin.Context, req domain.UserLoginRequest) (ginx.Result, error) {
	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		return ginx.Fail(401, err.Error()), nil
	}
	return ginx.Result{
		Code:    200,
		Msg:     "登录成功",
		Success: true,
		Data: gin.H{
			"user":  resp.User,
			"token": resp.Token,
		},
	}, nil
}
