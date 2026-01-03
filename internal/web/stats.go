package web

import (
	"classroom-analysis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	svc service.StatsService
}

func NewStatsHandler(svc service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/api/stats")
	{
		group.GET("/summary", h.GetSummary)
		group.GET("/trend", h.GetTrend)
	}
}

// GetSummary 获取统计概览
// @Summary      获取统计概览
// @Description  获取当前系统的统计概览数据，包括报警统计等
// @Tags         统计分析
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /stats/summary [get]
func (h *StatsHandler) GetSummary(c *gin.Context) {
	stats, err := h.svc.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetTrend 获取趋势分析
// @Summary      获取趋势分析
// @Description  获取最近24小时的报警趋势分析数据
// @Tags         统计分析
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Router       /stats/trend [get]
func (h *StatsHandler) GetTrend(c *gin.Context) {
	trend, err := h.svc.GetTrend(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trend)
}
