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
	group.GET("/check-in-info", h.GetCheckInInfo)
	group.POST("/check-out", h.CheckOut)
	group.GET("/check-out-list", h.GetCheckOutList)
	group.POST("/outgoing", h.Outgoing)
	group.GET("/outgoing-list", h.GetOutgoingList)
	group.POST("/return", h.Return)
	group.GET("", h.GetList)
	group.GET("/customer/:customer_id", h.GetByCustomerID)
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
// @Param        request  body      object  true  "退住信息"  example({"customer_id":"507f1f77bcf86cd799439011","note":"客户退住","created_by":"管理员"})
// @Success      200      {object}  map[string]interface{}  "退住登记成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /records/check-out [post]
func (h *RecordHandler) CheckOut(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		Note       string `json:"note"`
		CreatedBy  string `json:"created_by" binding:"required"`
	}
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
func (h *RecordHandler) Outgoing(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		Note       string `json:"note"`
		CreatedBy  string `json:"created_by" binding:"required"`
	}
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

	err = h.svc.Outgoing(c.Request.Context(), customerID, req.Note, req.CreatedBy)
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
func (h *RecordHandler) Return(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		Note       string `json:"note"`
		CreatedBy  string `json:"created_by" binding:"required"`
	}
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

	err = h.svc.Return(c.Request.Context(), customerID, req.Note, req.CreatedBy)
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
func (h *RecordHandler) GetCheckInInfo(c *gin.Context) {
	name := c.Query("name")
	roomNumber := c.Query("room_number")
	nursingLevel := c.Query("nursing_level")
	startDate := c.Query("check_in_start")
	endDate := c.Query("check_in_end")

	list, err := h.svc.GetCheckInList(c.Request.Context(), name, roomNumber, nursingLevel, startDate, endDate)
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
func (h *RecordHandler) GetCheckOutList(c *gin.Context) {
	name := c.Query("name")
	roomNumber := c.Query("room_number")
	reason := c.Query("reason")
	startDate := c.Query("check_out_start")
	endDate := c.Query("check_out_end")

	list, err := h.svc.GetCheckOutList(c.Request.Context(), name, roomNumber, reason, startDate, endDate)
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
