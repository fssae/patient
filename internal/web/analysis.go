package web

import (
	"classroom-analysis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	svc service.AnalysisService
}

func NewAnalysisHandler(svc service.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{svc: svc}
}

func (h *AnalysisHandler) RegisterRoutes(s *gin.RouterGroup) {
	s.GET("/analysis/logs", h.GetLogs)
	s.PUT("/alerts/:id/resolve", h.Resolve)
}

// GetLogs 获取分析与报警日志列表
// @Summary 获取分析与报警日志
// @Description 分页获取存储在 MongoDB 中的报警和会话分析记录
// @Tags 分析管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /analysis/logs [get]
func (h *AnalysisHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "20"), 10, 64)

	logs, total, err := h.svc.GetAnalysisLogs(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": logs,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Resolve 解除告警
// @Summary 解除告警
// @Description 处理前端点击“已处理”按钮的操作
// @Tags 分析管理
// @Accept json
// @Produce json
// @Param id path string true "告警ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/alerts/{id}/resolve [put]
func (h *AnalysisHandler) Resolve(c *gin.Context) {
	id := c.Param("id")
	// TODO: 实现真实的数据库更新逻辑
	// 这里模拟成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alert " + id + " resolved",
	})
}
