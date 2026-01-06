# 后端项目优化建议

## 📋 目录
- [关键安全问题](#关键安全问题)
- [性能优化](#性能优化)
- [架构改进](#架构改进)
- [代码质量](#代码质量)
- [数据库优化](#数据库优化)
- [可观测性](#可观测性)
- [API设计](#api设计)
- [配置管理](#配置管理)

---

## 🔐 关键安全问题（高优先级）

### 1. **敏感信息泄露**
**问题：** `config/conf.yaml` 包含明文密码和密钥
```yaml
mongodb:
  password: zjh770910  # ❌ 明文密码
jwt: "zjh770910"      # ❌ JWT密钥
```

**建议：**
- ✅ 使用环境变量管理敏感信息
- ✅ 使用 `.env` 文件（不提交到Git）
- ✅ 生产环境使用密钥管理服务（如AWS Secrets Manager, HashiCorp Vault）
- ✅ 更新 `.gitignore` 确保 `config/conf.yaml` 不被提交

**实施方案：**
```go
// 使用环境变量
mongoPassword := os.Getenv("MONGO_PASSWORD")
if mongoPassword == "" {
    mongoPassword = viper.GetString("mongodb.password") // fallback
}
```

### 2. **JWT密钥安全性**
**问题：** JWT密钥过于简单且与其他密码相同

**建议：**
- 使用至少32字节的随机密钥
- JWT密钥应该独立，不与数据库密码相同
- 定期轮换密钥（实现刷新token机制）

### 3. **配置文件管理**
```bash
# 添加到 .gitignore
config/conf.yaml
config/*.yaml
!config/conf.example.yaml
.env
```

---

## ⚡ 性能优化

### 1. **数据库连接池优化**

**问题：** MongoDB和Redis连接池配置不当

**当前代码（`internal/ioc/mongo.go`）：**
```go
client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
// 没有设置连接池参数
```

**优化建议：**
```go
func InitMongodb() *mongo.Client {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cfg := InitViper()
    uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
        cfg.Mongodb.Account,
        cfg.Mongodb.Password,
        cfg.Mongodb.Address1,
        cfg.Mongodb.Port1)

    clientOpts := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).           // 最大连接数
        SetMinPoolSize(10).            // 最小连接数
        SetMaxConnIdleTime(30 * time.Second). // 空闲连接超时
        SetConnectTimeout(10 * time.Second).  // 连接超时
        SetServerSelectionTimeout(10 * time.Second)

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return nil, fmt.Errorf("MongoDB连接失败: %w", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("MongoDB Ping失败: %w", err)
    }

    return client
}
```

**Redis连接池优化（`internal/ioc/redis.go`）：**
```go
func InitRedis() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:            viper.GetString("redis.addr"),
        Password:        viper.GetString("redis.password"),
        DB:              0,
        PoolSize:        50,              // 增加到50
        MinIdleConns:    10,              // 最小空闲连接
        MaxRetries:      3,               // 重试次数
        DialTimeout:     5 * time.Second,
        ReadTimeout:     3 * time.Second,
        WriteTimeout:    3 * time.Second,
        PoolTimeout:     4 * time.Second,
        IdleTimeout:     5 * time.Minute, // 空闲连接超时
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if _, err := client.Ping(ctx).Result(); err != nil {
        log.Printf("⚠️ Redis连接失败: %v", err)
        return nil
    }

    log.Println("✅ Redis连接成功")
    return client
}
```

### 2. **BaseDAO Count查询优化**

**问题：** `FindList` 总是执行 `CountDocuments`，即使不需要分页信息

**当前代码（`internal/repository/dao/base.go:71`）：**
```go
total, err := d.Coll.CountDocuments(ctx, filter)
```

**优化方案：**
```go
// 方案1: 添加参数控制是否查询总数
func (d *BaseDAO[T]) FindList(ctx context.Context, filter bson.M, skip, limit int64, sort bson.D, needCount bool) ([]*T, int64, error) {
    findOpts := options.Find()
    if sort != nil {
        findOpts.SetSort(sort)
    }
    if limit > 0 {
        findOpts.SetLimit(limit)
    }
    if skip > 0 {
        findOpts.SetSkip(skip)
    }

    cursor, err := d.Coll.Find(ctx, filter, findOpts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    var results []*T
    if err = cursor.All(ctx, &results); err != nil {
        return nil, 0, err
    }

    var total int64
    if needCount {
        total, err = d.Coll.CountDocuments(ctx, filter)
        if err != nil {
            return nil, 0, err
        }
    } else {
        total = int64(len(results))
    }

    return results, total, nil
}

// 方案2: 使用EstimatedDocumentCount替代（对于不需要精确总数的场景）
```

### 3. **批量查询优化**

**问题：** `CustomerService.GetList` 中虽然使用了批量查询，但可以进一步优化

**当前代码存在的问题：**
```go
// 收集ID时可能有重复
bedIDs := make([]primitive.ObjectID, 0)
for _, v := range list {
    if !v.BedID.IsZero() {
        bedIDs = append(bedIDs, v.BedID) // 可能重复
    }
}
```

**优化方案：**
```go
// 使用map去重
func (s *CustomerService) GetList(ctx context.Context, query domain.CustomerQuery, skip, limit int64) ([]*domain.CustomerResponse, int64, error) {
    // ... 查询逻辑 ...

    // 使用map去重收集ID
    bedIDMap := make(map[primitive.ObjectID]bool)
    careIDMap := make(map[primitive.ObjectID]bool)
    dietIDMap := make(map[primitive.ObjectID]bool)

    for _, v := range list {
        if !v.BedID.IsZero() {
            bedIDMap[v.BedID] = true
        }
        if !v.CareLevelID.IsZero() {
            careIDMap[v.CareLevelID] = true
        }
        if !v.DietPlanID.IsZero() {
            dietIDMap[v.DietPlanID] = true
        }
    }

    // 转换为切片
    bedIDs := make([]primitive.ObjectID, 0, len(bedIDMap))
    for id := range bedIDMap {
        bedIDs = append(bedIDs, id)
    }

    // 并发查询关联数据
    var bedMap, careMap, dietMap map[primitive.ObjectID]string
    var wg sync.WaitGroup
    wg.Add(3)

    go func() {
        defer wg.Done()
        bedMap = s.getBedMap(ctx, bedIDs)
    }()

    go func() {
        defer wg.Done()
        careMap = s.getCareMap(ctx, careIDs)
    }()

    go func() {
        defer wg.Done()
        dietMap = s.getDietMap(ctx, dietIDs)
    }()

    wg.Wait()

    // ... 组装结果 ...
}
```

### 4. **添加缓存层**

**建议在以下场景使用Redis缓存：**
- 护理级别列表（CareLevel）- 变化频率低
- 膳食计划列表（DietPlan）- 变化频率低
- 房间信息（Room）- 变化频率低
- 热点客户信息查询

**实施示例：**
```go
func (s *CareLevelService) GetList(ctx context.Context) ([]*domain.CareLevel, error) {
    cacheKey := "care_levels:all"
    
    // 1. 尝试从缓存读取
    var result []*domain.CareLevel
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil && cached != "" {
        if err := json.Unmarshal([]byte(cached), &result); err == nil {
            return result, nil
        }
    }

    // 2. 缓存未命中，查询数据库
    result, _, err = s.repo.FindList(ctx, bson.M{}, 0, 0)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存（设置1小时过期）
    if data, err := json.Marshal(result); err == nil {
        s.redis.Set(ctx, cacheKey, data, time.Hour)
    }

    return result, nil
}
```

---

## 🏗️ 架构改进

### 1. **优雅关闭（Graceful Shutdown）**

**问题：** 应用停止时没有优雅关闭机制，可能导致请求中断

**当前代码（`app.go:52`）：**
```go
return app.server.Run(":8081") // 阻塞直到错误
```

**优化方案：**
```go
func (app *App) Start() error {
    // 注册debug路由
    app.RegisterDebugRoute()
    ioc.InitApiColl(app.server, app.mongodb, app.redis)
    app.AlertService.Start()

    // 启动HTTP服务器（非阻塞）
    srv := &http.Server{
        Addr:           ":8081",
        Handler:        app.server,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    go func() {
        log.Println("🚀 服务启动在端口 8081")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("服务启动失败: %v", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("⏳ 正在优雅关闭服务...")

    // 设置5秒超时
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 关闭HTTP服务器
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("服务器强制关闭: %v", err)
    }

    // 关闭数据库连接
    if err := app.mongodb.Disconnect(ctx); err != nil {
        log.Printf("MongoDB断开连接失败: %v", err)
    }

    // 关闭Redis连接
    if err := app.redis.Close(); err != nil {
        log.Printf("Redis断开连接失败: %v", err)
    }

    // 停止报警服务
    app.AlertService.Stop()

    log.Println("✅ 服务已优雅关闭")
    return nil
}
```

### 2. **Context超时控制**

**问题：** 大部分数据库操作没有超时控制

**优化建议：**
```go
// 在middleware中添加请求超时控制
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        c.Request = c.Request.WithContext(ctx)
        
        finished := make(chan struct{})
        go func() {
            c.Next()
            finished <- struct{}{}
        }()

        select {
        case <-ctx.Done():
            c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
                "code": 408,
                "msg":  "请求超时",
            })
        case <-finished:
        }
    }
}

// 在gin.go中应用
server.Use(TimeoutMiddleware(30 * time.Second))
```

### 3. **添加限流中间件**

**建议使用令牌桶算法限流：**
```go
// 基于Redis的分布式限流
func RateLimitMiddleware(redis *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
        
        pipe := redis.Pipeline()
        incr := pipe.Incr(c.Request.Context(), key)
        pipe.Expire(c.Request.Context(), key, window)
        _, err := pipe.Exec(c.Request.Context())
        
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "code": 500,
                "msg":  "限流检查失败",
            })
            return
        }

        if incr.Val() > int64(limit) {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "code": 429,
                "msg":  "请求过于频繁，请稍后重试",
            })
            return
        }

        c.Next()
    }
}
```

### 4. **统一错误处理**

**创建统一的错误类型：**
```go
// internal/domain/errors.go
package domain

import "net/http"

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

var (
    ErrInvalidRequest   = &AppError{Code: 400, Message: "请求参数错误"}
    ErrUnauthorized     = &AppError{Code: 401, Message: "未授权"}
    ErrForbidden        = &AppError{Code: 403, Message: "禁止访问"}
    ErrNotFound         = &AppError{Code: 404, Message: "资源不存在"}
    ErrConflict         = &AppError{Code: 409, Message: "资源冲突"}
    ErrInternalServer   = &AppError{Code: 500, Message: "服务器内部错误"}
)

// 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            if appErr, ok := err.(*AppError); ok {
                c.JSON(appErr.Code, appErr)
            } else {
                c.JSON(http.StatusInternalServerError, &AppError{
                    Code:    500,
                    Message: "服务器内部错误",
                    Details: err.Error(),
                })
            }
        }
    }
}
```

---

## 📝 代码质量

### 1. **消除print调试代码**

**问题位置：** `internal/web/handler.go:100`
```go
print("err") // ❌ 应该使用日志
```

**修复：**
```go
log.Printf("患者注册失败: %v", err)
```

### 2. **完成TODO项**

**发现的TODO：**
- `app.go:56` - 定时任务处理信息
- `internal/domain/User.go:282` - Service结构体的Duration字段

**建议：** 尽快完成或删除过时的TODO注释

### 3. **重构重复代码**

**问题：** `CustomerService` 中床位释放逻辑重复

**优化建议：**
```go
// 提取公共方法
func (s *CustomerService) releaseBedIfExists(ctx context.Context, bedID primitive.ObjectID) error {
    if bedID.IsZero() {
        return nil
    }
    return s.bedRepo.Release(ctx, bedID)
}

func (s *CustomerService) assignBed(ctx context.Context, bedID, customerID primitive.ObjectID) error {
    bed, err := s.bedRepo.FindById(ctx, bedID)
    if err != nil {
        return err
    }
    if bed == nil {
        return errors.New("床位不存在")
    }
    if bed.Status == "占用" {
        return errors.New("该床位已被占用")
    }
    return s.bedRepo.AssignToCustomer(ctx, bedID, customerID)
}
```

### 4. **添加接口抽象**

**建议定义Repository接口：**
```go
// internal/repository/interface.go
package repository

type CustomerRepository interface {
    Create(ctx context.Context, customer *domain.Customer) error
    FindById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error)
    FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Customer, int64, error)
    Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
    Delete(ctx context.Context, id primitive.ObjectID) error
}
```

---

## 🗄️ 数据库优化

### 1. **添加数据库索引**

**创建索引初始化脚本：**
```go
// internal/repository/dao/indexes.go
package dao

func InitIndexes(db *mongo.Database) error {
    ctx := context.Background()

    // Customer集合索引
    customerColl := db.Collection("customers")
    _, err := customerColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true)},
        {Keys: bson.D{{Key: "id_card", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
        {Keys: bson.D{{Key: "status", Value: 1}}},
        {Keys: bson.D{{Key: "bed_id", Value: 1}}},
        {Keys: bson.D{{Key: "care_level_id", Value: 1}}},
        {Keys: bson.D{{Key: "health_manager_id", Value: 1}}},
        {Keys: bson.D{{Key: "created_at", Value: -1}}}, // 创建时间倒序
    })
    if err != nil {
        return fmt.Errorf("创建Customer索引失败: %w", err)
    }

    // User集合索引
    userColl := db.Collection("users")
    _, err = userColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true)},
    })
    if err != nil {
        return fmt.Errorf("创建User索引失败: %w", err)
    }

    // Bed集合索引
    bedColl := db.Collection("beds")
    _, err = bedColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "room_id", Value: 1}}},
        {Keys: bson.D{{Key: "status", Value: 1}}},
        {Keys: bson.D{{Key: "customer_id", Value: 1}}},
    })
    if err != nil {
        return fmt.Errorf("创建Bed索引失败: %w", err)
    }

    // Record集合索引
    recordColl := db.Collection("records")
    _, err = recordColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {Keys: bson.D{{Key: "customer_id", Value: 1}}},
        {Keys: bson.D{{Key: "type", Value: 1}}},
        {Keys: bson.D{{Key: "created_at", Value: -1}}},
        {Keys: bson.D{{Key: "customer_id", Value: 1}, {Key: "created_at", Value: -1}}}, // 复合索引
    })
    if err != nil {
        return fmt.Errorf("创建Record索引失败: %w", err)
    }

    log.Println("✅ 数据库索引创建成功")
    return nil
}
```

**在应用启动时调用：**
```go
// app.go
func (app *App) Start() error {
    // 初始化数据库索引
    if err := dao.InitIndexes(app.mongodb.Database("kongdong")); err != nil {
        log.Printf("⚠️ 索引初始化失败: %v", err)
    }
    // ...
}
```

### 2. **数据库健康检查**

```go
// internal/web/health.go
func (h *HealthHandler) CheckHealth(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    health := map[string]string{
        "status": "healthy",
    }

    // 检查MongoDB
    if err := h.mongodb.Ping(ctx, nil); err != nil {
        health["mongodb"] = "unhealthy"
        health["status"] = "degraded"
    } else {
        health["mongodb"] = "healthy"
    }

    // 检查Redis
    if _, err := h.redis.Ping(ctx).Result(); err != nil {
        health["redis"] = "unhealthy"
        health["status"] = "degraded"
    } else {
        health["redis"] = "healthy"
    }

    status := http.StatusOK
    if health["status"] != "healthy" {
        status = http.StatusServiceUnavailable
    }

    c.JSON(status, health)
}
```

---

## 📊 可观测性

### 1. **结构化日志**

**使用zap进行结构化日志：**
```go
// internal/ioc/logger.go
package ioc

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    logger, err := config.Build()
    if err != nil {
        panic(err)
    }
    
    return logger
}

// 使用示例
logger.Info("用户登录成功",
    zap.String("user_id", userID),
    zap.String("ip", clientIP),
    zap.Duration("duration", time.Since(start)),
)
```

### 2. **添加Trace ID**

```go
// 中间件添加Trace ID
func TraceIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = uuid.New().String()
        }
        c.Set("trace_id", traceID)
        c.Header("X-Trace-ID", traceID)
        c.Next()
    }
}
```

### 3. **完善Prometheus指标**

```go
// internal/ioc/prometheus.go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "collection"},
    )

    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        
        c.Next()

        duration := time.Since(start).Seconds()
        status := fmt.Sprintf("%d", c.Writer.Status())
        
        httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
    }
}
```

---

## 🌐 API设计

### 1. **统一响应格式**

**创建统一的Response结构：**
```go
// internal/web/response.go
package web

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    traceID, _ := c.Get("trace_id")
    c.JSON(http.StatusOK, Response{
        Code:    200,
        Message: "success",
        Data:    data,
        TraceID: traceID.(string),
    })
}

func Error(c *gin.Context, code int, message string) {
    traceID, _ := c.Get("trace_id")
    c.JSON(code, Response{
        Code:    code,
        Message: message,
        TraceID: traceID.(string),
    })
}

func BadRequest(c *gin.Context, message string) {
    Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
    Error(c, http.StatusUnauthorized, message)
}

func InternalError(c *gin.Context, message string) {
    Error(c, http.StatusInternalServerError, message)
}
```

### 2. **修正HTTP状态码**

**问题：** 很多错误情况返回200状态码

**修复示例（`internal/web/customer.go`）：**
```go
// 修改前
c.JSON(http.StatusOK, gin.H{
    "code": 400,
    "msg":  "请求参数错误: " + err.Error(),
})

// 修改后
c.JSON(http.StatusBadRequest, gin.H{
    "code": 400,
    "msg":  "请求参数错误: " + err.Error(),
})
```

### 3. **输入验证中间件**

```go
// 使用validator库进行验证
type CustomerCreateRequest struct {
    Name   string `json:"name" binding:"required,min=2,max=50"`
    Age    int    `json:"age" binding:"required,gte=0,lte=150"`
    Gender string `json:"gender" binding:"required,oneof=男 女"`
    Phone  string `json:"phone" binding:"required,len=11,numeric"`
    IDCard string `json:"id_card" binding:"omitempty,len=18"`
}
```

---

## ⚙️ 配置管理

### 1. **环境区分**

**创建多环境配置：**
```
config/
  ├── conf.dev.yaml      # 开发环境
  ├── conf.staging.yaml  # 测试环境
  ├── conf.prod.yaml     # 生产环境
  └── conf.example.yaml  # 示例配置（提交到Git）
```

**修改viper初始化：**
```go
func InitViper() *Config {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "dev"
    }

    configFile := fmt.Sprintf("config/conf.%s.yaml", env)
    viper.SetConfigFile(configFile)

    // 支持环境变量覆盖
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")
    
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("配置文件读取失败: %v", err)
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        log.Fatalf("配置解析失败: %v", err)
    }

    return &config
}
```

### 2. **配置验证**

```go
func (c *Config) Validate() error {
    if c.Mongodb.Address1 == "" {
        return errors.New("MongoDB地址不能为空")
    }
    if c.Mongodb.Port1 == 0 {
        return errors.New("MongoDB端口不能为空")
    }
    if c.General.JWT == "" || len(c.General.JWT) < 32 {
        return errors.New("JWT密钥至少需要32字符")
    }
    return nil
}
```

---

## 🚀 其他建议

### 1. **Dockerfile优化**

**当前问题：** 暴露端口8080但实际运行在8081

**修复：**
```dockerfile
# 修改
EXPOSE 8080
# 改为
EXPOSE 8081

# 或者使用环境变量
ENV PORT=8081
EXPOSE ${PORT}
```

### 2. **添加健康检查**

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1
```

### 3. **Git Hooks**

**添加pre-commit检查：**
```bash
#!/bin/bash
# .git/hooks/pre-commit

# 检查是否包含敏感信息
if git diff --cached | grep -E "password.*=|jwt.*=|secret.*=" ; then
    echo "⚠️  警告: 提交内容可能包含敏感信息!"
    exit 1
fi

# 运行格式化
go fmt ./...

# 运行lint
golangci-lint run
```

---

## 📋 优先级总结

### 🔴 高优先级（立即处理）
1. ✅ 移除敏感信息，使用环境变量
2. ✅ 修复HTTP状态码不一致问题
3. ✅ 添加优雅关闭机制
4. ✅ 添加数据库索引
5. ✅ 替换panic/log.Fatal为优雅错误处理

### 🟡 中优先级（近期处理）
1. ✅ 优化数据库连接池
2. ✅ 添加缓存层
3. ✅ 实现结构化日志
4. ✅ 添加请求限流
5. ✅ 完善Prometheus指标

### 🟢 低优先级（逐步改进）
1. ✅ 重构重复代码
2. ✅ 添加单元测试
3. ✅ 完善API文档
4. ✅ 添加集成测试
5. ✅ 代码风格统一

---

## 📖 参考资源

- [Go最佳实践](https://github.com/golang-standards/project-layout)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Redis Go Client](https://redis.uptrace.dev/)
- [Prometheus Go Client](https://prometheus.io/docs/guides/go-application/)

---

**生成时间：** 2025-01-08  
**版本：** v1.0
