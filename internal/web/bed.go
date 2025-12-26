package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BedHandler struct {
	svc *service.BedService
}

func NewBedHandler(svc *service.BedService) *BedHandler {
	return &BedHandler{
		svc: svc,
	}
}

// Create 创建床位
func (h *BedHandler) Create(c *gin.Context) {
	var req domain.Bed
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "创建成功",
		"success": true,
	})
}

// GetById 获取床位详情
func (h *BedHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	bed, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if bed == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "床位不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    bed,
	})
}

// GetByRoomID 根据房间ID获取床位列表
func (h *BedHandler) GetByRoomID(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := primitive.ObjectIDFromHex(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的房间ID",
		})
		return
	}

	beds, err := h.svc.GetByRoomID(c.Request.Context(), roomID)
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
		"data":    beds,
	})
}

// GetList 获取床位列表
func (h *BedHandler) GetList(c *gin.Context) {
	roomIDStr := c.Query("room_id")
	status := c.Query("status")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	var roomID primitive.ObjectID
	if roomIDStr != "" {
		roomID, _ = primitive.ObjectIDFromHex(roomIDStr)
	}

	list, total, err := h.svc.GetList(c.Request.Context(), roomID, status, skip, limit)
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

// AssignToCustomer 分配床位给客户
func (h *BedHandler) AssignToCustomer(c *gin.Context) {
	bedIDStr := c.Param("id")
	bedID, err := primitive.ObjectIDFromHex(bedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的床位ID",
		})
		return
	}

	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
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

	err = h.svc.AssignToCustomer(c.Request.Context(), bedID, customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "分配成功",
		"success": true,
	})
}

// Release 释放床位
func (h *BedHandler) Release(c *gin.Context) {
	bedIDStr := c.Param("id")
	bedID, err := primitive.ObjectIDFromHex(bedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的床位ID",
		})
		return
	}

	err = h.svc.Release(c.Request.Context(), bedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "释放成功",
		"success": true,
	})
}

// RegisterRoutes 注册路由
func (h *BedHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/api/beds")
	group.GET("", h.GetList)
	group.GET("/:id", h.GetById)
	group.GET("/room/:room_id", h.GetByRoomID)
	group.POST("", h.Create)
	group.PUT("/:id/assign", h.AssignToCustomer)
	group.PUT("/:id/release", h.Release)
}

