package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/middleware"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	svc *service.PatientService
}

func NewPatientHandler(svc *service.PatientService) *PatientHandler {
	return &PatientHandler{svc: svc}
}

func (h *PatientHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/login", ginx.WrapBody[domain.PatientLoginRequest](h.Login))
	server.POST("/register", ginx.WrapBody[domain.PatientRegisterRequest](h.Register))

	PatientGroup := server.Group("/Patient")
	PatientGroup.Use(
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/Patient/login").
			IgnorePaths("/Patient/register").
			Build(),
	)
}

func (h *PatientHandler) Login(c *gin.Context, req domain.PatientLoginRequest) (ginx.Result, error) {
	resp, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		return ginx.Fail(401, err.Error()), nil
	}
	return ginx.Result{
		Code:    200,
		Msg:     "登录成功",
		Success: true,
		Data: gin.H{
			"patient": resp.Patient,
			"token":   resp.Token,
		},
	}, nil
}

func (h *PatientHandler) Register(c *gin.Context, req domain.PatientRegisterRequest) (ginx.Result, error) {
	if err := h.svc.Register(c.Request.Context(), &req); err != nil {
		return ginx.Fail(401, err.Error()), nil
	}
	return ginx.OkMsg("注册成功"), nil
}
