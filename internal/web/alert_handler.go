package web

import (
	"classroom-analysis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AlertHandler 告警管理处理器
type AlertHandler struct {
	analysisSvc service.AnalysisService
}

// NewAlertHandler 创建告警处理器
func NewAlertHandler(analysisSvc service.AnalysisService) *AlertHandler {
	return &AlertHandler{analysisSvc: analysisSvc}
}

// RegisterRoutes 注册告警管理路由
func (h *AlertHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/alerts")
	{
		group.PUT("/:id/resolve", h.ResolveAlert)
	}
}

// ResolveAlert 解除告警
// @Summary      解除告警
// @Description  将指定告警标记为已处理
// @Tags         告警管理
// @Param        id   path  string  true  "告警ID"
// @Success      204  "解除成功"
// @Failure      404  {object}  map[string]interface{}  "告警不存在"
// @Failure      500  {object}  map[string]interface{}  "服务器错误"
// @Router       /alerts/{id}/resolve [put]
func (h *AlertHandler) ResolveAlert(c *gin.Context) {
	alertID := c.Param("id")

	err := h.analysisSvc.ResolveAlert(c.Request.Context(), alertID)
	if err != nil {
		if err.Error() == "告警不存在" || err.Error() == "无效的告警ID" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "ALERT_NOT_FOUND",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
