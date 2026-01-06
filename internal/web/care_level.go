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

type CareLevelHandler struct {
	svc *service.CareLevelService
}

func NewCareLevelHandler(svc *service.CareLevelService) *CareLevelHandler {
	return &CareLevelHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *CareLevelHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/care-levels")
	group.GET("", h.GetList)
	group.GET("/menu", h.GetMenu)
	group.GET("/:id", h.GetById)
	group.POST("/create", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create 创建护理级别
// @Summary      创建护理级别
// @Description  创建新的护理级别
// @Tags         护理级别管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CareLevel  true  "护理级别信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /care-levels [post]
func (h *CareLevelHandler) Create(c *gin.Context) {
	var req domain.CareLevel
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

// GetById 获取护理级别详情
// @Summary      获取护理级别详情
// @Description  根据ID获取护理级别详细信息
// @Tags         护理级别管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "护理级别ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Failure      400  {object}  map[string]interface{}  "无效的ID"
// @Failure      404  {object}  map[string]interface{}  "护理级别不存在"
// @Router       /care-levels/{id} [get]
func (h *CareLevelHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	level, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if level == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "护理级别不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    level,
	})
}

// GetList 获取护理级别列表
// @Summary      获取护理级别列表
// @Description  分页获取护理级别列表
// @Tags         护理级别管理
// @Accept       json
// @Produce      json
// @Param        skip    query     int     false  "跳过数量"  default(0)
// @Param        limit   query     int     false  "每页数量"  default(20)
// @Success      200     {object}  map[string]interface{}  "获取成功"
// @Router       /care-levels [get]
func (h *CareLevelHandler) GetList(c *gin.Context) {
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")
	level, _ := strconv.ParseInt(c.DefaultQuery("level", "0"), 10, 64)

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), level, skip, limit)
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

// Update 更新护理级别
// @Summary      更新护理级别
// @Description  更新护理级别信息
// @Tags         护理级别管理
// @Accept       json
// @Produce      json
// @Param        id       path      string            true  "护理级别ID"
// @Param        request  body      domain.CareLevel  true  "护理级别信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /care-levels/{id} [put]
func (h *CareLevelHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}
	// 验证请求参数，组合业务字段
	req, err := util.Validate(domain.CareLevel{}, c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err = h.svc.Update(c.Request.Context(), id, req)
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

// Delete 删除护理级别
// @Summary      删除护理级别
// @Description  删除指定的护理级别
// @Tags         护理级别管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "护理级别ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Failure      400  {object}  map[string]interface{}  "无效的ID"
// @Router       /care-levels/{id} [delete]
func (h *CareLevelHandler) Delete(c *gin.Context) {
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

// GetMenu 获取护理级别菜单
// @Summary      获取护理级别菜单
// @Description  用于下拉列表选择的护理级别数据
// @Tags         护理级别管理
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /care-levels/menu [get]
func (h *CareLevelHandler) GetMenu(c *gin.Context) {
	list, err := h.svc.GetMenu(c.Request.Context())
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
