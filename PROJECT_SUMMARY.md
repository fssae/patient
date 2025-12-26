# 东软顾养中心系统 - 项目总结

## ✅ 已完成功能模块

### 1. 用户认证模块 ✅
- **功能**: 用户注册、登录
- **API端点**:
  - `POST /api/user/register` - 用户注册
  - `POST /api/user/login` - 用户登录
- **文件**: `internal/domain/User.go`, `internal/repository/dao/user.go`, `internal/repository/user.go`, `internal/service/user.go`, `internal/web/user.go`

### 2. 健康管家模块 ✅
- **功能**: 健康管家信息管理
- **API端点**:
  - `GET /api/health-managers` - 获取健康管家列表
  - `GET /api/health-managers/:id` - 获取健康管家详情
  - `POST /api/health-managers` - 创建健康管家
  - `PUT /api/health-managers/:id` - 更新健康管家
  - `DELETE /api/health-managers/:id` - 删除健康管家

### 3. 房间管理模块 ✅
- **功能**: 房间信息管理
- **API端点**:
  - `GET /api/rooms` - 获取房间列表
  - `GET /api/rooms/:id` - 获取房间详情
  - `POST /api/rooms` - 创建房间
  - `PUT /api/rooms/:id` - 更新房间
  - `DELETE /api/rooms/:id` - 删除房间

### 4. 床位管理模块 ✅
- **功能**: 床位信息管理、床位分配
- **API端点**:
  - `GET /api/beds` - 获取床位列表
  - `GET /api/beds/:id` - 获取床位详情
  - `GET /api/beds/room/:room_id` - 根据房间ID获取床位列表
  - `POST /api/beds` - 创建床位
  - `PUT /api/beds/:id/assign` - 分配床位给客户
  - `PUT /api/beds/:id/release` - 释放床位

### 5. 护理级别模块 ✅
- **功能**: 护理级别管理
- **API端点**:
  - `GET /api/care-levels` - 获取护理级别列表
  - `GET /api/care-levels/:id` - 获取护理级别详情
  - `POST /api/care-levels` - 创建护理级别
  - `PUT /api/care-levels/:id` - 更新护理级别
  - `DELETE /api/care-levels/:id` - 删除护理级别

### 6. 膳食管理模块 ✅
- **功能**: 膳食计划管理、每周菜单设置
- **API端点**:
  - `GET /api/diet-plans` - 获取膳食计划列表
  - `GET /api/diet-plans/:id` - 获取膳食计划详情
  - `POST /api/diet-plans` - 创建膳食计划
  - `PUT /api/diet-plans/:id` - 更新膳食计划
  - `DELETE /api/diet-plans/:id` - 删除膳食计划

### 7. 客户管理模块 ✅
- **功能**: 客户信息管理、客户关系设置
- **API端点**:
  - `GET /api/customers` - 获取客户列表
  - `GET /api/customers/:id` - 获取客户详情
  - `POST /api/customers` - 创建客户
  - `PUT /api/customers/:id` - 更新客户信息
  - `PUT /api/customers/:id/health-manager` - 设置健康管家
  - `PUT /api/customers/:id/bed` - 设置床位
  - `PUT /api/customers/:id/diet-plan` - 设置膳食计划
  - `PUT /api/customers/:id/care-level` - 设置护理级别
  - `DELETE /api/customers/:id` - 删除客户

### 8. 入住/退住/外出登记模块 ✅
- **功能**: 入住登记、退住登记、外出登记、外出返回
- **API端点**:
  - `POST /api/records/check-in` - 入住登记
  - `POST /api/records/check-out` - 退住登记
  - `POST /api/records/outgoing` - 外出登记
  - `POST /api/records/return` - 外出返回
  - `GET /api/records` - 获取登记记录列表
  - `GET /api/records/customer/:customer_id` - 获取客户的登记记录

### 9. 服务管理模块 ✅
- **功能**: 服务项目管理、客户购买服务、服务管理
- **API端点**:
  - `GET /api/services` - 获取服务项目列表
  - `GET /api/services/:id` - 获取服务项目详情
  - `POST /api/services` - 创建服务项目
  - `PUT /api/services/:id` - 更新服务项目
  - `DELETE /api/services/:id` - 删除服务项目
  - `POST /api/services/purchase` - 客户购买服务
  - `GET /api/services/customer/:customer_id` - 获取客户购买的服务列表
  - `PUT /api/services/customer-service/:id/end` - 结束客户服务

## 📁 项目结构

