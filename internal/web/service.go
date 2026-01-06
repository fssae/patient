package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/ginx"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceHandler struct {
	svc *service.ServerService
}

func NewServiceHandler(svc *service.ServerService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

func (h *ServiceHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/services")
	group.GET("", ginx.Wrap(h.GetServiceList))
	group.GET("/customer-service-list", ginx.Wrap(h.GetCustomerServiceList))
	group.GET("/:id", ginx.Wrap(h.GetServiceById))
	group.POST("", ginx.WrapBody[domain.Service](h.CreateService))
	group.PUT("/:id", ginx.WrapBody[domain.Service](h.UpdateService))
	group.DELETE("/:id", ginx.Wrap(h.DeleteService))
	group.POST("/purchase", ginx.WrapBody[domain.PurchaseServiceRequest](h.PurchaseService))
	group.GET("/customer/:customer_id", ginx.Wrap(h.GetCustomerServices))
	group.PUT("/customer_service/:id/end", ginx.WrapBody[domain.EndServiceRequest](h.EndService))
	group.PUT("/customer_service/:id/end_date", ginx.WrapBody[domain.EndServiceRequest](h.UpdateCustomerServiceEndDate))
	group.PUT("/customer_service/:id/cancel", ginx.Wrap(h.CancelCustomerService))
	group.GET("/customer-services_by_group", ginx.Wrap(h.GetCustomerService))
}

func (h *ServiceHandler) CreateService(c *gin.Context, req domain.Service) (ginx.Result, error) {
	if err := h.svc.CreateService(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *ServiceHandler) GetServiceById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	service, err := h.svc.GetServiceById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if service == nil {
		return ginx.Fail(404, "服务项目不存在"), nil
	}
	return ginx.Ok("获取成功", service), nil
}

func (h *ServiceHandler) GetServiceList(c *gin.Context) (ginx.Result, error) {
	category := c.Query("category")
	status := c.Query("status")
	page := ginx.GetPage(c)

	list, total, err := h.svc.GetServiceList(c.Request.Context(), category, status, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *ServiceHandler) UpdateService(c *gin.Context, req domain.Service) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.UpdateService(c.Request.Context(), id, &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *ServiceHandler) DeleteService(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.DeleteService(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}

func (h *ServiceHandler) PurchaseService(c *gin.Context, req domain.PurchaseServiceRequest) (ginx.Result, error) {
	customerID, err := primitive.ObjectIDFromHex(req.CustomerID)
	if err != nil {
		return ginx.Fail(400, "无效的客户ID"), nil
	}
	serviceID, err := primitive.ObjectIDFromHex(req.ServiceID)
	if err != nil {
		return ginx.Fail(400, "无效的服务ID"), nil
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return ginx.Fail(400, "无效的日期格式，请使用 YYYY-MM-DD"), nil
	}
	if err := h.svc.PurchaseService(c.Request.Context(), customerID, serviceID, startDate); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("购买成功"), nil
}

func (h *ServiceHandler) GetCustomerServices(c *gin.Context) (ginx.Result, error) {
	customerIDStr := c.Param("customer_id")
	customerID, err := primitive.ObjectIDFromHex(customerIDStr)
	if err != nil {
		return ginx.Fail(400, "无效的客户ID"), nil
	}
	services, err := h.svc.GetCustomerServices(c.Request.Context(), customerID)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", services), nil
}

func (h *ServiceHandler) EndService(c *gin.Context, req domain.EndServiceRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return ginx.Fail(400, "无效的日期格式，请使用 YYYY-MM-DD"), nil
	}
	if err := h.svc.EndService(c.Request.Context(), id, endDate); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("结束成功"), nil
}

func (h *ServiceHandler) UpdateCustomerServiceEndDate(c *gin.Context, req domain.EndServiceRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return ginx.Fail(400, "无效的日期格式，请使用 YYYY-MM-DD"), nil
	}
	if err := h.svc.UpdateCustomerServiceEndDate(c.Request.Context(), id, endDate); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("修改成功"), nil
}

func (h *ServiceHandler) CancelCustomerService(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.CancelCustomerService(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("取消成功"), nil
}

func (h *ServiceHandler) GetCustomerServiceList(c *gin.Context) (ginx.Result, error) {
	customerIDStr := c.Query("customer_id")
	serviceIDStr := c.Query("service_id")
	status := c.Query("status")
	page := ginx.GetPage(c)

	customerID, _ := primitive.ObjectIDFromHex(customerIDStr)
	serviceID, _ := primitive.ObjectIDFromHex(serviceIDStr)

	list, total, err := h.svc.GetCustomerServiceList(c.Request.Context(), customerID, serviceID, status, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *ServiceHandler) GetCustomerService(c *gin.Context) (ginx.Result, error) {
	customerIDStr := c.Query("customer_id")
	serviceName := c.Query("service_name")
	status := c.Query("status")
	page := ginx.GetPage(c)

	customerID, _ := primitive.ObjectIDFromHex(customerIDStr)

	list, total, err := h.svc.GetCustomerServiceByGroup(c.Request.Context(), customerID, serviceName, status, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}
