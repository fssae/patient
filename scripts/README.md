# MongoDB 数据填充脚本使用说明

## 概述

`seed_data.go` 是一个用于填充 MongoDB 数据库的脚本，它会为养老院管理系统创建完整的测试数据。

## 数据集合

脚本会填充以下 MongoDB 集合：

1. **users** - 5个注册用户
2. **rooms** - 6个房间（A101, A102, A201, A202, B101, B201）
3. **beds** - 14个床位
4. **diet_plans** - 4种膳食计划（普通餐、低糖餐、流质餐、低盐餐）
5. **care_levels** - 4个护理级别（特级、一级、二级、三级）
6. **health_managers** - 5位健康管家
7. **customers** - 5位入住老人
8. **services** - 7个服务项目
9. **customer_services** - 客户购买的服务记录
10. **records** - 入住/外出登记记录
11. **care_records** - 护理记录（每位老人7天的护理记录）
12. **patients** - 5位患者（用于视频监控系统）
13. **analysis** - 分析日志（跌倒检测、情绪异常等告警）

## 使用方法

### 方法一：直接运行

```bash
cd d:\code\kongdong\scripts
go run seed_data.go
```

### 方法二：编译后运行

```bash
cd d:\code\kongdong\scripts
go build -o seed_data.exe seed_data.go
./seed_data.exe
```

## 配置说明

脚本中的数据库连接配置：

```go
const (
    MongoURI     = "mongodb://admin:zjh770910@82.156.64.69:27017"
    DatabaseName = "kongdong"
)
```

如需修改，请编辑 `seed_data.go` 文件中的这些常量。

## 数据详情

### 用户数据
- 5个注册用户，手机号从 13800138001 到 13800138005
- 默认密码：`password123`（已加密）

### 房间和床位
- A101: 双人间（2床）
- A102: 单人间（1床）
- A201: 多人间（4床）
- A202: 双人间（2床）
- B101: 单人间（1床，维护中）
- B201: 双人间（2床）

### 入住老人
1. **张老伯** (78岁) - 床位A101-1 - 低糖餐 - 二级护理
2. **李奶奶** (82岁) - 床位A101-2 - 流质餐 - 一级护理
3. **王老先生** (75岁) - 床位A102-1 - 普通餐 - 三级护理
4. **刘奶奶** (80岁) - 床位A201-1 - 低盐餐 - 二级护理
5. **陈老伯** (76岁) - 床位A201-2 - 低糖餐 - 三级护理

### 健康管家
1. 王医生 - 内科
2. 李护士 - 老年护理
3. 张医生 - 心血管科
4. 刘护士 - 康复护理
5. 陈医生 - 神经内科

### 服务项目
- 康复理疗 (200元/次)
- 定期体检 (500元/月)
- 洗衣服务 (300元/月)
- 理发服务 (50元/次)
- 文娱活动 (200元/月)
- 心理咨询 (300元/次)
- 陪同就医 (150元/次)

### 告警数据
- 5条跌倒检测事件
- 3条情绪异常事件
- 3条会话完成记录

## 注意事项

1. **数据重复问题**：脚本不会检查数据是否已存在，重复运行会插入重复数据
2. **清空数据库**：如需重新填充，请先手动清空相关集合
3. **网络连接**：确保能够连接到 MongoDB 服务器（82.156.64.69:27017）
4. **依赖包**：需要安装 `go.mongodb.org/mongo-driver` 包

## 清空数据库（可选）

如果需要清空数据库后重新填充，可以使用 MongoDB Shell：

```bash
mongosh "mongodb://admin:zjh770910@82.156.64.69:27017/kongdong"
```

然后执行：

```javascript
db.users.deleteMany({})
db.rooms.deleteMany({})
db.beds.deleteMany({})
db.customers.deleteMany({})
db.diet_plans.deleteMany({})
db.care_levels.deleteMany({})
db.health_managers.deleteMany({})
db.services.deleteMany({})
db.customer_services.deleteMany({})
db.records.deleteMany({})
db.care_records.deleteMany({})
db.patients.deleteMany({})
db.analysis.deleteMany({})
```

## 验证数据

填充完成后，可以通过以下方式验证：

```bash
mongosh "mongodb://admin:zjh770910@82.156.64.69:27017/kongdong"
```

```javascript
// 查看各集合的文档数量
db.users.countDocuments()
db.rooms.countDocuments()
db.beds.countDocuments()
db.customers.countDocuments()
db.diet_plans.countDocuments()
db.care_levels.countDocuments()
db.health_managers.countDocuments()
db.services.countDocuments()
db.customer_services.countDocuments()
db.records.countDocuments()
db.care_records.countDocuments()
db.patients.countDocuments()
db.analysis.countDocuments()

// 查看具体数据
db.customers.find().pretty()
db.analysis.find().pretty()
```

## 故障排除

### 连接失败
- 检查 MongoDB 服务是否运行
- 检查网络连接
- 验证用户名和密码是否正确

### 插入失败
- 检查 MongoDB 用户权限
- 查看错误日志获取详细信息

## 扩展

如需添加更多数据或修改数据内容，请编辑 `seed_data.go` 文件中相应的 `seed*` 函数。
