package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/util"
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

// RegisterRoutes 注册路由
func (h *BedHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/beds")
	group.GET("", h.GetList)
	group.GET("/options/beds", h.GetBedOptions)
	group.GET("/options/rooms", h.GetRoomOptions)
	group.GET("/:id", h.GetById)
	group.GET("/room/:room_id", h.GetByRoomID)
	group.POST("/create", h.Create)
	group.PUT("/:id", h.Update)
	group.PUT("/:id/release", h.Release)
	group.DELETE("/:id", h.DeleteBed)
}

// DeleteBed 删除床位
func (h *BedHandler) DeleteBed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "删除成功",
		"success": true,
	})
}

// Create 创建床位
func (h *BedHandler) Create(c *gin.Context) {
	var req domain.CreateBedRequest
	if util.HandleError(c, c.ShouldBindJSON(&req)) {
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
	roomNumber := c.Query("room_number")
	bedNumber := c.Query("bed_number")
	status := c.Query("status")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), roomNumber, bedNumber, status, skip, limit)
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

// Update 更新床位信息
// @Summary      更新床位
// @Description  更新床位信息
// @Tags         床位管理
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "床位ID"
// @Param        request  body      domain.Bed  true  "床位信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /beds/{id} [put]
func (h *BedHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的床位ID",
		})
		return
	}

	var req domain.Bed
	if util.HandleError(c, c.ShouldBindJSON(&req)) {
		return
	}

	err = h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "更新成功",
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

// GetBedOptions 获取床位选项列表
func (h *BedHandler) GetBedOptions(c *gin.Context) {
	options, err := h.svc.GetBedOptions(c.Request.Context())
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
		"data":    options,
	})
}

// GetRoomOptions 获取房间选项列表
func (h *BedHandler) GetRoomOptions(c *gin.Context) {
	options, err := h.svc.GetRoomOptions(c.Request.Context())
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
		"data":    options,
	})
}
