package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordHandler struct {
	svc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *RecordHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/records")
	group.POST("/check-in", h.CheckIn)
	//TODO根据房间id查

	group.GET("/check-in-info", h.GetCheckInInfo)
	//TODO
	group.POST("/check-out", h.CheckOut)
	group.GET("/check-out-list", h.GetCheckOutList)
	group.POST("/outgoing", h.Outgoing)
	group.GET("/outgoing-list", h.GetOutgoingList)
	group.POST("/return", h.Return)
	group.GET("", h.GetList)
	group.GET("/customer/:customer_id", h.GetByCustomerID)
	group.POST("/outgoing/update", h.OutgoingUpload)
}

// CheckIn 入住登记
// @Summary      入住登记
// @Description  为客户办理入住登记，创建客户档案，自动分配床位并更新客户状态
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.ElderlyRegisterRequest  true  "入住信息"
// @Success      200      {object}  map[string]interface{}  "入住登记成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /records/check-in [post]
func (h *RecordHandler) CheckIn(c *gin.Context) {
	var req domain.ElderlyRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 简单的参数校验
	if req.Name == "" || req.IDCard == "" || req.BedID == "" || req.HealthLevel == "" || req.NursingLevel == "" || req.DietaryType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "必填参数缺失",
		})
		return
	}

	err := h.svc.CheckIn(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "入住登记成功",
		"success": true,
	})
}

// CheckOut 退住登记
// @Summary      退住登记
// @Description  为客户办理退住登记，释放床位并更新客户状态
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CheckOutRequest  true  "退住信息"
// @Success      200      {object}  map[string]interface{}  "退住登记成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /records/check-out [post]
func (h *RecordHandler) CheckOut(c *gin.Context) {
	var req domain.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	customerID, err := primitive.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	err = h.svc.CheckOut(c.Request.Context(), customerID, req.Note, req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "退住登记成功",
		"success": true,
	})
}

// Outgoing 外出登记
// @Summary      外出登记
// @Description  办理客户外出登记
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.OutgoingRequest  true  "外出信息"
// @Success      200      {object}  map[string]interface{}  "外出登记成功"
// @Router       /records/outgoing [post]
func (h *RecordHandler) Outgoing(c *gin.Context) {
	var req domain.OutgoingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	customerID, err := primitive.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	err = h.svc.Outgoing(c.Request.Context(), customerID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "外出登记成功",
		"success": true,
	})
}

// Return 外出返回
// @Summary      外出返回
// @Description  办理客户外出返回登记
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.ReturnRequest  true  "返回信息"
// @Success      200      {object}  map[string]interface{}  "返回登记成功"
// @Router       /records/return [post]
func (h *RecordHandler) Return(c *gin.Context) {
	var req domain.ReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	RecordsID, err := primitive.ObjectIDFromHex(req.RecordsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	err = h.svc.Return(c.Request.Context(), RecordsID, req.Note, req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "返回登记成功",
		"success": true,
	})
}

// GetList 获取登记记录列表
// @Summary      获取登记记录列表
// @Description  分页获取所有登记记录
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        customer_id  query     string  false  "客户ID"
// @Param        type         query     string  false  "记录类型"
// @Param        skip         query     int     false  "跳过数量"  default(0)
// @Param        limit        query     int     false  "每页数量"  default(20)
// @Success      200          {object}  map[string]interface{}  "获取成功"
// @Router       /records [get]
func (h *RecordHandler) GetList(c *gin.Context) {
	customerIDStr := c.Query("customer_id")
	recordType := c.Query("type")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	var customerID primitive.ObjectID
	if customerIDStr != "" {
		customerID, _ = primitive.ObjectIDFromHex(customerIDStr)
	}

	list, total, err := h.svc.GetList(c.Request.Context(), customerID, recordType, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    list,
		"total":   total,
	})
}

// GetByCustomerID 获取客户的登记记录
// @Summary      获取客户登记记录
// @Description  根据客户ID获取其所有登记记录
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        customer_id  path      string  true  "客户ID"
// @Success      200          {object}  map[string]interface{}  "获取成功"
// @Router       /records/customer/{customer_id} [get]
func (h *RecordHandler) GetByCustomerID(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := primitive.ObjectIDFromHex(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	records, err := h.svc.GetByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    records,
	})
}

// GetCheckInInfo 获取入住登记信息列表
// @Summary      获取入住列表
// @Description  根据条件筛选获取入住登记记录
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        name           query     string  false  "客户姓名"
// @Param        room_number    query     string  false  "房间号"
// @Param        nursing_level  query     string  false  "护理级别"
// @Param        check_in_start query     string  false  "入住开始日期"
// @Param        check_in_end   query     string  false  "入住结束日期"
// @Success      200            {object}  map[string]interface{}  "获取成功"
// @Router       /records/check-in-info [get]
func (h *RecordHandler) GetCheckInInfo(c *gin.Context) {
	name := c.Query("name")
	bedId := c.Query("bed_id")
	nursingLevel := c.Query("nursing_level")
	startDate := c.Query("check_in_start")
	endDate := c.Query("check_in_end")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, err := h.svc.GetCheckInList(c.Request.Context(), name, bedId, nursingLevel, startDate, endDate, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    list,
	})
}

// GetCheckOutList 获取退住登记信息列表
// @Summary      获取退住列表
// @Description  根据条件筛选获取退住登记记录
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        name             query     string  false  "客户姓名"
// @Param        room_number      query     string  false  "房间号"
// @Param        reason           query     string  false  "退住原因"
// @Param        check_out_start  query     string  false  "退住开始日期"
// @Param        check_out_end    query     string  false  "退住结束日期"
// @Success      200              {object}  map[string]interface{}  "获取成功"
// @Router       /records/check-out-list [get]
func (h *RecordHandler) GetCheckOutList(c *gin.Context) {
	name := c.Query("name")
	bedNumber := c.Query("bed_id")
	reason := c.Query("reason")
	startDate := c.Query("check_out_start")
	endDate := c.Query("check_out_end")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, err := h.svc.GetCheckOutList(c.Request.Context(), name, bedNumber, reason, startDate, endDate, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    list,
	})
}

// GetOutgoingList 获取外出登记信息列表
// @Summary      获取外出列表
// @Description  根据条件筛选获取外出登记记录
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        name        query     string  false  "客户姓名"
// @Param        start_date  query     string  false  "开始日期"
// @Param        end_date    query     string  false  "结束日期"
// @Param        status      query     string  false  "状态"
// @Success      200         {object}  map[string]interface{}  "获取成功"
// @Router       /records/outgoing-list [get]
func (h *RecordHandler) GetOutgoingList(c *gin.Context) {
	name := c.Query("name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status") // "全部"/"已外出"/"已返回"

	list, err := h.svc.GetOutgoingList(c.Request.Context(), name, startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    list,
	})
}

func (h *RecordHandler) OutgoingUpload(c *gin.Context) {
	var req domain.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg":     "请求参数错误: " + err.Error(),
			"success": false,
		})
		return
	}

	// 调用服务层处理所有业务逻辑
	if err := h.svc.UpdateOutgoingRecord(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg":     "更新记录失败: " + err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "更新成功",
		"success": true,
	})
}
