package middleware

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/service"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginJWTMiddlewareBuilder struct {
	paths          []string
	jwtService     *service.JWTService
	blacklistSvc   *service.JWTBlacklistService
	ignorePaths    []string
}

func NewLoginJWTMiddlewareBuilder(jwtService *service.JWTService, blacklistSvc *service.JWTBlacklistService, ignorePaths []string) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		jwtService:   jwtService,
		blacklistSvc: blacklistSvc,
		ignorePaths:  ignorePaths,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// Build 构建JWT中间件
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为忽略路径
		if l.isIgnorePath(c.Request.URL.Path) {
			return
		}

		// 获取Token
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			tokenHeader = c.Request.URL.Query().Get("Authorization")
		}
		if tokenHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code": 401,
				"msg":  "未提供认证token",
			})
			return
		}

		// 解析Token格式
		tokenStr := l.extractToken(tokenHeader)
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code": 401,
				"msg":  "token格式应为 Bearer <Token> 或 <Token>",
			})
			return
		}

		// 检查黑名单
		isBlacklisted, err := l.blacklistSvc.IsBlacklisted(c.Request.Context(), tokenStr)
		if err == nil && isBlacklisted {
			c.AbortWithStatusJSON(401, gin.H{
				"code": 401,
				"msg":  "token已失效，请重新登录",
			})
			return
		}

		// 验证Token - 先尝试作为PatientClaims解析
		patientClaims, err := l.jwtService.ValidatePatientToken(tokenStr)
		if err == nil {
			// 成功解析为PatientClaims
			c.Set("claims", *patientClaims)
			c.Set("token", tokenStr)

			// 检查是否需要刷新Token
			if l.jwtService.ShouldRefreshToken(patientClaims.ExpiresAt) {
				newToken, err := l.jwtService.RefreshPatientToken(patientClaims)
				if err == nil {
					c.Header("X-New-Token", newToken)
				}
			}
			return
		}

		// 尝试作为UserClaims解析
		userClaims, err := l.jwtService.ValidateUserToken(tokenStr)
		if err == nil {
			// 成功解析为UserClaims
			c.Set("claims", *userClaims)
			c.Set("token", tokenStr)

			// 检查是否需要刷新Token
			if l.jwtService.ShouldRefreshToken(userClaims.ExpiresAt) {
				newToken, err := l.jwtService.RefreshUserToken(userClaims)
				if err == nil {
					c.Header("X-New-Token", newToken)
				}
			}
			return
		}

		// Token无效
		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "token无效或已过期",
		})
	}
}

// isIgnorePath 检查路径是否应该被忽略
func (l *LoginJWTMiddlewareBuilder) isIgnorePath(requestPath string) bool {
	// 合并配置文件中的忽略路径和代码中添加的忽略路径
	allPaths := append(l.ignorePaths, l.paths...)
	
	for _, pattern := range allPaths {
		if l.matchPathPattern(requestPath, pattern) {
			return true
		}
	}
	return false
}

// matchPathPattern 匹配路径模式
// 支持:
// 1. 精确匹配: /login
// 2. 通配符匹配: /swagger/* 匹配 /swagger/index.html 等
// 3. 前缀匹配: /api/* 匹配 /api/user/login 等
func (l *LoginJWTMiddlewareBuilder) matchPathPattern(requestPath, pattern string) bool {
	// 精确匹配
	if requestPath == pattern {
		return true
	}

	// 通配符匹配
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		// /swagger/* 应该匹配 /swagger/... 但不匹配 /swagger
		if strings.HasPrefix(requestPath, prefix+"/") {
			return true
		}
		// 也匹配精确的prefix路径
		if requestPath == prefix {
			return true
		}
	}

	// 使用path.Match进行通配符匹配
	matched, err := path.Match(pattern, requestPath)
	if err == nil && matched {
		return true
	}

	return false
}

// extractToken 从Authorization头中提取Token
func (l *LoginJWTMiddlewareBuilder) extractToken(tokenHeader string) string {
	segs := strings.Split(tokenHeader, " ")
	if len(segs) == 2 && strings.ToLower(segs[0]) == "bearer" {
		return segs[1]
	}
	if len(segs) == 1 {
		return segs[0]
	}
	return ""
}

// GetClaimsFromContext 从上下文中获取Claims
func GetClaimsFromContext(c *gin.Context) (interface{}, bool) {
	claims, exists := c.Get("claims")
	return claims, exists
}

// GetPatientClaimsFromContext 从上下文中获取PatientClaims
func GetPatientClaimsFromContext(c *gin.Context) (*domain.PatientClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	patientClaims, ok := claims.(domain.PatientClaims)
	if !ok {
		return nil, false
	}
	return &patientClaims, true
}

// GetUserClaimsFromContext 从上下文中获取UserClaims
func GetUserClaimsFromContext(c *gin.Context) (*domain.UserClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	userClaims, ok := claims.(domain.UserClaims)
	if !ok {
		return nil, false
	}
	return &userClaims, true
}
