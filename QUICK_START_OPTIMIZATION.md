# 快速开始 - 优化实施指南

本指南帮助你快速实施最关键的优化建议。

## 🚨 立即行动（15分钟）

### 1. 保护敏感信息

```bash
# 1.1 备份当前配置文件
cp config/conf.yaml config/conf.yaml.backup

# 1.2 复制示例配置
cp config/conf.example.yaml config/conf.yaml

# 1.3 创建 .env 文件
cp .env.example .env

# 1.4 编辑 .env 填入真实密钥
vim .env  # 或使用你喜欢的编辑器

# 1.5 确保敏感文件不被提交
git rm --cached config/conf.yaml
git add .gitignore
git commit -m "chore: 移除敏感配置文件，添加环境变量支持"
```

### 2. 修复Dockerfile端口不匹配

编辑 `Dockerfile` 第36行：
```dockerfile
# 修改前
EXPOSE 8080

# 修改后
EXPOSE 8081
```

### 3. 安装必要的依赖

```bash
# 如果需要环境变量支持
go get github.com/joho/godotenv
```

---

## ⚡ 快速胜利（1小时内）

### 1. 添加数据库索引

创建文件 `internal/repository/dao/indexes.go`，复制以下代码：

```go
package dao

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Customer集合索引
	customerColl := db.Collection("customers")
	_, err := customerColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "phone", Value: 1}}},
		{Keys: bson.D{{Key: "id_card", Value: 1}}, Options: options.Index().SetSparse(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "bed_id", Value: 1}}},
	})
	if err != nil {
		log.Printf("创建Customer索引失败: %v", err)
		return err
	}

	// User集合索引
	userColl := db.Collection("users")
	_, err = userColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		log.Printf("创建User索引失败: %v", err)
		return err
	}

	log.Println("✅ 数据库索引创建成功")
	return nil
}
```

在 `app.go` 的 `Start` 方法中添加：

```go
import "classroom-analysis/internal/repository/dao"

func (app *App) Start() error {
	// 初始化数据库索引
	db := app.mongodb.Database("kongdong")
	if err := dao.InitIndexes(db); err != nil {
		log.Printf("⚠️ 索引初始化失败: %v", err)
	}
	
	// ... 其他代码
}
```

### 2. 优化MongoDB连接池

编辑 `internal/ioc/mongo.go`：

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

	// 🆕 添加连接池配置
	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ MongoDB连接成功")
	return client
}
```

### 3. 优化Redis连接池

编辑 `internal/ioc/redis.go`：

```go
func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           0,
		PoolSize:     50,              // 🆕 增加连接池大小
		MinIdleConns: 10,              // 🆕 最小空闲连接
		MaxRetries:   3,               // 🆕 重试次数
		DialTimeout:  5 * time.Second, // 🆕 连接超时
		ReadTimeout:  3 * time.Second, // 🆕 读超时
		WriteTimeout: 3 * time.Second, // 🆕 写超时
		PoolTimeout:  4 * time.Second, // 🆕 连接池超时
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️ Redis连接失败: %v", err)
	} else {
		log.Println("✅ Redis连接成功")
	}

	return client
}
```

### 4. 修正HTTP状态码

查找并替换所有不正确的状态码：

```bash
# 查找所有返回200但code不是200的地方
grep -r "StatusOK.*code.*4[0-9][0-9]" internal/web/

# 示例修复（在每个handler中）
# 修改前：
c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "错误"})

