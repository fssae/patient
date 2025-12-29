# Apifox 导入 Swagger 文档指南

## Swagger文档位置

Swagger文档已生成在以下位置：
- **JSON格式**: `docs/swagger.json`
- **YAML格式**: `docs/swagger.yaml`

## 导入到Apifox

### 方法1: 通过文件导入（推荐）

1. 打开Apifox应用
2. 选择项目 → **导入** → **OpenAPI/Swagger**
3. 选择 **从文件导入**
4. 选择 `docs/swagger.json` 或 `docs/swagger.yaml` 文件
5. 点击 **导入**

### 方法2: 通过URL导入

如果服务正在运行，可以通过URL导入：

1. 打开Apifox应用
2. 选择项目 → **导入** → **OpenAPI/Swagger**
3. 选择 **从URL导入**
4. 输入URL: `http://localhost:8081/swagger/doc.json`
5. 点击 **导入**

### 方法3: 复制粘贴JSON内容

1. 打开 `docs/swagger.json` 文件
2. 复制全部内容
3. 打开Apifox应用
4. 选择项目 → **导入** → **OpenAPI/Swagger**
5. 选择 **从文本导入**
6. 粘贴JSON内容
7. 点击 **导入**

## 导入后配置

### 1. 设置环境变量

导入后，建议在Apifox中设置环境变量：

- **base_url**: `http://localhost:8081`
- **api_base**: `/api`

### 2. 配置认证

如果需要测试需要JWT认证的接口：

1. 在Apifox中设置 **环境变量**
2. 添加变量 `token`，值从登录接口获取
3. 在需要认证的接口中，添加Header：
   - Key: `Authorization`
   - Value: `Bearer {{token}}`

### 3. 测试登录接口

1. 先调用 `/api/user/login` 接口
2. 获取返回的 `token`
3. 将token保存到环境变量中
4. 后续接口会自动使用该token

## 当前包含的API

文档包含以下API模块：

### ✅ 已添加Swagger注释的API

1. **用户认证**
   - POST `/api/user/register` - 用户注册
   - POST `/api/user/login` - 用户登录

2. **客户管理**
   - GET `/api/customers` - 获取客户列表
   - GET `/api/customers/{id}` - 获取客户详情
   - POST `/api/customers` - 创建客户

3. **房间管理**
   - GET `/api/rooms` - 获取房间列表
   - POST `/api/rooms` - 创建房间

4. **床位管理**
   - PUT `/api/beds/{id}/assign` - 分配床位

5. **登记管理**
   - POST `/api/records/check-in` - 入住登记
   - POST `/api/records/check-out` - 退住登记

6. **服务管理**
   - POST `/api/services/purchase` - 客户购买服务

## 更新文档

当API有更新时，重新生成Swagger文档：

```bash
swag init -g main.go -o docs
```

然后重新导入到Apifox即可。

## 注意事项

1. **文件路径**: 确保使用正确的文件路径 `docs/swagger.json`
2. **服务运行**: 如果使用URL导入，确保服务正在运行
3. **格式选择**: Apifox支持JSON和YAML格式，推荐使用JSON
4. **环境配置**: 导入后记得配置环境变量和认证信息

## 快速测试流程

1. 导入Swagger文档到Apifox
2. 配置环境变量（base_url等）
3. 调用 `/api/user/login` 获取token
4. 设置token到环境变量
5. 测试其他需要认证的接口

## 问题排查

### Q: 导入失败怎么办？

A: 检查：
- JSON文件格式是否正确
- 文件路径是否正确
- Apifox版本是否支持OpenAPI 2.0

### Q: 接口显示不完整？

A: 可能原因：
- 某些API还没有添加Swagger注释
- 需要重新生成文档并导入

### Q: 如何添加更多API到文档？

A: 
1. 为Handler方法添加Swagger注释
2. 运行 `swag init -g main.go -o docs`
3. 重新导入到Apifox

