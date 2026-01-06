package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
)

type HealthManagerHandler struct {
	svc *service.HealthManagerService
}

func NewHealthManagerHandler(svc *service.HealthManagerService) *HealthManagerHandler {
	return &HealthManagerHandler{svc: svc}
}

func (h *HealthManagerHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/health-managers")
	group.POST("", ginx.WrapBody[domain.HealthManagerQuery](h.GetList))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.POST("/create", ginx.WrapBody[domain.HealthManager](h.Create))
	group.PUT("/:id", ginx.WrapBody[domain.HealthManager](h.Update))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *HealthManagerHandler) Create(c *gin.Context, req domain.HealthManager) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *HealthManagerHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	manager, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if manager == nil {
		return ginx.Fail(404, "健康管家不存在"), nil
	}
	return ginx.Ok("获取成功", manager), nil
}

func (h *HealthManagerHandler) GetList(c *gin.Context, req domain.HealthManagerQuery) (ginx.Result, error) {
	page := ginx.GetPage(c)
	list, total, err := h.svc.GetList(c.Request.Context(), req, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *HealthManagerHandler) Update(c *gin.Context, req domain.HealthManager) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *HealthManagerHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}
