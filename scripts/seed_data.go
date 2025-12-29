package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"classroom-analysis/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	mongoURI    = "mongodb://admin:zjh770910@82.156.64.69:27017"
	database    = "kongdong"
	defaultPass = "123456" // 默认密码
)

func main() {
	fmt.Println("开始填充MongoDB数据库...")

	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping MongoDB失败: %v", err)
	}

	db := client.Database(database)
	fmt.Println("✓ MongoDB连接成功")

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)
	passwordHash := string(hashedPassword)

	now := time.Now()

	// 1. 填充健康管家
	fmt.Println("\n1. 填充健康管家数据...")
	healthManagers := []interface{}{
		domain.HealthManager{
			ID:        primitive.NewObjectID(),
			Name:      "张医生",
			Phone:     "13800001001",
			Email:     "zhang@example.com",
			Specialty: "老年病科",
			Status:    "在职",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.HealthManager{
			ID:        primitive.NewObjectID(),
			Name:      "李护士",
			Phone:     "13800001002",
			Email:     "li@example.com",
			Specialty: "护理管理",
			Status:    "在职",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.HealthManager{
			ID:        primitive.NewObjectID(),
			Name:      "王营养师",
			Phone:     "13800001003",
			Email:     "wang@example.com",
			Specialty: "营养配餐",
			Status:    "在职",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	healthManagerIDs := insertMany(db, "health_managers", healthManagers)
	fmt.Printf("  ✓ 创建了 %d 个健康管家\n", len(healthManagerIDs))

	// 2. 填充护理级别
	fmt.Println("\n2. 填充护理级别数据...")
	careLevels := []interface{}{
		domain.CareLevel{
			ID:          primitive.NewObjectID(),
			Name:        "一级护理",
			Level:       1,
			Description: "生活完全自理，需要日常健康监测",
			Content:     "每日测量血压、体温，定期体检，提供健康咨询",
			Price:       500.0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.CareLevel{
			ID:          primitive.NewObjectID(),
			Name:        "二级护理",
			Level:       2,
			Description: "生活部分自理，需要协助日常活动",
			Content:     "协助洗漱、穿衣，每日健康监测，定期康复训练",
			Price:       1000.0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.CareLevel{
			ID:          primitive.NewObjectID(),
			Name:        "三级护理",
			Level:       3,
			Description: "生活不能自理，需要全面护理",
			Content:     "24小时护理，协助所有日常活动，医疗监测，康复训练",
			Price:       2000.0,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	careLevelIDs := insertMany(db, "care_levels", careLevels)
	fmt.Printf("  ✓ 创建了 %d 个护理级别\n", len(careLevelIDs))

	// 3. 填充膳食计划
	fmt.Println("\n3. 填充膳食计划数据...")
	dietPlans := []interface{}{
		domain.DietPlan{
			ID:          primitive.NewObjectID(),
			Name:        "普通餐",
			Description: "适合一般老年人的营养均衡餐",
			WeekMenu: []domain.WeekDayMenu{
				{Day: "周一", Breakfast: "小米粥、鸡蛋、小菜", Lunch: "米饭、红烧肉、青菜", Dinner: "面条、炒菜", Snack: "水果"},
				{Day: "周二", Breakfast: "豆浆、包子、小菜", Lunch: "米饭、鱼、青菜", Dinner: "粥、炒菜", Snack: "水果"},
				{Day: "周三", Breakfast: "牛奶、面包、小菜", Lunch: "米饭、鸡肉、青菜", Dinner: "饺子", Snack: "水果"},
				{Day: "周四", Breakfast: "小米粥、鸡蛋、小菜", Lunch: "米饭、排骨、青菜", Dinner: "面条、炒菜", Snack: "水果"},
				{Day: "周五", Breakfast: "豆浆、包子、小菜", Lunch: "米饭、鱼、青菜", Dinner: "粥、炒菜", Snack: "水果"},
				{Day: "周六", Breakfast: "牛奶、面包、小菜", Lunch: "米饭、鸡肉、青菜", Dinner: "饺子", Snack: "水果"},
				{Day: "周日", Breakfast: "小米粥、鸡蛋、小菜", Lunch: "米饭、红烧肉、青菜", Dinner: "面条、炒菜", Snack: "水果"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.DietPlan{
			ID:          primitive.NewObjectID(),
			Name:        "低糖餐",
			Description: "适合糖尿病患者的低糖饮食",
			WeekMenu: []domain.WeekDayMenu{
				{Day: "周一", Breakfast: "燕麦粥、水煮蛋、小菜", Lunch: "糙米饭、清蒸鱼、青菜", Dinner: "小米粥、炒菜", Snack: "黄瓜"},
				{Day: "周二", Breakfast: "豆浆、全麦面包、小菜", Lunch: "糙米饭、白切鸡、青菜", Dinner: "面条、炒菜", Snack: "西红柿"},
				{Day: "周三", Breakfast: "小米粥、水煮蛋、小菜", Lunch: "糙米饭、清蒸鱼、青菜", Dinner: "粥、炒菜", Snack: "黄瓜"},
				{Day: "周四", Breakfast: "燕麦粥、水煮蛋、小菜", Lunch: "糙米饭、白切鸡、青菜", Dinner: "小米粥、炒菜", Snack: "西红柿"},
				{Day: "周五", Breakfast: "豆浆、全麦面包、小菜", Lunch: "糙米饭、清蒸鱼、青菜", Dinner: "面条、炒菜", Snack: "黄瓜"},
				{Day: "周六", Breakfast: "小米粥、水煮蛋、小菜", Lunch: "糙米饭、白切鸡、青菜", Dinner: "粥、炒菜", Snack: "西红柿"},
				{Day: "周日", Breakfast: "燕麦粥、水煮蛋、小菜", Lunch: "糙米饭、清蒸鱼、青菜", Dinner: "小米粥、炒菜", Snack: "黄瓜"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.DietPlan{
			ID:          primitive.NewObjectID(),
			Name:        "流质餐",
			Description: "适合咀嚼困难老人的流质饮食",
			WeekMenu: []domain.WeekDayMenu{
				{Day: "周一", Breakfast: "米汤、蛋花", Lunch: "肉末粥、菜泥", Dinner: "鱼汤、米糊", Snack: "果汁"},
				{Day: "周二", Breakfast: "豆浆、米糊", Lunch: "鸡汤、菜泥", Dinner: "肉末粥、米糊", Snack: "果汁"},
				{Day: "周三", Breakfast: "米汤、蛋花", Lunch: "鱼汤、菜泥", Dinner: "鸡汤、米糊", Snack: "果汁"},
				{Day: "周四", Breakfast: "豆浆、米糊", Lunch: "肉末粥、菜泥", Dinner: "鱼汤、米糊", Snack: "果汁"},
				{Day: "周五", Breakfast: "米汤、蛋花", Lunch: "鸡汤、菜泥", Dinner: "肉末粥、米糊", Snack: "果汁"},
				{Day: "周六", Breakfast: "豆浆、米糊", Lunch: "鱼汤、菜泥", Dinner: "鸡汤、米糊", Snack: "果汁"},
				{Day: "周日", Breakfast: "米汤、蛋花", Lunch: "肉末粥、菜泥", Dinner: "鱼汤、米糊", Snack: "果汁"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	dietPlanIDs := insertMany(db, "diet_plans", dietPlans)
	fmt.Printf("  ✓ 创建了 %d 个膳食计划\n", len(dietPlanIDs))

	// 4. 填充房间
	fmt.Println("\n4. 填充房间数据...")
	rooms := []interface{}{
		domain.Room{
			ID:          primitive.NewObjectID(),
			Number:      "A101",
			Floor:       1,
			Type:        "单人间",
			Capacity:    1,
			Status:      "可用",
			Description: "朝南，采光好，适合独居老人",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Room{
			ID:          primitive.NewObjectID(),
			Number:      "A102",
			Floor:       1,
			Type:        "单人间",
			Capacity:    1,
			Status:      "可用",
			Description: "朝南，采光好，适合独居老人",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Room{
			ID:          primitive.NewObjectID(),
			Number:      "A201",
			Floor:       2,
			Type:        "双人间",
			Capacity:    2,
			Status:      "可用",
			Description: "朝南，双人居住，设施齐全",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Room{
			ID:          primitive.NewObjectID(),
			Number:      "A202",
			Floor:       2,
			Type:        "双人间",
			Capacity:    2,
			Status:      "可用",
			Description: "朝南，双人居住，设施齐全",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Room{
			ID:          primitive.NewObjectID(),
			Number:      "B101",
			Floor:       1,
			Type:        "多人间",
			Capacity:    4,
			Status:      "可用",
			Description: "多人居住，经济实惠",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	roomIDs := insertMany(db, "rooms", rooms)
	fmt.Printf("  ✓ 创建了 %d 个房间\n", len(roomIDs))

	// 5. 填充床位
	fmt.Println("\n5. 填充床位数据...")
	beds := []interface{}{}
	bedIDs := []primitive.ObjectID{}

	// 为每个房间创建床位
	for i, roomID := range roomIDs {
		room := rooms[i].(domain.Room)
		for j := 1; j <= room.Capacity; j++ {
			bedID := primitive.NewObjectID()
			bedIDs = append(bedIDs, bedID)
			beds = append(beds, domain.Bed{
				ID:        bedID,
				RoomID:    roomID,
				Number:    fmt.Sprintf("%s-%d", room.Number, j),
				Status:    "空闲",
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	}
	insertMany(db, "beds", beds)
	fmt.Printf("  ✓ 创建了 %d 个床位\n", len(beds))

	// 6. 填充用户
	fmt.Println("\n6. 填充用户数据...")
	users := []interface{}{
		domain.User{
			ID:        primitive.NewObjectID(),
			Phone:     "13800002001",
			Password:  passwordHash,
			Name:      "张三",
			Age:       65,
			Gender:    "男",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.User{
			ID:        primitive.NewObjectID(),
			Phone:     "13800002002",
			Password:  passwordHash,
			Name:      "李四",
			Age:       70,
			Gender:    "女",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.User{
			ID:        primitive.NewObjectID(),
			Phone:     "13800002003",
			Password:  passwordHash,
			Name:      "王五",
			Age:       68,
			Gender:    "男",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.User{
			ID:        primitive.NewObjectID(),
			Phone:     "13800002004",
			Password:  passwordHash,
			Name:      "赵六",
			Age:       72,
			Gender:    "女",
			CreatedAt: now,
			UpdatedAt: now,
		},
		domain.User{
			ID:        primitive.NewObjectID(),
			Phone:     "13800002005",
			Password:  passwordHash,
			Name:      "孙七",
			Age:       75,
			Gender:    "男",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	userIDs := insertMany(db, "users", users)
	fmt.Printf("  ✓ 创建了 %d 个用户\n", len(userIDs))

	// 7. 填充客户
	fmt.Println("\n7. 填充客户数据...")
	customers := []interface{}{
		domain.Customer{
			ID:              primitive.NewObjectID(),
			UserID:          userIDs[0],
			Name:            "张三",
			Age:             65,
			Gender:          "男",
			Phone:           "13800002001",
			IDCard:          "110101195801010001",
			BedID:           bedIDs[0],
			DietPlanID:      dietPlanIDs[0],
			CareLevelID:     careLevelIDs[0],
			HealthManagerID: healthManagerIDs[0],
			HealthManager:   "张医生",
			Status:          "入住中",
			CheckInDate:     now.AddDate(0, -2, 0), // 2个月前入住
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		domain.Customer{
			ID:              primitive.NewObjectID(),
			UserID:          userIDs[1],
			Name:            "李四",
			Age:             70,
			Gender:          "女",
			Phone:           "13800002002",
			IDCard:          "110101195401010002",
			BedID:           bedIDs[1],
			DietPlanID:      dietPlanIDs[1],
			CareLevelID:     careLevelIDs[1],
			HealthManagerID: healthManagerIDs[1],
			HealthManager:   "李护士",
			Status:          "入住中",
			CheckInDate:     now.AddDate(0, -1, 0), // 1个月前入住
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		domain.Customer{
			ID:              primitive.NewObjectID(),
			UserID:          userIDs[2],
			Name:            "王五",
			Age:             68,
			Gender:          "男",
			Phone:           "13800002003",
			IDCard:          "110101195601010003",
			BedID:           bedIDs[2],
			DietPlanID:      dietPlanIDs[0],
			CareLevelID:     careLevelIDs[0],
			HealthManagerID: healthManagerIDs[0],
			HealthManager:   "张医生",
			Status:          "入住中",
			CheckInDate:     now.AddDate(0, -3, 0), // 3个月前入住
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		domain.Customer{
			ID:              primitive.NewObjectID(),
			UserID:          userIDs[3],
			Name:            "赵六",
			Age:             72,
			Gender:          "女",
			Phone:           "13800002004",
			IDCard:          "110101195201010004",
			BedID:           bedIDs[3],
			DietPlanID:      dietPlanIDs[2],
			CareLevelID:     careLevelIDs[2],
			HealthManagerID: healthManagerIDs[2],
			HealthManager:   "王营养师",
			Status:          "入住中",
			CheckInDate:     now.AddDate(0, -1, -10), // 1个月10天前入住
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		domain.Customer{
			ID:              primitive.NewObjectID(),
			UserID:          userIDs[4],
			Name:            "孙七",
			Age:             75,
			Gender:          "男",
			Phone:           "13800002005",
			IDCard:          "110101194901010005",
			DietPlanID:      dietPlanIDs[1],
			CareLevelID:     careLevelIDs[1],
			HealthManagerID: healthManagerIDs[0],
			HealthManager:   "张医生",
			Status:          "未入住",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}
	customerIDs := insertMany(db, "customers", customers)
	fmt.Printf("  ✓ 创建了 %d 个客户\n", len(customerIDs))

	// 更新床位状态（已分配的床位）
	updateBedStatus(db, bedIDs[0], customerIDs[0], "占用")
	updateBedStatus(db, bedIDs[1], customerIDs[1], "占用")
	updateBedStatus(db, bedIDs[2], customerIDs[2], "占用")
	updateBedStatus(db, bedIDs[3], customerIDs[3], "占用")

	// 8. 填充服务项目
	fmt.Println("\n8. 填充服务项目数据...")
	services := []interface{}{
		domain.Service{
			ID:          primitive.NewObjectID(),
			Name:        "健康体检",
			Description: "全面健康体检服务",
			Category:    "医疗",
			Price:       200.0,
			Unit:        "次",
			Status:      "启用",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Service{
			ID:          primitive.NewObjectID(),
			Name:        "康复训练",
			Description: "专业康复训练指导",
			Category:    "医疗",
			Price:       150.0,
			Unit:        "次",
			Status:      "启用",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Service{
			ID:          primitive.NewObjectID(),
			Name:        "娱乐活动",
			Description: "组织各类娱乐活动",
			Category:    "娱乐",
			Price:       50.0,
			Unit:        "次",
			Status:      "启用",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Service{
			ID:          primitive.NewObjectID(),
			Name:        "心理辅导",
			Description: "专业心理辅导服务",
			Category:    "医疗",
			Price:       100.0,
			Unit:        "次",
			Status:      "启用",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.Service{
			ID:          primitive.NewObjectID(),
			Name:        "生活照料",
			Description: "日常生活照料服务",
			Category:    "生活",
			Price:       80.0,
			Unit:        "次",
			Status:      "启用",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	serviceIDs := insertMany(db, "services", services)
	fmt.Printf("  ✓ 创建了 %d 个服务项目\n", len(serviceIDs))

	// 9. 填充客户购买的服务
	fmt.Println("\n9. 填充客户购买的服务数据...")
	customerServices := []interface{}{
		domain.CustomerService{
			ID:          primitive.NewObjectID(),
			CustomerID:  customerIDs[0],
			ServiceID:   serviceIDs[0],
			ServiceName: "健康体检",
			StartDate:   now.AddDate(0, -1, 0),
			Status:      "进行中",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.CustomerService{
			ID:          primitive.NewObjectID(),
			CustomerID:  customerIDs[0],
			ServiceID:   serviceIDs[1],
			ServiceName: "康复训练",
			StartDate:   now.AddDate(0, -1, 0),
			Status:      "进行中",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.CustomerService{
			ID:          primitive.NewObjectID(),
			CustomerID:  customerIDs[1],
			ServiceID:   serviceIDs[2],
			ServiceName: "娱乐活动",
			StartDate:   now.AddDate(0, 0, -10),
			Status:      "进行中",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		domain.CustomerService{
			ID:          primitive.NewObjectID(),
			CustomerID:  customerIDs[2],
			ServiceID:   serviceIDs[3],
			ServiceName: "心理辅导",
			StartDate:   now.AddDate(0, -2, 0),
			EndDate:     now.AddDate(0, -1, 0),
			Status:      "已结束",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	insertMany(db, "customer_services", customerServices)
	fmt.Printf("  ✓ 创建了 %d 条客户服务记录\n", len(customerServices))

	// 10. 填充登记记录
	fmt.Println("\n10. 填充登记记录数据...")
	records := []interface{}{
		domain.Record{
			ID:         primitive.NewObjectID(),
			CustomerID: customerIDs[0],
			Type:       "入住",
			StartTime:  now.AddDate(0, -2, 0),
			Note:       "客户入住登记",
			CreatedBy:  "管理员",
			CreatedAt:  now,
		},
		domain.Record{
			ID:         primitive.NewObjectID(),
			CustomerID: customerIDs[1],
			Type:       "入住",
			StartTime:  now.AddDate(0, -1, 0),
			Note:       "客户入住登记",
			CreatedBy:  "管理员",
			CreatedAt:  now,
		},
		domain.Record{
			ID:         primitive.NewObjectID(),
			CustomerID: customerIDs[2],
			Type:       "入住",
			StartTime:  now.AddDate(0, -3, 0),
			Note:       "客户入住登记",
			CreatedBy:  "管理员",
			CreatedAt:  now,
		},
		domain.Record{
			ID:         primitive.NewObjectID(),
			CustomerID: customerIDs[0],
			Type:       "外出",
			StartTime:  now.AddDate(0, 0, -5),
			EndTime:    now.AddDate(0, 0, -4),
			Note:       "外出就医",
			CreatedBy:  "管理员",
			CreatedAt:  now,
		},
		domain.Record{
			ID:         primitive.NewObjectID(),
			CustomerID: customerIDs[1],
			Type:       "外出",
			StartTime:  now.AddDate(0, 0, -2),
			Note:       "外出探亲",
			CreatedBy:  "管理员",
			CreatedAt:  now,
		},
	}
	insertMany(db, "records", records)
	fmt.Printf("  ✓ 创建了 %d 条登记记录\n", len(records))

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✓ 数据库填充完成！")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\n默认密码: %s\n", defaultPass)
	fmt.Println("\n测试账号:")
	for i, user := range users {
		u := user.(domain.User)
		fmt.Printf("  手机号: %s, 密码: %s, 姓名: %s\n", u.Phone, defaultPass, u.Name)
		if i >= 2 {
			break
		}
	}
}

// insertMany 批量插入数据
func insertMany(db *mongo.Database, collectionName string, documents []interface{}) []primitive.ObjectID {
	if len(documents) == 0 {
		return []primitive.ObjectID{}
	}

	ctx := context.Background()
	coll := db.Collection(collectionName)

	// 清空集合（可选，如果需要重新填充）
	// coll.DeleteMany(ctx, bson.M{})

	result, err := coll.InsertMany(ctx, documents)
	if err != nil {
		log.Fatalf("插入 %s 失败: %v", collectionName, err)
	}

	ids := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		ids[i] = id.(primitive.ObjectID)
	}

	return ids
}

// updateBedStatus 更新床位状态
func updateBedStatus(db *mongo.Database, bedID, customerID primitive.ObjectID, status string) {
	ctx := context.Background()
	coll := db.Collection("beds")

	_, err := coll.UpdateOne(
		ctx,
		map[string]interface{}{"_id": bedID},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"customer_id": customerID,
				"status":      status,
				"updated_at":  time.Now(),
			},
		},
	)
	if err != nil {
		log.Printf("更新床位状态失败: %v", err)
	}
}

