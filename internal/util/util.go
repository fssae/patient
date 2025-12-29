package util

import (
	"classroom-analysis/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetClaims(c *gin.Context) (*domain.PatientClaims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		zap.String("error", "无法解析token")
		return nil, errors.New("无法解析token")
	}
	claim := claims.(domain.PatientClaims)
	return &claim, nil
}
func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return true // 表示有错误发生
	}
	return false
}
