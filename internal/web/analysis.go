package web

import (
	"classroom-analysis/internal/service"
	"strconv"

	"gitee.com/fssae/ginx"
	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	svc service.AnalysisService
}

func NewAnalysisHandler(svc service.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{svc: svc}
}

func (h *AnalysisHandler) RegisterRoutes(s *gin.RouterGroup) {
	s.GET("/analysis/logs", ginx.Wrap(h.GetLogs))
	s.PUT("/alerts/:id/resolve", ginx.Wrap(h.Resolve))
}

func (h *AnalysisHandler) GetLogs(c *gin.Context) (ginx.Result, error) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "20"), 10, 64)

	logs, total, err := h.svc.GetAnalysisLogs(c.Request.Context(), page, size)
	if err != nil {
		return ginx.Result{}, err
	}

	return ginx.Result{
		Code: 200,
		Data: gin.H{
			"data": logs,
			"meta": gin.H{
				"total": total,
				"page":  page,
				"size":  size,
			},
		},
	}, nil
}

func (h *AnalysisHandler) Resolve(c *gin.Context) (ginx.Result, error) {
	id := c.Param("id")
	return ginx.Result{
		Code: 200,
		Data: gin.H{
			"status":  "success",
			"message": "Alert " + id + " resolved",
		},
	}, nil
}
