# 东软顾养中心系统开发指南

## 系统架构概述

本系统采用**分层架构**设计，遵循**领域驱动设计（DDD）**思想，使用**依赖注入（Wire）**管理组件依赖。

### 架构层次

```
┌─────────────────────────────────────────┐
│  Web层 (Handler)                       │  ← HTTP请求处理
│  - 参数验证                            │
│  - 响应格式化                          │
│  - 路由注册                            │
├─────────────────────────────────────────┤
│  Service层 (业务逻辑)                   │  ← 核心业务逻辑
│  - 业务规则验证                        │
│  - 事务管理                            │
│  - 调用Repository                      │
├─────────────────────────────────────────┤
│  Repository层 (数据访问抽象)            │  ← 数据访问接口
│  - 定义数据访问接口                    │
│  - 实现数据访问逻辑                    │
├─────────────────────────────────────────┤
│  DAO层 (数据操作)                       │  ← 数据库操作
│  - MongoDB操作                         │
│  - 基础CRUD                            │
├─────────────────────────────────────────┤
│  Domain层 (领域模型)                     │  ← 业务实体
│  - 实体定义                            │
│  - 值对象                              │
│  - 领域事件                            │
└─────────────────────────────────────────┘
```

## 已完成的模块

### ✅ 1. 用户认证模块
- **功能**: 用户注册、登录
- **API端点**:
  - `POST /api/user/register` - 用户注册
  - `POST /api/user/login` - 用户登录
- **文件结构**:
  ```
  internal/
  ├── domain/
  │   ├── User.go          # 用户模型
  │   └── userClaims.go     # JWT Claims
  ├── repository/
  │   ├── dao/user.go      # 用户DAO
  │   └── user.go          # 用户Repository
  ├── service/
  │   └── user.go          # 用户Service
  └── web/
      └── user.go          # 用户Handler
  ```

## 待开发的模块

### 📋 2. 客户管理模块 (Customer)
**功能需求**:
- 客户信息管理（增删改查）
- 设置客户与健康管家关系
- 客户状态管理（入住中/已退住/外出中）

**开发步骤**:
1. 创建 `internal/repository/dao/customer.go` - Customer DAO
2. 创建 `internal/repository/customer.go` - Customer Repository
3. 创建 `internal/service/customer.go` - Customer Service
4. 创建 `internal/web/customer.go` - Customer Handler
5. 在 `internal/ioc/minimal_sets.go` 中添加依赖
6. 在 `internal/ioc/gin.go` 中注册路由

**API设计示例**:
```go
POST   /api/customers              # 创建客户
GET    /api/customers              # 获取客户列表
GET    /api/customers/:id          # 获取客户详情
PUT    /api/customers/:id          # 更新客户信息
DELETE /api/customers/:id          # 删除客户
PUT    /api/customers/:id/health-manager  # 设置健康管家
```

### 📋 3. 床位管理模块 (Bed & Room)
**功能需求**:
- 房间管理（创建、查询、更新）
- 床位管理（创建、查询、分配）
- 床位状态管理（空闲/占用/维护中）

**开发步骤**:
1. 创建 Room DAO、Repository、Service、Handler
2. 创建 Bed DAO、Repository、Service、Handler
3. 实现床位分配逻辑（关联客户）

**API设计示例**:
```go
# 房间管理
POST   /api/rooms                 # 创建房间
GET    /api/rooms                 # 获取房间列表
GET    /api/rooms/:id             # 获取房间详情
PUT    /api/rooms/:id             # 更新房间信息

# 床位管理
POST   /api/beds                  # 创建床位
GET    /api/beds                  # 获取床位列表
PUT    /api/beds/:id/assign       # 分配床位给客户
PUT    /api/beds/:id/release      # 释放床位
```

### 📋 4. 膳食管理模块 (DietPlan)
**功能需求**:
- 膳食计划管理
- 每周伙食菜单设置
- 为客户定制膳食计划

**开发步骤**:
1. 创建 DietPlan DAO、Repository、Service、Handler
2. 实现每周菜单管理（WeekDayMenu）
3. 实现客户膳食计划分配

**API设计示例**:
```go
POST   /api/diet-plans            # 创建膳食计划
GET    /api/diet-plans            # 获取膳食计划列表
GET    /api/diet-plans/:id        # 获取膳食计划详情
PUT    /api/diet-plans/:id        # 更新膳食计划
PUT    /api/customers/:id/diet-plan  # 为客户设置膳食计划
```

### 📋 5. 入住/退住/外出登记模块 (Record)
**功能需求**:
- 入住登记
- 退住登记
- 外出登记
- 登记记录查询

**开发步骤**:
1. 创建 Record DAO、Repository、Service、Handler
2. 实现入住登记逻辑（更新客户状态、床位状态）
3. 实现退住登记逻辑（释放床位、更新客户状态）
4. 实现外出登记逻辑（更新客户状态）

**API设计示例**:
```go
POST   /api/records/check-in       # 入住登记
POST   /api/records/check-out      # 退住登记
POST   /api/records/outgoing       # 外出登记
POST   /api/records/return         # 外出返回
GET    /api/records                # 查询登记记录
GET    /api/customers/:id/records  # 查询客户登记记录
```

### 📋 6. 护理级别模块 (CareLevel)
**功能需求**:
- 护理级别管理
- 为客户设置护理级别

**开发步骤**:
1. 创建 CareLevel DAO、Repository、Service、Handler
2. 实现护理级别CRUD
3. 实现客户护理级别分配

**API设计示例**:
```go
POST   /api/care-levels            # 创建护理级别
GET    /api/care-levels            # 获取护理级别列表
PUT    /api/care-levels/:id        # 更新护理级别
PUT    /api/customers/:id/care-level  # 为客户设置护理级别
```