# 修改后：
c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "错误"})
```

---

## 📅 本周完成（2-4小时）

### 1. 添加优雅关闭

编辑 `app.go`，替换 `Start` 方法：

```go
import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (app *App) Start() error {
	// 初始化
	app.RegisterDebugRoute()
	db := app.mongodb.Database("kongdong")
	if err := dao.InitIndexes(db); err != nil {
		log.Printf("⚠️ 索引初始化失败: %v", err)
	}
	ioc.InitApiColl(app.server, app.mongodb, app.redis)
	app.AlertService.Start()

	// 🆕 创建HTTP服务器
	srv := &http.Server{
		Addr:           ":8081",
		Handler:        app.server,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 🆕 在goroutine中启动服务器
	go func() {
		log.Println("🚀 服务启动在端口 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 🆕 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ 正在优雅关闭服务...")

	// 🆕 设置5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🆕 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	}

	// 🆕 关闭数据库连接
	if err := app.mongodb.Disconnect(ctx); err != nil {
		log.Printf("MongoDB断开连接失败: %v", err)
	}

	// 🆕 关闭Redis连接
	if err := app.redis.Close(); err != nil {
		log.Printf("Redis断开连接失败: %v", err)
	}

	// 🆕 停止报警服务
	app.AlertService.Stop()

	log.Println("✅ 服务已优雅关闭")
	return nil
}
```

### 2. 添加健康检查端点

创建 `internal/web/health.go`：

```go
package web

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	mongodb *mongo.Client
	redis   *redis.Client
}

func NewHealthHandler(mongodb *mongo.Client, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		mongodb: mongodb,
		redis:   redis,
	}
}

func (h *HealthHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/health", h.CheckHealth)
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	health := map[string]interface{}{
		"status": "healthy",
	}

	// 检查MongoDB
	if err := h.mongodb.Ping(ctx, nil); err != nil {
		health["mongodb"] = "unhealthy: " + err.Error()
		health["status"] = "degraded"
	} else {
		health["mongodb"] = "healthy"
	}

	// 检查Redis
	if _, err := h.redis.Ping(ctx).Result(); err != nil {
		health["redis"] = "unhealthy: " + err.Error()
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

在 `wire.go` 中添加 provider，并在路由注册时添加：

```go
// 注册健康检查路由
healthHandler := web.NewHealthHandler(mongodb, redis)
healthHandler.RegisterRoutes(server)
```

### 3. 移除print调试代码

编辑 `internal/web/handler.go:100`：

```go
// 修改前
print("err")

// 修改后
log.Printf("患者注册失败: %v", err)
```

---

## 🎯 下周计划（需要更多时间）

### 1. 实现统一响应格式

创建 `internal/web/response.go` 并逐步重构所有handler使用统一格式。

### 2. 添加结构化日志

引入zap库，替换所有`log.Printf`为结构化日志。

### 3. 添加限流中间件

基于Redis实现分布式限流。

### 4. 编写单元测试

为核心业务逻辑添加单元测试。

---

## 📋 检查清单

实施完成后，检查以下项：

- [ ] 敏感信息已移除，使用环境变量
- [ ] .gitignore已更新
- [ ] Dockerfile端口已修正
- [ ] 数据库索引已创建
- [ ] MongoDB连接池已优化
- [ ] Redis连接池已优化
- [ ] HTTP状态码已修正
- [ ] 优雅关闭已实现
- [ ] 健康检查端点已添加
- [ ] 调试代码已清理

---

## 🐛 常见问题

### Q: 环境变量如何加载？

A: 有两种方式：

1. **使用系统环境变量：**
```bash
export JWT_SECRET="your_secret_key"
go run main.go
```

2. **使用godotenv库（推荐）：**

在 `main.go` 开头添加：
```go
import _ "github.com/joho/godotenv/autoload"
```

或在 `setupEnvironment()` 中添加：
```go
import "github.com/joho/godotenv"

func setupEnvironment() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("没有找到.env文件，使用系统环境变量")
	}
	
	ioc.TimezoneInit()
	ioc.InitViper()
	ioc.InitPrometheus()
}
```

### Q: 如何在viper中使用环境变量？

A: 在 `internal/ioc/viper.go` 中添加：
```go
func InitViper() *Config {
	viper.SetConfigFile("config/conf.yaml")
	
	// 🆕 允许环境变量覆盖配置文件
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	
	// 环境变量会自动替换 ${VAR_NAME} 格式
	viper.SetConfigType("yaml")
	
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}
	
	// ... rest of code
}
```

### Q: 索引创建失败怎么办？

A: 可能是索引已存在或权限不足。可以：

1. 检查MongoDB用户权限
2. 手动删除冲突的索引
3. 在日志中添加更详细的错误信息

---

## 📞 需要帮助？

参考详细文档：`OPTIMIZATION_RECOMMENDATIONS.md`

---

**最后更新：** 2025-01-08
