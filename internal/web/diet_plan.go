package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DietPlanHandler struct {
	svc *service.DietPlanService
}

func NewDietPlanHandler(svc *service.DietPlanService) *DietPlanHandler {
	return &DietPlanHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *DietPlanHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/diet-plans")
	group.GET("", h.GetList)
	group.GET("/dietPlan-name-id", h.GetListNameAndID)
	group.GET("/:id", h.GetById)
	group.POST("/create", h.Create)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create 创建膳食计划
// @Summary      创建膳食计划
// @Description  创建新的膳食计划
// @Tags         膳食计划管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.DietPlan  true  "计划信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Router       /diet-plans/create [post]
func (h *DietPlanHandler) GetListNameAndID(c *gin.Context) {
	nameIDList, err := h.svc.GetListNameAndID(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取客户名称和ID失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": nameIDList,
	})
}

func (h *DietPlanHandler) Create(c *gin.Context) {
	var req domain.DietPlan
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

// GetById 获取膳食计划详情
// @Summary      获取详情
// @Description  根据ID获取膳食计划详情
// @Tags         膳食计划管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "计划ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /diet-plans/{id} [get]
func (h *DietPlanHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	plan, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "膳食计划不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    plan,
	})
}

// GetList 获取膳食计划列表
// @Summary      获取列表
// @Description  分页获取膳食计划列表
// @Tags         膳食计划管理
// @Accept       json
// @Produce      json
// @Param        skip   query     int  false  "跳过数量"  default(0)
// @Param        limit  query     int  false  "每页数量"  default(20)
// @Success      200    {object}  map[string]interface{}  "获取成功"
// @Router       /diet-plans [get]
func (h *DietPlanHandler) GetList(c *gin.Context) {
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), skip, limit)
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

// Update 更新膳食计划
// @Summary      更新计划
// @Description  根据ID更新膳食计划
// @Tags         膳食计划管理
// @Accept       json
// @Produce      json
// @Param        id       path      string           true  "计划ID"
// @Param        request  body      domain.DietPlan  true  "计划信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Router       /diet-plans/{id} [put]
func (h *DietPlanHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.DietPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
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

// Delete 删除膳食计划
// @Summary      删除计划
// @Description  根据ID删除膳食计划
// @Tags         膳食计划管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "计划ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Router       /diet-plans/{id} [delete]
func (h *DietPlanHandler) Delete(c *gin.Context) {
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
