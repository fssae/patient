package util

import (
	"classroom-analysis/internal/domain"
	"errors"
	"fmt"
	"net/http"
	"reflect"

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
func Validate(model interface{}, c *gin.Context) (map[string]interface{}, error) {
	fieldMap := GenerateFieldMap(model)
	fields, err := BindAndValidateFields(fieldMap)
	return fields, err
}
func GenerateFieldMap(model interface{}) map[string]bool {
	result := make(map[string]bool)

	// 获取结构体类型和值
	t := reflect.TypeOf(model)

	// 解包指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 获取 bson 标签作为字段名
		bsonTag := field.Tag.Get("bson")
		if bsonTag != "" && bsonTag != "-" {
			// 提取标签中的字段名（去掉可能的选项）
			if idx := index(bsonTag, ","); idx != -1 {
				bsonTag = bsonTag[:idx]
			}
			if bsonTag != "" {
				result[bsonTag] = true
			}
		} else {
			// 如果没有 bson 标签，使用 json 标签
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				if idx := index(jsonTag, ","); idx != -1 {
					jsonTag = jsonTag[:idx]
				}
				if jsonTag != "" {
					result[jsonTag] = true
				}
			} else {
				// 如果都没有，使用字段名的小写形式
				result[toLowerFirst(field.Name)] = true
			}
		}
	}

	return result
}

// index 是一个简单的字符串查找函数
func index(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// toLowerFirst 将字符串的首字母转为小写
func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]|' ') + s[1:]
}

// BindAndValidateFields 绑定JSON并验证字段是否在允许的字段列表中
func BindAndValidateFields(c *gin.Context, allowedFields map[string]bool) (map[string]interface{}, error) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, fmt.Errorf("请求参数错误: %w", err)
	}

	// 验证请求字段是否在允许的字段列表中
	for field := range req {
		if !allowedFields[field] {
			return nil, fmt.Errorf("无效的业务字段: %s", field)
		}
	}

	return req, nil
}
