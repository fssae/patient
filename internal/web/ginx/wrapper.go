package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Result 统一响应结构
type Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
}

// Page 分页参数
type Page struct {
	Skip  int64
	Limit int64
}

// GetPage 从gin.Context解析分页参数
func GetPage(c *gin.Context) Page {
	var page Page
	page.Skip = 0
	page.Limit = 20
	if skipStr := c.Query("skip"); skipStr != "" {
		for _, ch := range skipStr {
			if ch >= '0' && ch <= '9' {
				page.Skip = page.Skip*10 + int64(ch-'0')
			}
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		for _, ch := range limitStr {
			if ch >= '0' && ch <= '9' {
				page.Limit = page.Limit*10 + int64(ch-'0')
			}
		}
	}
	return page
}

// GetId 从路径参数获取ObjectID
func GetId(c *gin.Context) (primitive.ObjectID, bool) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 400, Msg: "无效的ID"})
		return primitive.NilObjectID, false
	}
	return id, true
}

// WrapBody 包装需要请求体的处理函数
func WrapBody[Req any](fn func(c *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, Result{
				Code: 400,
				Msg:  "请求参数错误: " + err.Error(),
			})
			return
		}
		result, err := fn(c, req)
		if err != nil {
			c.JSON(http.StatusOK, Result{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// Wrap 包装不需要请求体的处理函数
func Wrap(fn func(c *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := fn(c)
		if err != nil {
			c.JSON(http.StatusOK, Result{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// Ok 成功响应
func Ok(msg string, data interface{}) Result {
	return Result{Code: 200, Msg: msg, Success: true, Data: data}
}

// OkMsg 仅消息的成功响应
func OkMsg(msg string) Result {
	return Result{Code: 200, Msg: msg, Success: true}
}

// OkList 列表成功响应
func OkList(msg string, data interface{}, total int64) Result {
	return Result{Code: 200, Msg: msg, Success: true, Data: data, Total: total}
}

// Fail 失败响应
func Fail(code int, msg string) Result {
	return Result{Code: code, Msg: msg}
}
