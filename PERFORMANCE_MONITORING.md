# 性能监控与优化指南

## 📊 当前性能指标

### 数据库连接池监控

```go
// internal/web/metrics.go
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type MetricsHandler struct {
	mongodb *mongo.Client
	redis   *redis.Client
}

func NewMetricsHandler(mongodb *mongo.Client, redis *redis.Client) *MetricsHandler {
	return &MetricsHandler{
		mongodb: mongodb,
		redis:   redis,
	}
}

func (h *MetricsHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/debug/stats", h.GetStats)
}

func (h *MetricsHandler) GetStats(c *gin.Context) {
	stats := make(map[string]interface{})

	// Redis连接池统计
	poolStats := h.redis.PoolStats()
	stats["redis"] = map[string]interface{}{
		"hits":        poolStats.Hits,
		"misses":      poolStats.Misses,
		"timeouts":    poolStats.Timeouts,
		"total_conns": poolStats.TotalConns,
		"idle_conns":  poolStats.IdleConns,
		"stale_conns": poolStats.StaleConns,
	}

	// MongoDB连接数（需要通过serverStatus命令获取）
	// 这需要管理员权限，这里提供示例
	/*
	var result bson.M
	err := h.mongodb.Database("admin").RunCommand(
		context.Background(),
		bson.D{{Key: "serverStatus", Value: 1}},
	).Decode(&result)
	if err == nil {
		connections := result["connections"].(bson.M)
		stats["mongodb"] = connections
	}
	*/

	c.JSON(http.StatusOK, stats)
}
```

## ⚡ 性能优化建议

### 1. 查询性能优化

#### 使用Explain分析查询

```go
// 示例：分析Customer查询性能
func (d *CustomerDAO) AnalyzeQuery(ctx context.Context, filter bson.M) {
	cursor, err := d.Coll.Find(ctx, filter, options.Find().SetComment("performance_test"))
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// 获取查询计划
	explainResult := d.Coll.Database().RunCommand(
		ctx,
		bson.D{
			{Key: "explain", Value: bson.D{
				{Key: "find", Value: d.Coll.Name()},
				{Key: "filter", Value: filter},
			}},
		},
	)

	var result bson.M
	if err := explainResult.Decode(&result); err == nil {
		log.Printf("查询计划: %+v", result)
	}
}
```

### 2. 缓存策略

#### 实现多级缓存

```go
// internal/service/cached_care_level.go
package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
)

type CachedCareLevelService struct {
	repo       repository.CareLevelRepository
	redis      *redis.Client
	localCache *sync.Map // 本地缓存
}

func NewCachedCareLevelService(repo repository.CareLevelRepository, redis *redis.Client) *CachedCareLevelService {
	return &CachedCareLevelService{
		repo:       repo,
		redis:      redis,
		localCache: &sync.Map{},
	}
}

func (s *CachedCareLevelService) GetList(ctx context.Context) ([]*domain.CareLevel, error) {
	cacheKey := "care_levels:all"

	// 1. 尝试从本地缓存获取（最快）
	if cached, ok := s.localCache.Load(cacheKey); ok {
		if data, ok := cached.(*cacheItem); ok && time.Now().Before(data.expiresAt) {
			return data.value.([]*domain.CareLevel), nil
		}
	}

	// 2. 尝试从Redis获取
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var result []*domain.CareLevel
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			// 更新本地缓存
			s.localCache.Store(cacheKey, &cacheItem{
				value:     result,
				expiresAt: time.Now().Add(5 * time.Minute),
			})
			return result, nil
		}
	}

	// 3. 从数据库查询
	result, _, err := s.repo.FindList(ctx, bson.M{}, 0, 0)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存（1小时）
	if data, err := json.Marshal(result); err == nil {
		s.redis.Set(ctx, cacheKey, data, time.Hour)
	}

	// 5. 更新本地缓存（5分钟）
	s.localCache.Store(cacheKey, &cacheItem{
		value:     result,
		expiresAt: time.Now().Add(5 * time.Minute),
	})

	return result, nil
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}
```

### 3. 批量操作优化

#### 使用MongoDB的BulkWrite

```go
// 批量更新床位状态
func (d *BedDAO) BatchUpdateStatus(ctx context.Context, updates map[primitive.ObjectID]string) error {
	if len(updates) == 0 {
		return nil
	}

	var models []mongo.WriteModel
	for id, status := range updates {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"status":     status,
					"updated_at": time.Now(),
				},
			})
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false) // 并发执行
	_, err := d.Coll.BulkWrite(ctx, models, opts)
	return err
}
```

### 4. 并发查询优化

#### 使用errgroup管理并发