```
kongdong/
├── main.go                    # 应用入口
├── app.go                     # 应用启动配置
├── wire.go                    # Wire依赖注入配置
├── wire_gen.go               # Wire自动生成的代码
│
├── internal/
│   ├── domain/               # 领域模型层
│   │   ├── User.go           # 用户、客户、房间、床位等模型
│   │   ├── userClaims.go     # 用户JWT Claims
│   │   └── ...
│   │
│   ├── repository/           # 数据访问层
│   │   ├── dao/              # 数据访问对象（DAO）
│   │   │   ├── base.go       # 通用DAO基类
│   │   │   ├── user.go       # 用户DAO
│   │   │   ├── customer.go   # 客户DAO
│   │   │   ├── room.go       # 房间DAO
│   │   │   ├── bed.go        # 床位DAO
│   │   │   ├── care_level.go # 护理级别DAO
│   │   │   ├── diet_plan.go  # 膳食计划DAO
│   │   │   ├── record.go     # 登记记录DAO
│   │   │   ├── service.go    # 服务项目DAO
│   │   │   ├── customer_service.go # 客户服务DAO
│   │   │   └── health_manager.go  # 健康管家DAO
│   │   │
│   │   ├── user.go           # 用户Repository
│   │   ├── customer.go       # 客户Repository
│   │   ├── room.go           # 房间Repository
│   │   ├── bed.go            # 床位Repository
│   │   ├── care_level.go     # 护理级别Repository
│   │   ├── diet_plan.go      # 膳食计划Repository
│   │   ├── record.go         # 登记记录Repository
│   │   ├── service.go        # 服务项目Repository
│   │   ├── customer_service.go # 客户服务Repository
│   │   └── health_manager.go  # 健康管家Repository
│   │
│   ├── service/              # 业务逻辑层
│   │   ├── user.go           # 用户服务
│   │   ├── customer.go       # 客户服务
│   │   ├── room.go           # 房间服务
│   │   ├── bed.go            # 床位服务
│   │   ├── care_level.go     # 护理级别服务
│   │   ├── diet_plan.go      # 膳食计划服务
│   │   ├── record.go         # 登记记录服务
│   │   ├── service.go        # 服务项目服务
│   │   └── health_manager.go # 健康管家服务
│   │
│   ├── web/                  # Web层（HTTP处理）
│   │   ├── user.go           # 用户Handler
│   │   ├── customer.go       # 客户Handler
│   │   ├── room.go           # 房间Handler
│   │   ├── bed.go            # 床位Handler
│   │   ├── care_level.go     # 护理级别Handler
│   │   ├── diet_plan.go      # 膳食计划Handler
│   │   ├── record.go         # 登记记录Handler
│   │   ├── service.go        # 服务项目Handler
│   │   └── health_manager.go # 健康管家Handler
│   │
│   └── ioc/                  # 依赖注入配置
│       ├── minimal_sets.go   # Wire依赖集合
│       ├── gin.go            # Gin路由初始化
│       └── ...
│
└── config/                   # 配置文件
    └── conf.yaml
```

## 🚀 快速开始

### 1. 启动服务
```bash
go run main.go
```

服务将在 `http://localhost:8081` 启动

### 2. 测试API

#### 用户注册
```bash
POST http://localhost:8081/api/user/register
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "123456",
  "name": "张三",
  "age": 65,
  "gender": "男"
}
```

#### 用户登录
```bash
POST http://localhost:8081/api/user/login
Content-Type: application/json

{
  "phone": "13800138000",
  "password": "123456"
}
```

#### 创建房间
```bash
POST http://localhost:8081/api/rooms
Content-Type: application/json

{
  "number": "A101",
  "floor": 1,
  "type": "单人间",
  "capacity": 1,
  "description": "朝南，采光好"
}
```

#### 创建床位
```bash
POST http://localhost:8081/api/beds
Content-Type: application/json

{
  "room_id": "房间ID",
  "number": "A101-1"
}
```

#### 创建客户
```bash
POST http://localhost:8081/api/customers
Content-Type: application/json

{
  "user_id": "用户ID",
  "name": "李四",
  "age": 70,
  "gender": "女",
  "phone": "13900139000",
  "id_card": "110101195001011234"
}
```

#### 入住登记
```bash
POST http://localhost:8081/api/records/check-in
Content-Type: application/json

{
  "customer_id": "客户ID",
  "bed_id": "床位ID",
  "note": "客户入住",
  "created_by": "管理员"
}
```

## 📊 数据库集合

系统使用MongoDB，包含以下集合：
- `users` - 用户表
- `customers` - 客户表
- `rooms` - 房间表
- `beds` - 床位表
- `diet_plans` - 膳食计划表
- `care_levels` - 护理级别表
- `records` - 登记记录表
- `services` - 服务项目表
- `customer_services` - 客户服务表
- `health_managers` - 健康管家表

## 🔧 技术栈

- **语言**: Go 1.25.1
- **Web框架**: Gin
- **数据库**: MongoDB
- **缓存**: Redis
- **对象存储**: MinIO
- **消息队列**: Kafka
- **依赖注入**: Google Wire
- **JWT认证**: golang-jwt/jwt/v5
- **密码加密**: bcrypt

## 📝 注意事项

1. **JWT配置**: 需要在配置文件中设置 `general.jwt` 作为JWT密钥
2. **数据库连接**: 确保MongoDB、Redis等服务已启动并配置正确
3. **密码加密**: 所有密码使用bcrypt加密存储
4. **错误处理**: 统一使用JSON格式返回错误信息
5. **数据验证**: 使用Gin的binding标签进行参数验证

## 🎯 下一步优化建议

1. **权限控制**: 添加基于角色的访问控制（RBAC）
2. **数据验证**: 增强输入数据验证和校验
3. **日志记录**: 完善操作日志记录
4. **单元测试**: 为各模块编写单元测试
5. **API文档**: 使用Swagger生成API文档
6. **性能优化**: 添加缓存策略，优化查询性能
7. **监控告警**: 集成Prometheus监控和告警

## ✨ 项目特点

- ✅ 完整的分层架构设计
- ✅ 使用Wire进行依赖注入
- ✅ 统一的错误处理和响应格式
- ✅ 完整的CRUD操作
- ✅ 业务逻辑封装在Service层
- ✅ 数据访问抽象在Repository层
- ✅ 支持分页查询
- ✅ 支持状态管理和业务流转

项目已完整实现所有需求功能，可以直接使用！

