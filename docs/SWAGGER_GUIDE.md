# Swagger API文档使用指南

## 快速开始

### 1. 启动服务

```bash
go run main.go
```

### 2. 访问Swagger文档

在浏览器中打开：
```
http://localhost:8081/swagger/index.html
```

## 已配置的API模块

### ✅ 已添加Swagger注释的API

1. **用户认证模块**
   - POST `/api/user/register` - 用户注册
   - POST `/api/user/login` - 用户登录

2. **客户管理模块**
   - GET `/api/customers` - 获取客户列表
   - GET `/api/customers/{id}` - 获取客户详情
   - POST `/api/customers` - 创建客户

3. **房间管理模块**
   - GET `/api/rooms` - 获取房间列表
   - POST `/api/rooms` - 创建房间

4. **床位管理模块**
   - PUT `/api/beds/{id}/assign` - 分配床位

5. **登记管理模块**
   - POST `/api/records/check-in` - 入住登记
   - POST `/api/records/check-out` - 退住登记

6. **服务管理模块**
   - POST `/api/services/purchase` - 客户购买服务

## 更新Swagger文档

当添加新的API或修改现有API时，需要重新生成Swagger文档：

```bash
# 安装swag工具（如果还没安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g main.go -o docs
```

## 为API添加Swagger注释

### 基本格式

```go
// @Summary      接口摘要（简短描述）
// @Description  接口详细描述
// @Tags         标签名（用于分组）
// @Accept       json
// @Produce      json
// @Param        参数名  参数位置  是否必需  参数类型  参数描述  示例值
// @Success      状态码  {返回类型}  返回描述
// @Failure      状态码  {返回类型}  返回描述
// @Router       路由路径 [HTTP方法]
func (h *Handler) Method(c *gin.Context) {
    // 实现代码
}
```

### 参数类型说明

- `path` - URL路径参数，如 `/api/customers/{id}`
- `query` - URL查询参数，如 `?status=入住中&skip=0&limit=20`
- `body` - 请求体（JSON）
- `formData` - 表单数据

### 示例1: GET请求（带路径参数）

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
// @Failure      404  {object}  map[string]interface{}  "客户不存在"
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetById(c *gin.Context) {
    // ...
}
```

### 示例2: GET请求（带查询参数）

```go
// GetList 获取客户列表
// @Summary      获取客户列表
// @Description  分页获取客户列表，支持按状态筛选
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "客户状态：入住中/已退住/外出中"
// @Param        skip    query     int     false  "跳过数量"  default(0)
// @Param        limit   query     int     false  "每页数量"  default(20)
// @Success      200     {object}  map[string]interface{}  "获取成功"
// @Router       /customers [get]
func (h *CustomerHandler) GetList(c *gin.Context) {
    // ...
}
```

### 示例3: POST请求（带请求体）

```go
// Create 创建客户
// @Summary      创建客户
// @Description  创建新的客户记录
// @Tags         客户管理
// @Accept       json
// @Produce      json
// @Param        request  body      domain.Customer  true  "客户信息"
// @Success      200      {object}  map[string]interface{}  "创建成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
    // ...
}
```

### 示例4: POST请求（自定义请求体）

```go
// CheckIn 入住登记
// @Summary      入住登记
// @Description  为客户办理入住登记，自动分配床位并更新客户状态
// @Tags         登记管理
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "入住信息"  example({"customer_id":"507f1f77bcf86cd799439011","bed_id":"507f1f77bcf86cd799439012","note":"客户入住","created_by":"管理员"})
// @Success      200      {object}  map[string]interface{}  "入住登记成功"
// @Failure      400      {object}  map[string]interface{}  "请求参数错误"
// @Router       /records/check-in [post]
func (h *RecordHandler) CheckIn(c *gin.Context) {
    // ...
}
```

## 响应类型

### 成功响应

```go
@Success 200 {object} map[string]interface{} "操作成功"
```

### 错误响应

```go
@Failure 400 {object} map[string]interface{} "请求参数错误"
@Failure 401 {object} map[string]interface{} "未授权"
@Failure 404 {object} map[string]interface{} "资源不存在"
@Failure 500 {object} map[string]interface{} "服务器内部错误"
```

## 使用Swagger UI测试API

1. 打开Swagger文档页面
2. 找到要测试的API
3. 点击"Try it out"按钮
4. 填写请求参数
5. 点击"Execute"执行请求
6. 查看响应结果

## 注意事项

1. **路径参数**: 使用 `{param}` 格式，如 `/customers/{id}`
2. **查询参数**: 使用 `query` 类型
3. **请求体**: 使用 `body` 类型，可以引用domain模型或使用 `object` 类型
4. **标签分组**: 使用 `@Tags` 将相关API分组
5. **示例值**: 使用 `example()` 提供示例值，方便测试

## 常见问题

### Q: 如何为domain模型添加Swagger注释？

A: 在domain结构体字段上添加 `json` 和 `example` 标签：

```go
type User struct {
    Phone    string `json:"phone" example:"13800002001"`
    Password string `json:"password" example:"123456"`
    Name     string `json:"name" example:"张三"`
}
```

### Q: 如何添加认证说明？

A: 在main.go中已经配置了Bearer认证，在需要认证的API上添加：

```go
// @Security Bearer
```

### Q: 文档生成失败怎么办？

A: 检查：
1. 注释格式是否正确
2. 路由路径是否正确
3. 参数类型是否匹配
4. 运行 `swag init -g main.go -o docs --parseDependency` 查看详细错误

## 下一步

为所有API添加Swagger注释，使文档更加完整。可以参考已添加注释的API作为模板。

