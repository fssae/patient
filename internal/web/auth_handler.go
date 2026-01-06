package web

import (
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	jwtService       *service.JWTService
	blacklistService *service.JWTBlacklistService
}

func NewAuthHandler(jwtService *service.JWTService, blacklistService *service.JWTBlacklistService) *AuthHandler {
	return &AuthHandler{
		jwtService:       jwtService,
		blacklistService: blacklistService,
	}
}

// RegisterRoutes 注册认证相关路由
func (h *AuthHandler) RegisterRoutes(server gin.IRouter) {
	server.POST("/api/auth/logout", h.Logout)
}

// Logout 用户登出
// @Summary      用户登出
// @Description  将当前Token加入黑名单，使其失效
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "登出成功"
// @Failure      401  {object}  map[string]interface{}  "未提供认证token"
// @Failure      500  {object}  map[string]interface{}  "登出失败"
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取Token
	tokenStr, exists := c.Get("token")
	if !exists {
		// 尝试从Header获取
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供认证token",
			})
			return
		}
		tokenStr = h.extractToken(tokenHeader)
	}

	token, ok := tokenStr.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "token格式错误",
		})
		return
	}

	// 获取Claims以确定Token过期时间
	claims, exists := middleware.GetClaimsFromContext(c)
	if !exists {
		// 如果上下文中没有claims，尝试重新解析
		patientClaims, err := h.jwtService.ValidatePatientToken(token)
		if err == nil {
			claims = patientClaims
		} else {
			userClaims, err := h.jwtService.ValidateUserToken(token)
			if err == nil {
				claims = userClaims
			}
		}
	}

	// 获取过期时间
	var expiresAt *jwt.NumericDate
	switch v := claims.(type) {
	case *jwt.RegisteredClaims:
		expiresAt = v.ExpiresAt
	default:
		// 尝试从接口获取
		if rc, ok := claims.(interface{ GetExpirationTime() (*jwt.NumericDate, error) }); ok {
			expiresAt, _ = rc.GetExpirationTime()
		}
	}

	if expiresAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无法获取token过期时间",
		})
		return
	}

	// 将Token添加到黑名单
	err := h.blacklistService.AddToBlacklist(c.Request.Context(), token, expiresAt.Time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "登出失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "登出成功",
		"success": true,
	})
}

// extractToken 从Authorization头中提取Token
func (h *AuthHandler) extractToken(tokenHeader string) string {
	segs := strings.Split(tokenHeader, " ")
	if len(segs) == 2 && strings.ToLower(segs[0]) == "bearer" {
		return segs[1]
	}
	if len(segs) == 1 {
		return segs[0]
	}
	return ""
}
