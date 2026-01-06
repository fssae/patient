package ioc

import (
    "classroom-analysis/internal/service"
    "classroom-analysis/internal/web/middleware"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "github.com/spf13/viper"
)

func InitJWTMiddleware(redisClient *redis.Client) gin.HandlerFunc {
    // 初始化JWT服务
    jwtService := service.NewJWTService()
    
    // 初始化JWT黑名单服务
    blacklistService := service.NewJWTBlacklistService(redisClient)
    
    // 从配置文件读取忽略路径
    ignorePaths := viper.GetStringSlice("jwt.ignore_paths")
    if len(ignorePaths) == 0 {
        // 默认忽略路径
        ignorePaths = []string{
            "/login",
            "/register",
            "/api/user/login",
            "/api/user/register",
            "/health",
            "/swagger/*",
            "/metrics",
            "/debug/*",
        }
    }
    
    // 创建JWT中间件
    return middleware.NewLoginJWTMiddlewareBuilder(
        jwtService,
        blacklistService,
        ignorePaths,
    ).Build()
}
