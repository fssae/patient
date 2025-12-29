# Apifox 导入说明

## 📁 文件说明

- `swagger.json` - Swagger JSON格式文档（推荐使用）
- `swagger.yaml` - Swagger YAML格式文档

## 🚀 快速导入

### 步骤1: 打开Apifox
打开Apifox应用，选择要导入的项目

### 步骤2: 导入文档
1. 点击 **导入** → **OpenAPI/Swagger**
2. 选择 **从文件导入**
3. 选择 `swagger.json` 文件
4. 点击 **导入**

### 步骤3: 配置环境
导入后，建议配置以下环境变量：

```
base_url = http://localhost:8081
api_base = /api
```

## 📋 包含的API

当前文档包含以下API：

- ✅ 用户注册/登录
- ✅ 客户管理（列表、详情、创建）
- ✅ 房间管理（列表、创建）
- ✅ 床位管理（分配）
- ✅ 登记管理（入住、退住）
- ✅ 服务管理（购买服务）

## 🔄 更新文档

当API有更新时：

```bash
swag init -g main.go -o docs
```

然后重新导入到Apifox。

## 💡 提示

- 推荐使用 `swagger.json` 格式导入
- 导入后记得配置环境变量
- 先测试登录接口获取token，再测试其他接口