```go
import "golang.org/x/sync/errgroup"

func (s *CustomerService) GetListOptimized(ctx context.Context, query domain.CustomerQuery, skip, limit int64) ([]*domain.CustomerResponse, int64, error) {
	// 查询客户列表
	list, total, err := s.customerRepo.FindList(ctx, buildFilter(query), skip, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(list) == 0 {
		return []*domain.CustomerResponse{}, total, nil
	}

	// 收集关联ID（去重）
	bedIDs, careIDs, dietIDs := collectUniqueIDs(list)

	// 使用errgroup并发查询
	var bedMap, careMap, dietMap map[primitive.ObjectID]string
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		bedMap = s.getBedMap(ctx, bedIDs)
		return nil
	})

	g.Go(func() error {
		careMap = s.getCareMap(ctx, careIDs)
		return nil
	})

	g.Go(func() error {
		dietMap = s.getDietMap(ctx, dietIDs)
		return nil
	})

	// 等待所有查询完成
	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	// 组装结果
	return assembleResults(list, bedMap, careMap, dietMap), total, nil
}

// 辅助函数：收集唯一ID
func collectUniqueIDs(list []*domain.Customer) ([]primitive.ObjectID, []primitive.ObjectID, []primitive.ObjectID) {
	bedIDSet := make(map[primitive.ObjectID]bool)
	careIDSet := make(map[primitive.ObjectID]bool)
	dietIDSet := make(map[primitive.ObjectID]bool)

	for _, v := range list {
		if !v.BedID.IsZero() {
			bedIDSet[v.BedID] = true
		}
		if !v.CareLevelID.IsZero() {
			careIDSet[v.CareLevelID] = true
		}
		if !v.DietPlanID.IsZero() {
			dietIDSet[v.DietPlanID] = true
		}
	}

	return setToSlice(bedIDSet), setToSlice(careIDSet), setToSlice(dietIDSet)
}

func setToSlice(set map[primitive.ObjectID]bool) []primitive.ObjectID {
	slice := make([]primitive.ObjectID, 0, len(set))
	for id := range set {
		slice = append(slice, id)
	}
	return slice
}
```

## 📈 性能测试

### 使用pprof进行性能分析

```go
// main.go 添加pprof支持
import (
	_ "net/http/pprof"
	"net/http"
)

func main() {
	// 启动pprof服务器（在另一个端口）
	go func() {
		log.Println("pprof服务启动在 http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// ... 其他启动代码
}
```

### 性能测试命令

```bash
# 1. CPU性能分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 2. 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 3. Goroutine分析
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 4. 生成火焰图
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
```

### 压力测试

使用wrk或hey进行压力测试：

```bash
# 安装hey
go install github.com/rakyll/hey@latest

# 测试健康检查端点
hey -n 10000 -c 100 http://localhost:8081/health

# 测试API端点（带认证）
hey -n 1000 -c 50 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8081/api/customers
```

## 🎯 性能目标

### 响应时间目标

| 端点类型 | P50 | P95 | P99 |
|---------|-----|-----|-----|
| 健康检查 | <10ms | <20ms | <50ms |
| 简单查询 | <50ms | <100ms | <200ms |
| 复杂查询 | <200ms | <500ms | <1s |
| 批量操作 | <500ms | <1s | <2s |

### 资源使用目标

- **CPU**: 平均 < 50%，峰值 < 80%
- **内存**: 平均 < 500MB，峰值 < 1GB
- **数据库连接**: 平均 < 20，峰值 < 50
- **Goroutines**: 平均 < 100，峰值 < 500

### 监控指标

```go
// 使用Prometheus监控关键指标
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 请求延迟
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// 数据库查询延迟
	dbDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"collection", "operation"},
	)

	// 缓存命中率
	cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits.",
		},
		[]string{"cache_type"},
	)

	cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses.",
		},
		[]string{"cache_type"},
	)
)

// 使用中间件记录指标
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpDuration.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Observe(duration)
	}
}
```

## 📊 监控面板

### Grafana Dashboard配置

创建 `monitoring/grafana-dashboard.json`，包含以下面板：

1. **请求速率** (QPS)
2. **响应时间** (P50, P95, P99)
3. **错误率** (4xx, 5xx)
4. **数据库连接数**
5. **缓存命中率**
6. **Goroutine数量**
7. **内存使用**
8. **CPU使用率**

### Prometheus配置

```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'classroom-analysis'
    static_configs:
      - targets: ['localhost:8081']
    metrics_path: '/metrics'
```

## 🔧 优化检查清单

- [ ] 数据库查询使用了索引
- [ ] 热点数据已缓存
- [ ] 避免了N+1查询
- [ ] 使用了连接池
- [ ] 并发操作使用了goroutine
- [ ] 有适当的超时控制
- [ ] 大数据集使用了分页
- [ ] 避免了全表扫描
- [ ] 定期清理过期数据
- [ ] 监控了关键性能指标

---

**最后更新：** 2025-01-08
