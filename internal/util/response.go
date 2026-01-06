package util

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ 响应函数 ============

// Success 成功响应
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     msg,
		"success": true,
		"data":    data,
	})
}

// SuccessWithTotal 带总数的成功响应（用于列表查询）
func SuccessWithTotal(c *gin.Context, msg string, data interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     msg,
		"success": true,
		"data":    data,
		"total":   total,
	})
}

// SuccessMsg 仅消息的成功响应
func SuccessMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     msg,
		"success": true,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// FailServer 服务器错误响应
func FailServer(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"msg":  err.Error(),
	})
}

// ============ 请求解析函数 ============

// ParseObjectID 解析URL路径中的ObjectID参数，失败时自动返回错误响应
func ParseObjectID(c *gin.Context, paramName string) (primitive.ObjectID, bool) {
	idStr := c.Param(paramName)
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		Fail(c, 400, "无效的ID")
		return primitive.NilObjectID, false
	}
	return id, true
}

// Pagination 分页参数
type Pagination struct {
	Skip  int64
	Limit int64
}

// ParsePagination 解析分页参数（skip和limit）
func ParsePagination(c *gin.Context) Pagination {
	skipStr := c.DefaultQuery("skip", "0")
	limitStr := c.DefaultQuery("limit", "20")
	skip, _ := strconv.ParseInt(skipStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	return Pagination{Skip: skip, Limit: limit}
}

// ParseInt64Query 解析int64类型的查询参数
func ParseInt64Query(c *gin.Context, key string, defaultVal int64) int64 {
	str := c.Query(key)
	if str == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseIntQuery 解析int类型的查询参数
func ParseIntQuery(c *gin.Context, key string, defaultVal int) int {
	str := c.Query(key)
	if str == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return val
}

// BindJSON 绑定JSON请求体，失败时自动返回错误响应
func BindJSON[T any](c *gin.Context) (*T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, 400, "请求参数错误: "+err.Error())
		return nil, false
	}
	return &req, true
}
