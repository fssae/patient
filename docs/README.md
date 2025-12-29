# Swagger API 文档

## 访问方式

启动服务后，访问以下地址查看Swagger文档：

```
http://localhost:8081/swagger/index.html
```

## 生成文档

如果需要更新Swagger文档，运行以下命令：

```bash
swag init -g main.go -o docs
```

## 文档说明

Swagger文档包含了以下API模块：

1. **用户认证** - 用户注册、登录
2. **客户管理** - 客户CRUD操作、关系设置
3. **房间管理** - 房间CRUD操作
4. **床位管理** - 床位CRUD、分配、释放
5. **护理级别** - 护理级别管理
6. **膳食管理** - 膳食计划管理
7. **健康管家** - 健康管家管理
8. **登记管理** - 入住、退住、外出登记
9. **服务管理** - 服务项目、客户服务管理

## 添加API注释

为API添加Swagger注释的格式：

```go
// @Summary      接口摘要
// @Description  接口详细描述
// @Tags         标签名
// @Accept       json
// @Produce      json
// @Param        param_name  param_type  required  description  example(value)
// @Success      200  {object}  map[string]interface{}  "成功响应"
// @Failure      400  {object}  map[string]interface{}  "错误响应"
// @Router       /path [method]
func (h *Handler) Method(c *gin.Context) {
    // ...
}
```

## 参数说明

- `@Summary`: API简要说明
- `@Description`: API详细描述
- `@Tags`: API分组标签
- `@Accept`: 接受的Content-Type
- `@Produce`: 返回的Content-Type
- `@Param`: 参数定义
  - `param_name`: 参数名
  - `param_type`: 参数类型（path/query/body/formData）
  - `required`: 是否必需（true/false）
  - `description`: 参数描述
  - `example`: 示例值
- `@Success`: 成功响应
- `@Failure`: 失败响应
- `@Router`: 路由路径和方法

## 示例

```go
// GetById 获取客户详情
// @Summary      获取客户详情
// @Description  根据ID获取客户详细信息
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "客户ID"
// @Success      200  {object}  map[string]interface{}  "获取成功"
// @Failure      400  {object}  map[string]interface{}  "无效的ID"
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetById(c *gin.Context) {
    // ...
}
```

