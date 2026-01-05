package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知推送处理器
type NotificationHandler struct{}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// RegisterRoutes 注册通知相关路由
func (h *NotificationHandler) RegisterRoutes(server gin.IRouter) {
	group := server.Group("/api/notifications")
	{
		group.POST("/call", h.CallStaff)
	}
}

// CallRequest 呼叫请求
type CallRequest struct {
	AlertID    string `json:"alert_id" binding:"required"`
	TargetRole string `json:"target_role" binding:"required"` // doctor, nurse, manager
}

// CallStaff 呼叫医护人员
// @Summary      呼叫医护人员
// @Description  触发对指定角色的通知推送
// @Tags         通知推送
// @Accept       json
// @Produce      json
// @Param        request  body      CallRequest  true  "呼叫请求"
// @Success      200      {object}  map[string]interface{}  "呼叫成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /notifications/call [post]
func (h *NotificationHandler) CallStaff(c *gin.Context) {
	var req CallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "BAD_REQUEST",
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 根据角色返回模拟响应
	var targetName string
	switch req.TargetRole {
	case "doctor":
		targetName = "值班医生"
	case "nurse":
		targetName = "值班护士"
	case "manager":
		targetName = "护士长"
	default:
		targetName = "医护人员"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已通知" + targetName + "，预计2分钟内响应",
	})
}
