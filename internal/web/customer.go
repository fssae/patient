package web

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/util"

	"gitee.com/fssae/ginx"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/customers")
	group.POST("", ginx.WrapBody[domain.CustomerQuery](h.GetList))
	group.GET("/:id", ginx.Wrap(h.GetById))
	group.GET("", ginx.Wrap(h.GetListNameAndID))
	group.POST("/create", ginx.WrapBody[domain.Customer](h.Create))
	group.POST("/:id", ginx.Wrap(h.Update))
	group.PUT("/:id/health-manager", ginx.WrapBody[domain.SetHealthManagerRequest](h.SetHealthManager))
	group.PUT("/:id/bed", ginx.WrapBody[domain.SetBedRequest](h.SetBed))
	group.PUT("/:id/diet-plan", ginx.WrapBody[domain.SetDietPlanRequest](h.SetDietPlan))
	group.PUT("/:id/care-level", ginx.WrapBody[domain.SetCareLevelRequest](h.SetCareLevel))
	group.DELETE("/:id", ginx.Wrap(h.Delete))
}

func (h *CustomerHandler) Create(c *gin.Context, req domain.Customer) (ginx.Result, error) {
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("创建成功"), nil
}

func (h *CustomerHandler) GetListNameAndID(c *gin.Context) (ginx.Result, error) {
	nameIDList, err := h.svc.GetListNameAndID(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Ok("获取成功", nameIDList), nil
}

func (h *CustomerHandler) GetById(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	customer, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		return ginx.Result{}, err
	}
	if customer == nil {
		return ginx.Fail(404, "客户不存在"), nil
	}

	patient := gin.H{
		"patient_id":     customer.ID.Hex(),
		"name":           customer.Name,
		"bed_id":         customer.BedID.Hex(),
		"age":            customer.Age,
		"gender":         customer.Gender,
		"care_level":     "一级护理",
		"health_manager": customer.HealthManager,
		"avatar":         "https://picsum.photos/seed/" + customer.ID.Hex() + "/200/200",
		"diagnosis":      customer.MedicalHistory,
	}

	return ginx.Result{Code: 200, Data: patient}, nil
}

func (h *CustomerHandler) GetList(c *gin.Context, req domain.CustomerQuery) (ginx.Result, error) {
	page := ginx.GetPage(c)
	list, total, err := h.svc.GetList(c.Request.Context(), req, page.Skip, page.Limit)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.OkList("获取成功", list, total), nil
}

func (h *CustomerHandler) Update(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	req, err := util.Validate(&domain.Customer{}, c)
	if err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	if err := h.svc.Update(c.Request.Context(), id, req); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("更新成功"), nil
}

func (h *CustomerHandler) SetHealthManager(c *gin.Context, req domain.SetHealthManagerRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	managerID, err := primitive.ObjectIDFromHex(req.ManagerID)
	if err != nil {
		return ginx.Fail(400, "无效的健康管家ID"), nil
	}
	if err := h.svc.SetHealthManager(c.Request.Context(), id, managerID, req.ManagerName); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("设置成功"), nil
}

func (h *CustomerHandler) SetBed(c *gin.Context, req domain.SetBedRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	bedID, err := primitive.ObjectIDFromHex(req.BedID)
	if err != nil {
		return ginx.Fail(400, "无效的床位ID"), nil
	}
	if err := h.svc.SetBed(c.Request.Context(), id, bedID); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("设置成功"), nil
}

func (h *CustomerHandler) SetDietPlan(c *gin.Context, req domain.SetDietPlanRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	dietPlanID, err := primitive.ObjectIDFromHex(req.DietPlanID)
	if err != nil {
		return ginx.Fail(400, "无效的膳食计划ID"), nil
	}
	if err := h.svc.SetDietPlan(c.Request.Context(), id, dietPlanID); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("设置成功"), nil
}

func (h *CustomerHandler) SetCareLevel(c *gin.Context, req domain.SetCareLevelRequest) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	careLevelID, err := primitive.ObjectIDFromHex(req.CareLevelID)
	if err != nil {
		return ginx.Fail(400, "无效的护理级别ID"), nil
	}
	if err := h.svc.SetCareLevel(c.Request.Context(), id, careLevelID); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("设置成功"), nil
}

func (h *CustomerHandler) Delete(c *gin.Context) (ginx.Result, error) {
	id, ok := ginx.GetId(c)
	if !ok {
		return ginx.Result{}, nil
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		return ginx.Fail(400, err.Error()), nil
	}
	return ginx.OkMsg("删除成功"), nil
}
