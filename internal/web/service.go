package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceHandler struct {
	svc *service.ServerService
}

func NewServiceHandler(svc *service.ServerService) *ServiceHandler {
	return &ServiceHandler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *ServiceHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/services")
	group.GET("", h.GetServiceList)
	group.GET("/:id", h.GetServiceById)
	group.POST("", h.CreateService)
	group.PUT("/:id", h.UpdateService)
	group.DELETE("/:id", h.DeleteService)
	group.POST("/purchase", h.PurchaseService)
	group.GET("/customer/:customer_id", h.GetCustomerServices)
	group.PUT("/customer-service/:id/end", h.EndService)
}

// CreateService 创建服务项目
// @Summary      创建服务项目
// @Description  添加新的服务项目定义
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.Service  true  "服务项目信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Router       /services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req domain.Service
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.svc.CreateService(c.Request.Context(), &req)
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

// GetServiceById 获取服务项目详情
// @Summary      获取详情
// @Description  根据ID获取服务项目详情
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "服务项目ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /services/{id} [get]
func (h *ServiceHandler) GetServiceById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	service, err := h.svc.GetServiceById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "服务项目不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "获取成功",
		"success": true,
		"data":    service,
	})
}

// GetServiceList 获取服务项目列表
// @Summary      获取服务项目列表
// @Description  分页获取服务项目定义列表
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        category  query     string  false  "类别"
// @Param        status    query     string  false  "状态"
// @Param        skip      query     int     false  "跳过数量"  default(0)
// @Param        limit     query     int     false  "每页数量"  default(20)
// @Success      200       {object}  map[string]interface{}  "获取成功"
// @Router       /services [get]
func (h *ServiceHandler) GetServiceList(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")

	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	list, total, err := h.svc.GetServiceList(c.Request.Context(), category, status, skip, limit)
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

// UpdateService 更新服务项目
// @Summary      更新服务项目
// @Description  根据ID更新服务项目信息
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "服务项目ID"
// @Param        request  body      domain.Service  true  "服务项目信息"
// @Success      200      {object}  map[string]interface{}  "更新成功"
// @Router       /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.Service
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	err = h.svc.UpdateService(c.Request.Context(), id, &req)
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

// DeleteService 删除服务项目
// @Summary      删除服务项目
// @Description  根据ID删除服务项目定义
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "服务项目ID"
// @Success      200  {object}  map[string]interface{}  "删除成功"
// @Router       /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	err = h.svc.DeleteService(c.Request.Context(), id)
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

// PurchaseService 客户购买服务
// @Summary      客户购买服务
// @Description  为客户购买指定的服务项目
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.PurchaseServiceRequest  true  "购买信息"
// @Success      200      {object}  map[string]interface{}  "购买成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /services/purchase [post]
func (h *ServiceHandler) PurchaseService(c *gin.Context) {
	var req domain.PurchaseServiceRequest
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

	serviceID, err := primitive.ObjectIDFromHex(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的服务ID",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的日期格式，请使用 YYYY-MM-DD",
		})
		return
	}

	err = h.svc.PurchaseService(c.Request.Context(), customerID, serviceID, startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "购买成功",
		"success": true,
	})
}

// GetCustomerServices 获取客户购买的服务列表
// @Summary      获取客户服务列表
// @Description  获取指定客户已购买的服务列表
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        customer_id  path      string  true  "客户ID"
// @Success      200          {object}  map[string]interface{}  "获取成功"
// @Router       /services/customer/{customer_id} [get]
func (h *ServiceHandler) GetCustomerServices(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := primitive.ObjectIDFromHex(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的客户ID",
		})
		return
	}

	services, err := h.svc.GetCustomerServices(c.Request.Context(), customerID)
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
		"data":    services,
	})
}

// EndService 结束客户服务
// @Summary      结束客户服务
// @Description  根据购买记录ID手动结束一项服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "购买记录ID"
// @Param        request  body      domain.EndServiceRequest  true  "结束信息"
// @Success      200      {object}  map[string]interface{}  "结束成功"
// @Router       /services/customer-service/{id}/end [put]
func (h *ServiceHandler) EndService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	var req domain.EndServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的日期格式，请使用 YYYY-MM-DD",
		})
		return
	}

	err = h.svc.EndService(c.Request.Context(), id, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "结束成功",
		"success": true,
	})
}
