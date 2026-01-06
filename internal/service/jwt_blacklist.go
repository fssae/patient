package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type JWTBlacklistService struct {
	redis *redis.Client
}

func NewJWTBlacklistService(redisClient *redis.Client) *JWTBlacklistService {
	return &JWTBlacklistService{
		redis: redisClient,
	}
}

// AddToBlacklist 将Token添加到黑名单
func (s *JWTBlacklistService) AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error {
	key := s.makeKey(token)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token已过期，无需添加到黑名单
		return nil
	}
	return s.redis.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查Token是否在黑名单中
func (s *JWTBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := s.makeKey(token)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// RemoveFromBlacklist 从黑名单中移除Token（通常不需要，Redis会自动过期）
func (s *JWTBlacklistService) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := s.makeKey(token)
	return s.redis.Del(ctx, key).Err()
}

// makeKey 生成Redis键
func (s *JWTBlacklistService) makeKey(token string) string {
	return fmt.Sprintf("jwt:blacklist:%s", token)
}
