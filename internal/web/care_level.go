package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/util"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
)

type CareLevelHandler struct {
	svc *service.CareLevelService
}

func NewCareLevelHandler(svc *service.CareLevelService) *CareLevelHandler {
	return &CareLevelHandler{svc: svc}
}

func (h *CareLevelHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/care-levels")
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/menu", ginx.Wrap(h.GetMenu))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.POST("/create", ginx.WrapBody[domain.CareLevel](h.Create))
	group.PUT("/:id", ginx.Wrap(h.Update))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *CareLevelHandler) Create(c *gin.Context, req domain.CareLevel) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *CareLevelHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	level, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if level == nil {
		return ginx.Fail(404, "护理级别不存在"), nil
	}
	return ginx.Ok("获取成功", level), nil
}

func (h *CareLevelHandler) GetList(c *gin.Context) (ginx.Result, error) {
	level := util.ParseInt64Query(c, "level", 0)
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetList(c.Request.Context(), level, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *CareLevelHandler) Update(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	req, err := util.Validate(domain.CareLevel{}, c)
	if err != nil {
		return ginx.Fail(400, "请求参数错误: "+err.Error()), nil
	}

	if err := h.svc.Update(c.Request.Context(), id, req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *CareLevelHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}

func (h *CareLevelHandler) GetMenu(c *gin.Context) (ginx.Result, error) {
	list, err := h.svc.GetMenu(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", list), nil
}
