package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
)

type DietPlanHandler struct {
	svc *service.DietPlanService
}

func NewDietPlanHandler(svc *service.DietPlanService) *DietPlanHandler {
	return &DietPlanHandler{svc: svc}
}

func (h *DietPlanHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/diet-plans")
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/dietPlan-name-id", ginx.Wrap(h.GetListNameAndID))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.POST("/create", ginx.WrapBody[domain.DietPlan](h.Create))
	group.PUT("/:id", ginx.WrapBody[domain.DietPlan](h.Update))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *DietPlanHandler) Create(c *gin.Context, req domain.DietPlan) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *DietPlanHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	plan, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if plan == nil {
		return ginx.Fail(404, "膳食计划不存在"), nil
	}
	return ginx.Ok("获取成功", plan), nil
}

func (h *DietPlanHandler) GetList(c *gin.Context) (ginx.Result, error) {
	page := ginx.GetPage(c)
	list, total, err := h.svc.GetList(c.Request.Context(), page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *DietPlanHandler) GetListNameAndID(c *gin.Context) (ginx.Result, error) {
	nameIDList, err := h.svc.GetListNameAndID(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", nameIDList), nil
}

func (h *DietPlanHandler) Update(c *gin.Context, req domain.DietPlan) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Update(c.Request.Context(), id, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *DietPlanHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}
