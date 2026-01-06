package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordHandler struct {
	svc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{svc: svc}
}

func (h *RecordHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/records")
	group.POST("/check-in", ginx.WrapBody[domain.ElderlyRegisterRequest](h.CheckIn))
	group.GET("/check-in-info", ginx.Wrap(h.GetCheckInInfo))
	group.POST("/check-out", ginx.WrapBody[domain.CheckOutRequest](h.CheckOut))
	group.GET("/check-out-list", ginx.Wrap(h.GetCheckOutList))
	group.POST("/outgoing", ginx.WrapBody[domain.OutgoingRequest](h.Outgoing))
	group.GET("/outgoing-list", ginx.Wrap(h.GetOutgoingList))
	group.POST("/return", ginx.WrapBody[domain.ReturnRequest](h.Return))
	group.GET("", ginx.Wrap(h.GetList))
	group.GET("/customer/:customer_id", ginx.Wrap(h.GetByCustomerID))
	group.POST("/outgoing/update", ginx.WrapBody[domain.UpdateRecordRequest](h.OutgoingUpload))
}

func (h *RecordHandler) CheckIn(c *gin.Context, req domain.ElderlyRegisterRequest) (ginx.Result, error) {
	if req.Name == "" || req.IDCard == "" || req.BedID == "" || req.HealthLevel == "" || req.NursingLevel == "" || req.DietaryType == "" || req.CheckInDate == "" {
		return ginx.Fail(400, "必填参数缺失"), nil
	}
	if err := h.svc.CheckIn(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("入住登记成功"), nil
}

func (h *RecordHandler) CheckOut(c *gin.Context, req domain.CheckOutRequest) (ginx.Result, error) {
	customerID, err := primitive.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return ginx.Fail(400, "无效的客户ID"), nil
	}
	if err := h.svc.CheckOut(c.Request.Context(), customerID, req.Note, req.CreatedBy); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("退住登记成功"), nil
}

func (h *RecordHandler) Outgoing(c *gin.Context, req domain.OutgoingRequest) (ginx.Result, error) {
	customerID, err := primitive.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return ginx.Fail(400, "无效的客户ID"), nil
	}
	if err := h.svc.Outgoing(c.Request.Context(), customerID, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("外出登记成功"), nil
}

func (h *RecordHandler) Return(c *gin.Context, req domain.ReturnRequest) (ginx.Result, error) {
	recordsID, err := primitive.ObjectIDFromHex(req.RecordsID)
	if err != nil {
		return ginx.Fail(400, "无效的记录ID"), nil
	}
	if err := h.svc.Return(c.Request.Context(), recordsID, req.Note, req.CreatedBy); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("返回登记成功"), nil
}

func (h *RecordHandler) GetList(c *gin.Context) (ginx.Result, error) {
	customerIDStr := c.Query("customer_id")
	recordType := c.Query("type")
	page := ginx.GetPage(c)

	var customerID primitive.ObjectID
	if customerIDStr != "" {
		customerID, _ = primitive.ObjectIDFromHex(customerIDStr)
	}

	list, total, err := h.svc.GetList(c.Request.Context(), customerID, recordType, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *RecordHandler) GetByCustomerID(c *gin.Context) (ginx.Result, error) {
	customerIDStr := c.Param("customer_id")
	customerID, err := primitive.ObjectIDFromHex(customerIDStr)
	if err != nil {
		return ginx.Fail(400, "无效的客户ID"), nil
	}
	records, err := h.svc.GetByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", records), nil
}

func (h *RecordHandler) GetCheckInInfo(c *gin.Context) (ginx.Result, error) {
	name := c.Query("name")
	bedId := c.Query("bed_id")
	nursingLevel := c.Query("nursing_level")
	startDate := c.Query("check_in_start")
	endDate := c.Query("check_in_end")
	page := ginx.GetPage(c)

	list, err := h.svc.GetCheckInList(c.Request.Context(), name, bedId, nursingLevel, startDate, endDate, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", list), nil
}

func (h *RecordHandler) GetCheckOutList(c *gin.Context) (ginx.Result, error) {
	name := c.Query("name")
	bedNumber := c.Query("bed_id")
	reason := c.Query("reason")
	startDate := c.Query("check_out_start")
	endDate := c.Query("check_out_end")
	page := ginx.GetPage(c)

	list, err := h.svc.GetCheckOutList(c.Request.Context(), name, bedNumber, reason, startDate, endDate, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", list), nil
}

func (h *RecordHandler) GetOutgoingList(c *gin.Context) (ginx.Result, error) {
	name := c.Query("name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetOutgoingList(c.Request.Context(), name, startDate, endDate, status, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *RecordHandler) OutgoingUpload(c *gin.Context, req domain.UpdateRecordRequest) (ginx.Result, error) {
	if err := h.svc.UpdateOutgoingRecord(c.Request.Context(), &req); err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkMsg("更新成功"), nil
}
