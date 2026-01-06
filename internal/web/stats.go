package web

import (
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/ginx"

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
	group.GET("/summary", ginx.Wrap(h.GetSummary))
	group.GET("/trend", ginx.Wrap(h.GetTrend))
}

func (h *StatsHandler) GetSummary(c *gin.Context) (ginx.Result, error) {
	stats, err := h.svc.GetSummary(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{Code: 200, Data: stats}, nil
}

func (h *StatsHandler) GetTrend(c *gin.Context) (ginx.Result, error) {
	trend, err := h.svc.GetTrend(c.Request.Context())
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{Code: 200, Data: trend}, nil
}
