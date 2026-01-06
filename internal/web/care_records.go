package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/util"
	"classroom-analysis/internal/web/ginx"

	"github.com/gin-gonic/gin"
)

type CareRecordHandler struct {
	svc *service.CareRecordService
}

func NewCareRecordHandler(svc *service.CareRecordService) *CareRecordHandler {
	return &CareRecordHandler{svc: svc}
}

func (h *CareRecordHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/care-records")
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.POST("", ginx.WrapBody[domain.CareRecords](h.Create))
	group.POST("/:id/records", ginx.WrapBody[domain.RecordItems](h.AddRecord))
	group.PUT("/:id", ginx.Wrap(h.Update))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *CareRecordHandler) Create(c *gin.Context, req domain.CareRecords) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *CareRecordHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	record, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if record == nil {
		return ginx.Fail(404, "记录不存在"), nil
	}
	return ginx.Ok("获取成功", record), nil
}

func (h *CareRecordHandler) GetList(c *gin.Context) (ginx.Result, error) {
	customerName := c.Query("customer_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	customerID := c.Query("customer_id")
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetList(c.Request.Context(), customerName, page.Skip, page.Limit, startDate, endDate, customerID)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *CareRecordHandler) Update(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	updates, err := util.Validate(domain.CareRecords{}, c)
	if err != nil {
		return ginx.Fail(400, err.Error()), nil
	}

	if err := h.svc.Update(c.Request.Context(), id, updates); err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *CareRecordHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkMsg("删除成功"), nil
}

func (h *CareRecordHandler) AddRecord(c *gin.Context, req domain.RecordItems) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.AddRecord(c.Request.Context(), id, req); err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkMsg("追加成功"), nil
}
