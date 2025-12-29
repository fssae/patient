package web

import (
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
	group.POST("/check-out", h.CheckOut)
	group.POST("/outgoing", h.Outgoing)
	group.POST("/return", h.Return)
	group.GET("", h.GetList)
	group.GET("/customer/:customer_id", h.GetByCustomerID)
}

// CheckIn 入住登记
// @Summary      入住登记
// @Description  为客户办理入住登记，自动分配床位并更新客户状态
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "入住信息"  example({"customer_id":"507f1f77bcf86cd799439011","bed_id":"507f1f77bcf86cd799439012","note":"客户入住","created_by":"管理员"})
// @Success      200      {object}  map[string]interface{}  "入住登记成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /records/check-in [post]
func (h *RecordHandler) CheckIn(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		BedID      string `json:"bed_id"`
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

	var bedID primitive.ObjectID
	if req.BedID != "" {
		bedID, err = primitive.ObjectIDFromHex(req.BedID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "无效的床位ID",
			})
			return
		}
	}

	err = h.svc.CheckIn(c.Request.Context(), customerID, bedID, req.Note, req.CreatedBy)
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
