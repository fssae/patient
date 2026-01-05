package web

import (
	"classroom-analysis/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CompatHandler 前端兼容接口处理器
// 单独开文件处理前端规范的接口路径
type CompatHandler struct {
	customerSvc *service.CustomerService
	statsSvc    service.StatsService
}

// NewCompatHandler 创建兼容处理器
func NewCompatHandler(customerSvc *service.CustomerService, statsSvc service.StatsService) *CompatHandler {
	return &CompatHandler{
		customerSvc: customerSvc,
		statsSvc:    statsSvc,
	}
}

// RegisterRoutes 注册兼容路由
func (h *CompatHandler) RegisterRoutes(server gin.IRouter) {
	// 患者详情接口（前端规范路径）
	server.GET("/api/patients/:patient_id", h.GetPatient)
	// 趋势接口（前端规范路径）
	server.GET("/api/stats/trends", h.GetTrends)
	// 值班人员接口
	server.GET("/api/staff/on-duty", h.GetOnDutyStaff)
}

// GetPatient 获取患者详情（前端规范格式）
// @Summary      获取患者详情
// @Description  根据患者ID获取详细信息
// @Tags         患者信息
// @Produce      json
// @Param        patient_id  path      string  true  "患者ID"
// @Success      200         {object}  map[string]interface{}  "获取成功"
// @Failure      400         {object}  map[string]interface{}  "无效的ID"
// @Failure      404         {object}  map[string]interface{}  "患者不存在"
// @Router       /patients/{patient_id} [get]
func (h *CompatHandler) GetPatient(c *gin.Context) {
	patientID := c.Param("patient_id")
	id, err := primitive.ObjectIDFromHex(patientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "BAD_REQUEST",
			"message": "无效的患者ID",
		})
		return
	}

	customer, err := h.customerSvc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": err.Error(),
		})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "PATIENT_NOT_FOUND",
			"message": "患者不存在",
		})
		return
	}

	// 转换为前端规范的 Patient 格式
	patient := gin.H{
		"patient_id":     customer.ID.Hex(),
		"name":           customer.Name,
		"bed_id":         customer.BedID.Hex(),
		"age":            customer.Age,
		"gender":         customer.Gender,
		"care_level":     "一级护理", // TODO: 从 CareLevelID 获取真实名称
		"health_manager": customer.HealthManager,
		"avatar":         "https://picsum.photos/seed/" + customer.ID.Hex() + "/200/200",
		"diagnosis":      customer.MedicalHistory,
	}

	c.JSON(http.StatusOK, patient)
}

// TrendDataCompat 趋势数据（前端规范格式）
type TrendDataCompat struct {
	Time  string `json:"time"`
	Count int    `json:"count"`
}

// GetTrends 获取24小时告警趋势（前端规范格式）
// @Summary      获取24小时告警趋势
// @Description  获取最近24小时的告警趋势数据
// @Tags         统计分析
// @Produce      json
// @Success      200  {array}  TrendDataCompat  "获取成功"
// @Router       /stats/trends [get]
func (h *CompatHandler) GetTrends(c *gin.Context) {
	// 调用原有服务获取数据
	trends, err := h.statsSvc.GetTrend(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": err.Error(),
		})
		return
	}

	// 转换字段名称：name -> time
	result := make([]TrendDataCompat, len(trends))
	for i, t := range trends {
		result[i] = TrendDataCompat{
			Time:  t.Name,
			Count: t.Count,
		}
	}

	c.JSON(http.StatusOK, result)
}

// DutyStaffInfo 值班人员信息
type DutyStaffInfo struct {
	OnDutyCount int      `json:"onDutyCount"`
	Doctors     []string `json:"doctors"`
	Nurses      []string `json:"nurses"`
	LastUpdated int64    `json:"lastUpdated"`
}

// GetOnDutyStaff 获取当前值班人员
// @Summary      获取值班人员
// @Description  获取当前在岗的医护人员列表
// @Tags         值班管理
// @Produce      json
// @Success      200  {object}  DutyStaffInfo  "获取成功"
// @Router       /staff/on-duty [get]
func (h *CompatHandler) GetOnDutyStaff(c *gin.Context) {
	// 返回模拟数据
	info := DutyStaffInfo{
		OnDutyCount: 12,
		Doctors:     []string{"张医生", "李医生"},
		Nurses:      []string{"陈护士", "王护士", "刘护士"},
		LastUpdated: time.Now().Unix(),
	}

	c.JSON(http.StatusOK, info)
}
