# 优化实施检查清单 ✅

使用本清单追踪优化进度。在完成每一项后打勾。

---

## 🔴 关键问题（立即处理）

### 安全加固
- [ ] 已将 `config/conf.yaml` 添加到 `.gitignore`
- [ ] 已创建 `.env` 文件并填入真实配置
- [ ] 已从Git历史中移除敏感信息 (`git rm --cached config/conf.yaml`)
- [ ] JWT密钥已更改为至少32字符的强密钥
- [ ] 运行 `./scripts/check-secrets.sh` 通过
- [ ] 确认 `.env` 和 `config/conf.yaml` 未被Git跟踪

### HTTP状态码修正
- [ ] 修改 `internal/web/customer.go` 中的状态码
- [ ] 修改 `internal/web/user.go` 中的状态码  
- [ ] 修改 `internal/web/handler.go` 中的状态码
- [ ] 修改其他handler文件中的状态码
- [ ] 验证API响应状态码与实际情况一致

### 基础修复
- [ ] 修复 `Dockerfile` 端口（改为8081）
- [ ] 移除 `internal/web/handler.go:100` 的 `print("err")`
- [ ] 检查并移除所有其他调试代码

---

## 🟡 重要优化（本周完成）

### 数据库优化
- [ ] 创建 `internal/repository/dao/indexes.go`
- [ ] 实现索引初始化函数
- [ ] 在 `app.go` 中调用索引初始化
- [ ] 优化 `internal/ioc/mongo.go` 连接池配置
- [ ] 优化 `internal/ioc/redis.go` 连接池配置
- [ ] 测试数据库连接池是否生效

### 优雅关闭
- [ ] 修改 `app.go` 的 `Start` 方法实现优雅关闭
- [ ] 添加信号处理（SIGINT, SIGTERM）
- [ ] 实现5秒超时关闭逻辑
- [ ] 确保MongoDB连接正确关闭
- [ ] 确保Redis连接正确关闭
- [ ] 确保AlertService正确停止
- [ ] 测试重启时无请求丢失

### 健康检查
- [ ] 创建 `internal/web/health.go`
- [ ] 实现 `HealthHandler` 结构体
- [ ] 实现 `CheckHealth` 方法（检查MongoDB和Redis）
- [ ] 在路由中注册 `/health` 端点
- [ ] 在 `wire.go` 中添加依赖注入
- [ ] 测试健康检查端点可访问
- [ ] 在 `Dockerfile` 中添加 HEALTHCHECK

### BaseDAO优化
- [ ] 修改 `internal/repository/dao/base.go` 的 `FindList` 方法
- [ ] 添加 `needCount` 参数控制是否查询总数
- [ ] 更新所有调用 `FindList` 的地方
- [ ] 测试优化后的查询性能

---

## 🟢 性能提升（逐步实施）

### 缓存实现
- [ ] 安装必要的依赖（如需要）
- [ ] 创建缓存服务包装器
- [ ] 为 CareLevelService 添加缓存
- [ ] 为 DietPlanService 添加缓存  
- [ ] 为 RoomService 添加缓存
- [ ] 实现缓存失效机制
- [ ] 测试缓存命中率

### 批量查询优化
- [ ] 重构 `CustomerService.GetList` 使用map去重
- [ ] 实现并发查询关联数据
- [ ] 使用 `sync.WaitGroup` 或 `errgroup`
- [ ] 测试并发查询正确性
- [ ] 对比优化前后性能

### 性能监控
- [ ] 创建 `internal/web/metrics.go`
- [ ] 实现 `/debug/stats` 端点
- [ ] 完善Prometheus指标
- [ ] 添加性能监控中间件
- [ ] 配置Grafana面板（如使用）

---

## 📊 可观测性

### 结构化日志
- [ ] 评估是否引入zap库
- [ ] 创建 `internal/ioc/logger.go`
- [ ] 初始化全局logger
- [ ] 逐步替换 `log.Printf` 为结构化日志
- [ ] 添加日志级别控制

### Trace ID
- [ ] 创建 Trace ID 中间件
- [ ] 在请求头中传递 Trace ID
- [ ] 在响应中返回 Trace ID
- [ ] 在日志中记录 Trace ID
- [ ] 测试追踪功能

### 统一响应格式
- [ ] 创建 `internal/web/response.go`
- [ ] 定义统一的 `Response` 结构
- [ ] 实现辅助函数（Success, Error等）
- [ ] 逐步重构handler使用统一格式
- [ ] 更新API文档

---

## 🔧 代码质量

### 代码清理
- [ ] 完成或删除所有TODO注释
- [ ] 移除所有调试代码（print/println）
- [ ] 运行 `go fmt ./...`
- [ ] 运行 `golangci-lint run`（如已安装）
- [ ] 修复linter报告的问题

### 重构重复代码
- [ ] 重构 `CustomerService` 中的床位释放逻辑
- [ ] 提取公共的验证逻辑
- [ ] 提取公共的错误处理逻辑
- [ ] 创建辅助函数减少重复