### 📋 7. 服务管理模块 (Service)
**功能需求**:
- 服务项目管理
- 客户购买的服务信息管理
- 服务关注设置（服务对象）

**开发步骤**:
1. 创建 Service DAO、Repository、Service、Handler
2. 创建 CustomerService DAO、Repository、Service、Handler
3. 创建 ServiceAttention DAO、Repository、Service、Handler
4. 实现服务购买逻辑
5. 实现服务关注设置逻辑

**API设计示例**:
```go
# 服务项目管理
POST   /api/services               # 创建服务项目
GET    /api/services               # 获取服务项目列表

# 客户服务
POST   /api/customer-services      # 客户购买服务
GET    /api/customers/:id/services # 查询客户服务

# 服务关注
POST   /api/service-attentions     # 设置服务关注
GET    /api/customers/:id/attentions  # 查询客户服务关注
```

### 📋 8. 健康管家模块 (HealthManager)
**功能需求**:
- 健康管家信息管理
- 健康管家与客户关系管理

**开发步骤**:
1. 创建 HealthManager DAO、Repository、Service、Handler
2. 实现健康管家CRUD
3. 在Customer模块中实现关系管理

## 开发模板

### DAO层模板
```go
package dao

import (
    "classroom-analysis/internal/domain"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type XxxDAO struct {
    collection *mongo.Collection
}

func NewXxxDAO(db *mongo.Database) *XxxDAO {
    return &XxxDAO{
        collection: db.Collection("xxx_collection"),
    }
}

func (dao *XxxDAO) Create(ctx context.Context, entity *domain.Xxx) error {
    entity.CreatedAt = time.Now()
    entity.UpdatedAt = time.Now()
    result, err := dao.collection.InsertOne(ctx, entity)
    if err != nil {
        return err
    }
    entity.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

// 其他方法...
```

### Repository层模板
```go
package repository

import (
    "classroom-analysis/internal/domain"
    "classroom-analysis/internal/repository/dao"
    "context"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type XxxRepository interface {
    Create(ctx context.Context, entity *domain.Xxx) error
    FindById(ctx context.Context, id primitive.ObjectID) (*domain.Xxx, error)
    // 其他方法...
}

type xxxRepository struct {
    dao *dao.XxxDAO
}

func NewXxxRepository(xxxDAO *dao.XxxDAO) XxxRepository {
    return &xxxRepository{
        dao: xxxDAO,
    }
}

func (r *xxxRepository) Create(ctx context.Context, entity *domain.Xxx) error {
    return r.dao.Create(ctx, entity)
}

// 其他方法实现...
```

### Service层模板
```go
package service

import (
    "classroom-analysis/internal/domain"
    "classroom-analysis/internal/repository"
    "context"
    "errors"
)

type XxxService struct {
    xxxRepo repository.XxxRepository
}

func NewXxxService(xxxRepo repository.XxxRepository) *XxxService {
    return &XxxService{
        xxxRepo: xxxRepo,
    }
}

func (s *XxxService) Create(ctx context.Context, req *domain.XxxRequest) error {
    // 业务逻辑验证
    // 调用Repository
    return s.xxxRepo.Create(ctx, entity)
}

// 其他方法...
```

### Handler层模板
```go
package web

import (
    "classroom-analysis/internal/domain"
    "classroom-analysis/internal/service"
    "net/http"
    "github.com/gin-gonic/gin"
)

type XxxHandler struct {
    svc *service.XxxService
}

func NewXxxHandler(svc *service.XxxService) *XxxHandler {
    return &XxxHandler{
        svc: svc,
    }
}

func (h *XxxHandler) Create(c *gin.Context) {
    var req domain.XxxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 400,
            "msg":  "请求参数错误: " + err.Error(),
        })
        return
    }

    err := h.svc.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": 400,
            "msg":  err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "msg":     "创建成功",
        "success": true,
    })
}

func (h *XxxHandler) RegisterRoutes(server *gin.Engine) {
    group := server.Group("/api/xxx")
    // 需要JWT验证的路由
    group.Use(middleware.NewLoginJWTMiddlewareBuilder().Build())
    group.POST("", h.Create)
    // 其他路由...
}
```

## 依赖注入配置

在 `internal/ioc/minimal_sets.go` 中添加新模块：

```go
var MinimalSet = wire.NewSet(
    // ... 现有配置
    
    // 新模块DAO
    dao.NewXxxDAO,
    
    // 新模块Repository
    repository.NewXxxRepository,
    
    // 新模块Service
    service.NewXxxService,
    
    // 新模块Handler
    web.NewXxxHandler,
    
    // 更新Gin初始化
    InitGin,
)
```

在 `internal/ioc/gin.go` 中注册路由：

```go
func InitGin(..., xxxHandler *web.XxxHandler) *gin.Engine {
    // ...
    xxxHandler.RegisterRoutes(engine)
    return engine
}
```

然后运行 `go generate ./...` 重新生成Wire代码。

## 测试建议

1. **单元测试**: 为Service层编写单元测试
2. **集成测试**: 测试API端点
3. **数据验证**: 确保数据校验正确
4. **错误处理**: 测试各种错误场景

## 注意事项

1. **密码加密**: 使用 `bcrypt` 加密存储密码
2. **JWT Token**: 用户登录后生成JWT token，后续请求需要携带
3. **数据验证**: 使用 `binding` 标签进行参数验证
4. **错误处理**: 统一错误响应格式
5. **事务管理**: 复杂业务操作需要考虑事务（如入住登记需要更新多个表）

## 下一步

1. 按照上述模板创建各个模块
2. 实现业务逻辑
3. 编写测试用例
4. 完善API文档
5. 部署和上线

