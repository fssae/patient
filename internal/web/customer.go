package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *CustomerHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/customers")
	group.GET("", h.GetList)
	group.GET("/:id", h.GetById)
	group.POST("", h.Create)
	group.PUT("/:id", h.Update)
	group.PUT("/:id/health-manager", h.SetHealthManager)
	group.PUT("/:id/bed", h.SetBed)
	group.PUT("/:id/diet-plan", h.SetDietPlan)
	group.PUT("/:id/care-level", h.SetCareLevel)
	group.DELETE("/:id", h.Delete)
}

// Create 创建客户
// @Summary      创建客户
// @Description  创建新的客户记录
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.Customer  true  "客户信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req domain.Customer
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

// GetById 获取客户详情
// @Summary      获取客户详情
// @Description  根据ID获取客户详细信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "客户ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Failure      400  {object}  map[string]interface{}  "无效的ID"
// @Failure      404  {object}  map[string]interface{}  "客户不存在"
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	customer, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "客户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    customer,
	})
}

// GetList 获取客户列表
// @Summary      获取客户列表
// @Description  分页获取客户列表，支持按状态筛选
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "客户状态：入住中/已退住/外出中"
// @Param        skip    query     int     false  "跳过数量"  default(0)
// @Param        limit   query     int     false  "每页数量"  default(20)
// @Success      200     {object}  map[string]interface{}  "获取成功"
// @Router       /customers [get]
func (h *CustomerHandler) GetList(c *gin.Context) {
	status := c.Query("status")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetList(c.Request.Context(), status, skip, limit)
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

// Update 更新客户
func (h *CustomerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.Customer
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

// SetHealthManager 设置健康管家
func (h *CustomerHandler) SetHealthManager(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	var req struct {
		ManagerID   string `json:"manager_id" binding:"required"`
		ManagerName string `json:"manager_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	managerID, err := primitive.ObjectIDFromHex(req.ManagerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的健康管家ID",
		})
		return
	}

	err = h.svc.SetHealthManager(c.Request.Context(), id, managerID, req.ManagerName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "设置成功",
		"success": true,
	})
}

// SetBed 设置床位
func (h *CustomerHandler) SetBed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	var req struct {
		BedID string `json:"bed_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	bedID, err := primitive.ObjectIDFromHex(req.BedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的床位ID",
		})
		return
	}

	err = h.svc.SetBed(c.Request.Context(), id, bedID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "设置成功",
		"success": true,
	})
}

// SetDietPlan 设置膳食计划
func (h *CustomerHandler) SetDietPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	var req struct {
		DietPlanID string `json:"diet_plan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	dietPlanID, err := primitive.ObjectIDFromHex(req.DietPlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的膳食计划ID",
		})
		return
	}

	err = h.svc.SetDietPlan(c.Request.Context(), id, dietPlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "设置成功",
		"success": true,
	})
}

// SetCareLevel 设置护理级别
func (h *CustomerHandler) SetCareLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	var req struct {
		CareLevelID string `json:"care_level_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	careLevelID, err := primitive.ObjectIDFromHex(req.CareLevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的护理级别ID",
		})
		return
	}

	err = h.svc.SetCareLevel(c.Request.Context(), id, careLevelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "设置成功",
		"success": true,
	})
}

// Delete 删除客户
func (h *CustomerHandler) Delete(c *gin.Context) {
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
