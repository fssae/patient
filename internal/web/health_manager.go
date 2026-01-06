package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HealthManagerHandler struct {
	svc *service.HealthManagerService
}

func NewHealthManagerHandler(svc *service.HealthManagerService) *HealthManagerHandler {
	return &HealthManagerHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *HealthManagerHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/health-managers")
	group.POST("", h.GetList)
	group.GET("/:id", h.GetById)
	group.POST("/create", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create 创建健康管家
// @Summary      创建健康管家
// @Description  创建新的健康管家
// @Tags         健康管家管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.HealthManager  true  "管家信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Router       /health-managers/create [post]
func (h *HealthManagerHandler) Create(c *gin.Context) {
	var req domain.HealthManager
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
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

// GetById 获取健康管家详情
// @Summary      获取详情
// @Description  根据ID获取健康管家详情
// @Tags         健康管家管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "管家ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /health-managers/{id} [get]
func (h *HealthManagerHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	manager, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if manager == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "健康管家不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    manager,
	})
}

// GetList 获取健康管家列表
// @Summary      获取列表
// @Description  根据查询条件分页获取健康管家列表
// @Tags         健康管家管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.HealthManagerQuery  true  "查询条件"
// @Param        skip     query     int                       false  "跳过数量"  default(0)
// @Param        limit    query     int                       false  "每页数量"  default(20)
// @Success      200      {object}  map[string]interface{}           "获取成功"
// @Router       /health-managers [post]
func (h *HealthManagerHandler) GetList(c *gin.Context) {
	var req domain.HealthManagerQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), req, skip, limit)
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

// Update 更新健康管家
// @Summary      更新管家
// @Description  根据ID更新健康管家信息
// @Tags         健康管家管理
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "管家ID"
// @Param        request  body      domain.HealthManager  true  "更新信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Router       /health-managers/{id} [put]
func (h *HealthManagerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.HealthManager
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err = h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
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

// Delete 删除健康管家
// @Summary      删除管家
// @Description  根据ID删除健康管家
// @Tags         健康管家管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "管家ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Router       /health-managers/{id} [delete]
func (h *HealthManagerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
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