### 接口抽象
- [ ] 创建 `internal/repository/interface.go`
- [ ] 定义Repository接口
- [ ] 评估是否需要Service接口

---

## ⚙️ 配置管理

### 环境配置
- [ ] 创建 `config/conf.dev.yaml`
- [ ] 创建 `config/conf.staging.yaml`（如需要）
- [ ] 创建 `config/conf.prod.yaml`
- [ ] 修改 `internal/ioc/viper.go` 支持多环境
- [ ] 测试环境切换功能

### 配置验证
- [ ] 实现 `Config.Validate()` 方法
- [ ] 验证必填配置项
- [ ] 验证配置项格式
- [ ] 在启动时调用验证

### 环境变量支持
- [ ] 安装 `github.com/joho/godotenv`（如需要）
- [ ] 在 `main.go` 中加载 `.env` 文件
- [ ] 修改配置文件支持 `${VAR}` 语法
- [ ] 测试环境变量覆盖配置文件

---

## 🧪 测试

### 单元测试
- [ ] 为 `CustomerService` 编写测试
- [ ] 为 `UserService` 编写测试
- [ ] 为关键业务逻辑编写测试
- [ ] 运行 `go test ./...` 确保通过
- [ ] 检查测试覆盖率 `go test -cover ./...`

### 集成测试
- [ ] 编写健康检查端点测试
- [ ] 编写API端点测试
- [ ] 编写数据库操作测试
- [ ] 测试优雅关闭流程

### 性能测试
- [ ] 使用hey/wrk进行压力测试
- [ ] 记录优化前性能基准
- [ ] 记录优化后性能数据
- [ ] 对比性能提升百分比

---

## 📚 文档

### API文档
- [ ] 更新Swagger注释
- [ ] 运行 `swag init` 生成文档
- [ ] 验证 `/swagger/index.html` 可访问
- [ ] 确保所有端点都有文档

### 开发文档
- [ ] 阅读 `OPTIMIZATION_RECOMMENDATIONS.md`
- [ ] 阅读 `QUICK_START_OPTIMIZATION.md`
- [ ] 阅读 `PERFORMANCE_MONITORING.md`
- [ ] 更新项目README（如需要）

### 运维文档
- [ ] 编写部署文档
- [ ] 编写配置说明
- [ ] 编写监控配置
- [ ] 编写故障排查指南

---

## 🚀 部署准备

### 容器化
- [ ] 验证Dockerfile构建成功
- [ ] 添加 `.dockerignore` 文件
- [ ] 测试容器运行正常
- [ ] 优化镜像大小（如需要）

### CI/CD
- [ ] 配置构建流程
- [ ] 配置测试流程
- [ ] 配置部署流程
- [ ] 配置回滚流程

### 监控告警
- [ ] 配置Prometheus抓取
- [ ] 配置Grafana面板
- [ ] 配置告警规则
- [ ] 测试告警通知

---

## ✅ 最终验证

### 功能验证
- [ ] 所有API端点正常工作
- [ ] 数据库操作正确
- [ ] 缓存功能正常
- [ ] WebSocket连接稳定
- [ ] Kafka消息正常收发

### 性能验证
- [ ] 响应时间达标（P95 < 300ms）
- [ ] QPS达标（> 1500）
- [ ] 资源使用合理（CPU < 80%, Memory < 1GB）
- [ ] 数据库连接稳定（< 50）

### 安全验证
- [ ] 运行 `./scripts/check-secrets.sh` 通过
- [ ] 无敏感信息泄露
- [ ] JWT认证正常工作
- [ ] 输入验证充分

### 可靠性验证
- [ ] 健康检查端点正常
- [ ] 优雅关闭功能正常
- [ ] 错误处理完善
- [ ] 日志记录完整

---

## 📊 进度追踪

### 总体进度

- 关键问题: __ / 12 项完成
- 重要优化: __ / 23 项完成  
- 性能提升: __ / 14 项完成
- 可观测性: __ / 11 项完成
- 代码质量: __ / 12 项完成
- 配置管理: __ / 8 项完成
- 测试: __ / 12 项完成
- 文档: __ / 9 项完成
- 部署准备: __ / 12 项完成
- 最终验证: __ / 14 项完成

**总计: __ / 127 项完成 (___%)**

---

## 🎯 里程碑

### 第一周结束
- [ ] 所有关键问题已修复
- [ ] 重要优化完成50%以上
- [ ] 测试环境部署成功

### 第二周结束  
- [ ] 所有重要优化已完成
- [ ] 性能提升已实现
- [ ] 监控系统已配置

### 第三周结束
- [ ] 所有优化项已完成
- [ ] 文档已更新
- [ ] 生产环境已部署

---

## 💡 提示

- 每完成一个大类别，进行一次提交
- 定期运行测试确保功能正常
- 遇到问题及时查阅对应的文档
- 优化过程中保持代码可运行
- 记录遇到的问题和解决方案

---

**最后更新**: 2025-01-08  
**使用说明**: 在Markdown编辑器中打开此文件，勾选已完成的项目。
