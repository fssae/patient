# 东软顾养中心系统架构说明

## 系统分层架构

本系统采用经典的分层架构模式，分为以下层次：

```
┌─────────────────────────────────────┐
│         Web层 (Handler)             │  处理HTTP请求/响应
├─────────────────────────────────────┤
│         Service层 (业务逻辑)         │  业务逻辑处理
├─────────────────────────────────────┤
│      Repository层 (数据访问)        │  数据访问抽象
├─────────────────────────────────────┤
│         DAO层 (数据操作)            │  数据库操作
├─────────────────────────────────────┤
│         Domain层 (领域模型)          │  业务实体定义
└─────────────────────────────────────┘
```

## 目录结构说明

```
kongdong/
├── main.go                    # 应用入口
├── app.go                     # 应用启动配置
├── wire.go                    # Wire依赖注入配置
├── wire_gen.go               # Wire自动生成的代码
│
├── internal/
│   ├── domain/               # 领域模型层
│   │   ├── User.go           # 用户模型（注册用户）
│   │   ├── userClaims.go     # 用户JWT Claims
│   │   └── ...               # 其他领域模型
│   │
│   ├── repository/           # 数据访问层
│   │   ├── dao/              # 数据访问对象（DAO）
│   │   │   ├── base.go       # 通用DAO基类
│   │   │   ├── user.go       # 用户DAO
│   │   │   └── ...           # 其他DAO
│   │   │
│   │   ├── user.go           # 用户Repository接口和实现
│   │   └── ...               # 其他Repository
│   │
│   ├── service/              # 业务逻辑层
│   │   ├── user.go           # 用户服务
│   │   └── ...               # 其他服务
│   │
│   ├── web/                  # Web层（HTTP处理）
│   │   ├── user.go           # 用户Handler
│   │   ├── handler.go        # Handler基类
│   │   └── middleware/       # 中间件
│   │
│   └── ioc/                  # 依赖注入配置
│       ├── minimal_sets.go   # Wire依赖集合
│       ├── gin.go            # Gin路由初始化
│       └── ...               # 其他IOC配置
│
└── config/                   # 配置文件
    └── conf.yaml
```

## 模块划分

### 1. 用户认证模块 (User)
- **功能**: 用户注册、登录
- **文件**: 
  - `domain/User.go` - 用户模型
  - `repository/dao/user.go` - 用户DAO
  - `repository/user.go` - 用户Repository
  - `service/user.go` - 用户Service
  - `web/user.go` - 用户Handler

### 2. 客户管理模块 (Customer)
- **功能**: 客户信息管理、客户与健康管家关系设置
- **领域模型**: `domain/User.go` 中的 `Customer` 结构

### 3. 床位管理模块 (Bed)
- **功能**: 房间和床位设置、床位分配
- **领域模型**: `domain/User.go` 中的 `Bed`、`Room` 结构

### 4. 膳食管理模块 (DietPlan)
- **功能**: 膳食计划设置、每周伙食菜单定制
- **领域模型**: `domain/User.go` 中的 `DietPlan`、`WeekDayMenu` 结构

### 5. 入住/退住/外出登记模块 (Record)
- **功能**: 入住登记、退住登记、外出登记
- **领域模型**: `domain/User.go` 中的 `Record` 结构

### 6. 护理级别模块 (CareLevel)
- **功能**: 护理级别设置、客户护理级别分配
- **领域模型**: `domain/User.go` 中的 `CareLevel` 结构

### 7. 服务管理模块 (Service)
- **功能**: 服务项目设置、客户购买的服务信息、服务关注设置
- **领域模型**: `domain/User.go` 中的 `Service`、`CustomerService`、`ServiceAttention` 结构

### 8. 健康管家模块 (HealthManager)
- **功能**: 健康管家信息管理
- **领域模型**: `domain/User.go` 中的 `HealthManager` 结构

## 开发规范

### 1. 命名规范
- **Domain**: 使用大驼峰，如 `User`、`Customer`
- **DAO**: 以 `DAO` 结尾，如 `UserDAO`
- **Repository**: 以 `Repository` 结尾，如 `UserRepository`
- **Service**: 以 `Service` 结尾，如 `UserService`
- **Handler**: 以 `Handler` 结尾，如 `UserHandler`

### 2. 接口定义
- Repository层定义接口，实现结构体使用小写开头
- Service层直接使用结构体（可后续扩展为接口）

### 3. 依赖注入
- 使用 Google Wire 进行依赖注入
- 在 `internal/ioc/minimal_sets.go` 中配置依赖关系
- 运行 `go generate ./...` 生成 Wire 代码

### 4. 路由注册
- 每个Handler实现 `RegisterRoutes` 方法
- 在 `internal/ioc/gin.go` 中统一注册路由

### 5. 错误处理
- Service层返回业务错误
- Handler层处理HTTP状态码和响应格式

## API设计规范

### 请求格式
```json
{
  "phone": "13800138000",
  "password": "password123"
}
```

### 响应格式
```json
{
  "code": 200,
  "msg": "操作成功",
  "success": true,
  "data": {...},
  "token": "jwt_token_string"
}
```

### 状态码
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权（登录失败）
- `500`: 服务器内部错误

## 数据库设计

使用 MongoDB，集合命名规范：
- `users` - 用户表
- `customers` - 客户表
- `rooms` - 房间表
- `beds` - 床位表
- `diet_plans` - 膳食计划表
- `care_levels` - 护理级别表
- `records` - 登记记录表
- `services` - 服务项目表
- `customer_services` - 客户服务表
- `service_attentions` - 服务关注表
- `health_managers` - 健康管家表

## 下一步开发建议

1. **完成核心模块开发**:
   - [x] 用户认证模块（注册、登录）
   - [ ] 客户管理模块
   - [ ] 床位管理模块
   - [ ] 膳食管理模块
   - [ ] 入住/退住/外出登记模块
   - [ ] 护理级别模块
   - [ ] 服务管理模块

2. **完善功能**:
   - JWT中间件支持用户认证
   - 权限控制（管理员、普通用户等）
   - 数据验证和校验
   - 日志记录
   - 单元测试

3. **优化**:
   - 缓存策略（Redis）
   - 文件上传（MinIO）
   - 消息队列（Kafka）
   - 监控和告警（Prometheus）

