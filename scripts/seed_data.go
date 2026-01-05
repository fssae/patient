package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB连接配置
const (
	MongoURI     = "mongodb://admin:zjh770910@82.156.64.69:27017"
	DatabaseName = "kongdong"
)

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoURI))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping MongoDB失败:", err)
	}

	db := client.Database(DatabaseName)
	fmt.Println("✅ 成功连接到MongoDB数据库:", DatabaseName)

	// 执行数据填充
	if err := seedData(ctx, db); err != nil {
		log.Fatal("数据填充失败:", err)
	}

	fmt.Println("\n🎉 所有数据填充完成!")
}

func seedData(ctx context.Context, db *mongo.Database) error {
	// 1. 填充用户数据
	userIDs, err := seedUsers(ctx, db)
	if err != nil {
		return fmt.Errorf("填充用户数据失败: %w", err)
	}

	// 2. 填充房间数据
	roomIDs, err := seedRooms(ctx, db)
	if err != nil {
		return fmt.Errorf("填充房间数据失败: %w", err)
	}

	// 3. 填充床位数据
	bedIDs, err := seedBeds(ctx, db, roomIDs)
	if err != nil {
		return fmt.Errorf("填充床位数据失败: %w", err)
	}

	// 4. 填充膳食计划
	dietPlanIDs, err := seedDietPlans(ctx, db)
	if err != nil {
		return fmt.Errorf("填充膳食计划失败: %w", err)
	}

	// 5. 填充护理级别
	careLevelIDs, err := seedCareLevels(ctx, db)
	if err != nil {
		return fmt.Errorf("填充护理级别失败: %w", err)
	}

	// 6. 填充健康管家
	healthManagerIDs, err := seedHealthManagers(ctx, db)
	if err != nil {
		return fmt.Errorf("填充健康管家失败: %w", err)
	}

	// 7. 填充客户(入住老人)
	customerIDs, err := seedCustomers(ctx, db, userIDs, bedIDs, dietPlanIDs, careLevelIDs, healthManagerIDs)
	if err != nil {
		return fmt.Errorf("填充客户数据失败: %w", err)
	}

	// 8. 更新床位的客户关联
	if err := updateBedsWithCustomers(ctx, db, bedIDs, customerIDs); err != nil {
		return fmt.Errorf("更新床位关联失败: %w", err)
	}

	// 9. 填充服务项目
	serviceIDs, err := seedServices(ctx, db)
	if err != nil {
		return fmt.Errorf("填充服务项目失败: %w", err)
	}

	// 10. 填充客户服务
	if err := seedCustomerServices(ctx, db, customerIDs, serviceIDs); err != nil {
		return fmt.Errorf("填充客户服务失败: %w", err)
	}

	// 11. 填充登记记录
	if err := seedRecords(ctx, db, customerIDs); err != nil {
		return fmt.Errorf("填充登记记录失败: %w", err)
	}

	// 12. 填充护理记录
	if err := seedCareRecords(ctx, db, customerIDs); err != nil {
		return fmt.Errorf("填充护理记录失败: %w", err)
	}

	// 13. 填充患者数据
	patientIDs, err := seedPatients(ctx, db)
	if err != nil {
		return fmt.Errorf("填充患者数据失败: %w", err)
	}

	// 14. 填充分析日志(告警数据)
	if err := seedAnalysisLogs(ctx, db, patientIDs, bedIDs); err != nil {
		return fmt.Errorf("填充分析日志失败: %w", err)
	}

	return nil
}

// 1. 填充用户数据
func seedUsers(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("users")

	now := time.Now()
	users := []interface{}{
		map[string]interface{}{
			"phone":      "13800138001",
			"password":   "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH", // 加密后的密码: password123
			"name":       "张伟",
			"age":        45,
			"gender":     "男",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"phone":      "13800138002",
			"password":   "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"name":       "李娜",
			"age":        42,
			"gender":     "女",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"phone":      "13800138003",
			"password":   "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"name":       "王强",
			"age":        50,
			"gender":     "男",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"phone":      "13800138004",
			"password":   "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"name":       "刘芳",
			"age":        38,
			"gender":     "女",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"phone":      "13800138005",
			"password":   "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"name":       "陈明",
			"age":        48,
			"gender":     "男",
			"created_at": now,
			"updated_at": now,
		},
	}

	result, err := collection.InsertMany(ctx, users)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条用户数据\n", len(ids))
	return ids, nil
}

// 2. 填充房间数据
func seedRooms(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("rooms")

	now := time.Now()
	rooms := []interface{}{
		map[string]interface{}{
			"number":      "A101",
			"floor":       1,
			"type":        "双人间",
			"capacity":    2,
			"status":      "可用",
			"description": "一楼阳光房，采光好",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"number":      "A102",
			"floor":       1,
			"type":        "单人间",
			"capacity":    1,
			"status":      "可用",
			"description": "一楼独立房间，安静舒适",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"number":      "A201",
			"floor":       2,
			"type":        "多人间",
			"capacity":    4,
			"status":      "可用",
			"description": "二楼大房间，适合多人居住",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"number":      "A202",
			"floor":       2,
			"type":        "双人间",
			"capacity":    2,
			"status":      "可用",
			"description": "二楼标准双人间",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"number":      "B101",
			"floor":       1,
			"type":        "单人间",
			"capacity":    1,
			"status":      "维护中",
			"description": "B栋一楼单人间，正在维护",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"number":      "B201",
			"floor":       2,
			"type":        "双人间",
			"capacity":    2,
			"status":      "可用",
			"description": "B栋二楼双人间",
			"created_at":  now,
			"updated_at":  now,
		},
	}

	result, err := collection.InsertMany(ctx, rooms)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条房间数据\n", len(ids))
	return ids, nil
}

// 3. 填充床位数据
func seedBeds(ctx context.Context, db *mongo.Database, roomIDs []primitive.ObjectID) ([]primitive.ObjectID, error) {
	collection := db.Collection("beds")

	now := time.Now()
	var beds []interface{}

	// A101 双人间 - 2张床
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[0],
			"number":     "A101-1",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[0],
			"number":     "A101-2",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
	)

	// A102 单人间 - 1张床
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[1],
			"number":     "A102-1",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
	)

	// A201 多人间 - 4张床
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[2],
			"number":     "A201-1",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[2],
			"number":     "A201-2",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[2],
			"number":     "A201-3",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[2],
			"number":     "A201-4",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
	)

	// A202 双人间 - 2张床
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[3],
			"number":     "A202-1",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[3],
			"number":     "A202-2",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
	)

	// B101 单人间 - 1张床 (维护中)
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[4],
			"number":     "B101-1",
			"status":     "维护中",
			"created_at": now,
			"updated_at": now,
		},
	)

	// B201 双人间 - 2张床
	beds = append(beds,
		map[string]interface{}{
			"room_id":    roomIDs[5],
			"number":     "B201-1",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"room_id":    roomIDs[5],
			"number":     "B201-2",
			"status":     "空闲",
			"created_at": now,
			"updated_at": now,
		},
	)

	result, err := collection.InsertMany(ctx, beds)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条床位数据\n", len(ids))
	return ids, nil
}

// 4. 填充膳食计划
func seedDietPlans(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("diet_plans")

	now := time.Now()
	dietPlans := []interface{}{
		map[string]interface{}{
			"name":        "普通餐",
			"description": "适合身体健康的老人",
			"week_menu": []map[string]interface{}{
				{"day": "周一", "breakfast": "小米粥、鸡蛋、馒头", "lunch": "米饭、红烧肉、青菜", "dinner": "面条、炒鸡蛋、凉拌黄瓜"},
				{"day": "周二", "breakfast": "豆浆、油条、咸菜", "lunch": "米饭、清蒸鱼、西红柿炒蛋", "dinner": "米饭、炖排骨、炒青菜"},
				{"day": "周三", "breakfast": "牛奶、面包、水果", "lunch": "米饭、宫保鸡丁、炒豆角", "dinner": "馄饨、凉拌海带丝"},
				{"day": "周四", "breakfast": "八宝粥、煮鸡蛋、小菜", "lunch": "米饭、红烧鱼、炒白菜", "dinner": "米饭、炒肉片、紫菜汤"},
				{"day": "周五", "breakfast": "豆腐脑、油饼、咸菜", "lunch": "米饭、糖醋里脊、炒菠菜", "dinner": "饺子、凉拌木耳"},
				{"day": "周六", "breakfast": "南瓜粥、鸡蛋、馒头", "lunch": "米饭、炖鸡、炒土豆丝", "dinner": "炒面、凉拌黄瓜"},
				{"day": "周日", "breakfast": "牛奶、蛋糕、水果", "lunch": "米饭、红烧肉、炒青菜", "dinner": "米饭、清蒸鱼、西红柿鸡蛋汤"},
			},
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":        "低糖餐",
			"description": "适合糖尿病患者",
			"week_menu": []map[string]interface{}{
				{"day": "周一", "breakfast": "燕麦粥、煮鸡蛋、全麦面包", "lunch": "糙米饭、清蒸鱼、炒青菜", "dinner": "荞麦面、凉拌豆腐"},
				{"day": "周二", "breakfast": "无糖豆浆、全麦馒头、小菜", "lunch": "糙米饭、炖鸡胸肉、炒西兰花", "dinner": "杂粮粥、炒鸡蛋、凉拌黄瓜"},
				{"day": "周三", "breakfast": "无糖酸奶、全麦面包、水果", "lunch": "糙米饭、清蒸鲈鱼、炒芹菜", "dinner": "玉米面粥、炒豆腐、凉拌海带"},
				{"day": "周四", "breakfast": "燕麦粥、煮鸡蛋、小菜", "lunch": "糙米饭、炖牛肉、炒菠菜", "dinner": "荞麦面、炒鸡蛋、紫菜汤"},
				{"day": "周五", "breakfast": "无糖豆浆、全麦馒头、咸菜", "lunch": "糙米饭、清蒸鱼、炒豆角", "dinner": "杂粮粥、炒豆腐、凉拌木耳"},
				{"day": "周六", "breakfast": "燕麦粥、煮鸡蛋、全麦面包", "lunch": "糙米饭、炖鸡肉、炒白菜", "dinner": "玉米面粥、炒鸡蛋、凉拌黄瓜"},
				{"day": "周日", "breakfast": "无糖酸奶、全麦面包、水果", "lunch": "糙米饭、清蒸鱼、炒青菜", "dinner": "荞麦面、炒豆腐、紫菜汤"},
			},
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":        "流质餐",
			"description": "适合咀嚼困难的老人",
			"week_menu": []map[string]interface{}{
				{"day": "周一", "breakfast": "小米粥、蒸鸡蛋羹", "lunch": "烂面条、肉末粥", "dinner": "南瓜粥、豆腐羹"},
				{"day": "周二", "breakfast": "燕麦粥、蒸鸡蛋羹", "lunch": "鱼肉粥、蔬菜泥", "dinner": "小米粥、肉末羹"},
				{"day": "周三", "breakfast": "牛奶、蒸鸡蛋羹", "lunch": "烂面条、鸡肉粥", "dinner": "南瓜粥、豆腐羹"},
				{"day": "周四", "breakfast": "小米粥、蒸鸡蛋羹", "lunch": "肉末粥、蔬菜泥", "dinner": "燕麦粥、鱼肉羹"},
				{"day": "周五", "breakfast": "牛奶、蒸鸡蛋羹", "lunch": "烂面条、鸡肉粥", "dinner": "小米粥、豆腐羹"},
				{"day": "周六", "breakfast": "燕麦粥、蒸鸡蛋羹", "lunch": "鱼肉粥、蔬菜泥", "dinner": "南瓜粥、肉末羹"},
				{"day": "周日", "breakfast": "牛奶、蒸鸡蛋羹", "lunch": "烂面条、鸡肉粥", "dinner": "小米粥、豆腐羹"},
			},
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":        "低盐餐",
			"description": "适合高血压患者",
			"week_menu": []map[string]interface{}{
				{"day": "周一", "breakfast": "小米粥、煮鸡蛋、无盐馒头", "lunch": "米饭、清蒸鱼、炒青菜", "dinner": "面条、炒鸡蛋、凉拌黄瓜"},
				{"day": "周二", "breakfast": "豆浆、无盐面包、水果", "lunch": "米饭、炖鸡肉、炒西兰花", "dinner": "米饭、清蒸鱼、紫菜汤"},
				{"day": "周三", "breakfast": "牛奶、无盐馒头、小菜", "lunch": "米饭、炖牛肉、炒菠菜", "dinner": "面条、炒豆腐、凉拌海带"},
				{"day": "周四", "breakfast": "八宝粥、煮鸡蛋、无盐面包", "lunch": "米饭、清蒸鱼、炒白菜", "dinner": "米饭、炖鸡肉、紫菜汤"},
				{"day": "周五", "breakfast": "豆浆、无盐馒头、水果", "lunch": "米饭、炖牛肉、炒豆角", "dinner": "面条、炒鸡蛋、凉拌黄瓜"},
				{"day": "周六", "breakfast": "小米粥、煮鸡蛋、无盐面包", "lunch": "米饭、清蒸鱼、炒青菜", "dinner": "米饭、炖鸡肉、紫菜汤"},
				{"day": "周日", "breakfast": "牛奶、无盐馒头、水果", "lunch": "米饭、炖牛肉、炒西兰花", "dinner": "面条、炒豆腐、凉拌海带"},
			},
			"created_at": now,
			"updated_at": now,
		},
	}

	result, err := collection.InsertMany(ctx, dietPlans)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条膳食计划数据\n", len(ids))
	return ids, nil
}

// 5. 填充护理级别
func seedCareLevels(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("care_levels")

	now := time.Now()
	careLevels := []interface{}{
		map[string]interface{}{
			"name":        "一级护理",
			"level":       1,
			"description": "适合生活完全不能自理的老人",
			"content":     "24小时专人护理，协助进食、洗漱、如厕、翻身等日常生活",
			"price":       3000.00,
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "二级护理",
			"level":       2,
			"description": "适合生活部分不能自理的老人",
			"content":     "定时巡视，协助部分日常生活，如洗澡、更衣等",
			"price":       2000.00,
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "三级护理",
			"level":       3,
			"description": "适合生活基本能自理的老人",
			"content":     "定期巡视，提供必要的生活协助和健康监测",
			"price":       1000.00,
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "特级护理",
			"level":       0,
			"description": "适合病情危重需要重点监护的老人",
			"content":     "24小时专业医护人员监护，随时观察病情变化",
			"price":       5000.00,
			"created_at":  now,
			"updated_at":  now,
		},
	}

	result, err := collection.InsertMany(ctx, careLevels)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条护理级别数据\n", len(ids))
	return ids, nil
}

// 6. 填充健康管家
func seedHealthManagers(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("health_managers")

	now := time.Now()
	managers := []interface{}{
		map[string]interface{}{
			"name":       "王医生",
			"phone":      "13900139001",
			"email":      "wangdoctor@example.com",
			"specialty":  "内科",
			"status":     "在职",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":       "李护士",
			"phone":      "13900139002",
			"email":      "linurse@example.com",
			"specialty":  "老年护理",
			"status":     "在职",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":       "张医生",
			"phone":      "13900139003",
			"email":      "zhangdoctor@example.com",
			"specialty":  "心血管科",
			"status":     "在职",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":       "刘护士",
			"phone":      "13900139004",
			"email":      "liunurse@example.com",
			"specialty":  "康复护理",
			"status":     "在职",
			"created_at": now,
			"updated_at": now,
		},
		map[string]interface{}{
			"name":       "陈医生",
			"phone":      "13900139005",
			"email":      "chendoctor@example.com",
			"specialty":  "神经内科",
			"status":     "在职",
			"created_at": now,
			"updated_at": now,
		},
	}

	result, err := collection.InsertMany(ctx, managers)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条健康管家数据\n", len(ids))
	return ids, nil
}

// 7. 填充客户(入住老人)
func seedCustomers(ctx context.Context, db *mongo.Database, userIDs, bedIDs, dietPlanIDs, careLevelIDs, healthManagerIDs []primitive.ObjectID) ([]primitive.ObjectID, error) {
	collection := db.Collection("customers")

	now := time.Now()
	checkInDate := now.AddDate(0, -3, 0) // 3个月前入住

	customers := []interface{}{
		map[string]interface{}{
			"user_id":           userIDs[0],
			"name":              "张老伯",
			"age":               78,
			"gender":            "男",
			"phone":             "13700137001",
			"id_card":           "310101194601011234",
			"bed_id":            bedIDs[0],
			"diet_plan_id":      dietPlanIDs[1],  // 低糖餐
			"care_level_id":     careLevelIDs[1], // 二级护理
			"health_manager_id": healthManagerIDs[0],
			"health_manager":    "王医生",
			"status":            "入住中",
			"check_in_date":     checkInDate,
			"health_level":      "一般",
			"medical_history":   "高血压、糖尿病",
			"medication":        "降压药、降糖药",
			"allergy_history":   "青霉素过敏",
			"contact_name":      "张伟",
			"relationship":      "儿子",
			"contact_phone":     "13800138001",
			"contact_address":   "上海市浦东新区XX路XX号",
			"remarks":           "需要定期测血糖血压",
			"created_at":        checkInDate,
			"updated_at":        now,
		},
		map[string]interface{}{
			"user_id":           userIDs[1],
			"name":              "李奶奶",
			"age":               82,
			"gender":            "女",
			"phone":             "13700137002",
			"id_card":           "310101194201021234",
			"bed_id":            bedIDs[1],
			"diet_plan_id":      dietPlanIDs[2],  // 流质餐
			"care_level_id":     careLevelIDs[0], // 一级护理
			"health_manager_id": healthManagerIDs[1],
			"health_manager":    "李护士",
			"status":            "入住中",
			"check_in_date":     checkInDate.AddDate(0, -1, 0),
			"health_level":      "较差",
			"medical_history":   "脑梗、高血压",
			"medication":        "阿司匹林、降压药",
			"allergy_history":   "无",
			"contact_name":      "李娜",
			"relationship":      "女儿",
			"contact_phone":     "13800138002",
			"contact_address":   "上海市徐汇区XX路XX号",
			"remarks":           "行动不便，需要轮椅",
			"created_at":        checkInDate.AddDate(0, -1, 0),
			"updated_at":        now,
		},
		map[string]interface{}{
			"user_id":           userIDs[2],
			"name":              "王老先生",
			"age":               75,
			"gender":            "男",
			"phone":             "13700137003",
			"id_card":           "310101194901031234",
			"bed_id":            bedIDs[2],
			"diet_plan_id":      dietPlanIDs[0],  // 普通餐
			"care_level_id":     careLevelIDs[2], // 三级护理
			"health_manager_id": healthManagerIDs[2],
			"health_manager":    "张医生",
			"status":            "入住中",
			"check_in_date":     checkInDate.AddDate(0, 0, -15),
			"health_level":      "健康",
			"medical_history":   "无重大疾病",
			"medication":        "无",
			"allergy_history":   "无",
			"contact_name":      "王强",
			"relationship":      "儿子",
			"contact_phone":     "13800138003",
			"contact_address":   "上海市静安区XX路XX号",
			"remarks":           "身体健康，生活能自理",
			"created_at":        checkInDate.AddDate(0, 0, -15),
			"updated_at":        now,
		},
		map[string]interface{}{
			"user_id":           userIDs[3],
			"name":              "刘奶奶",
			"age":               80,
			"gender":            "女",
			"phone":             "13700137004",
			"id_card":           "310101194401041234",
			"bed_id":            bedIDs[3],
			"diet_plan_id":      dietPlanIDs[3],  // 低盐餐
			"care_level_id":     careLevelIDs[1], // 二级护理
			"health_manager_id": healthManagerIDs[3],
			"health_manager":    "刘护士",
			"status":            "入住中",
			"check_in_date":     checkInDate.AddDate(0, 0, -20),
			"health_level":      "较好",
			"medical_history":   "高血压",
			"medication":        "降压药",
			"allergy_history":   "海鲜过敏",
			"contact_name":      "刘芳",
			"relationship":      "女儿",
			"contact_phone":     "13800138004",
			"contact_address":   "上海市黄浦区XX路XX号",
			"remarks":           "需要低盐饮食",
			"created_at":        checkInDate.AddDate(0, 0, -20),
			"updated_at":        now,
		},
		map[string]interface{}{
			"user_id":           userIDs[4],
			"name":              "陈老伯",
			"age":               76,
			"gender":            "男",
			"phone":             "13700137005",
			"id_card":           "310101194801051234",
			"bed_id":            bedIDs[4],
			"diet_plan_id":      dietPlanIDs[1],  // 低糖餐
			"care_level_id":     careLevelIDs[2], // 三级护理
			"health_manager_id": healthManagerIDs[4],
			"health_manager":    "陈医生",
			"status":            "入住中",
			"check_in_date":     checkInDate.AddDate(0, 0, -10),
			"health_level":      "一般",
			"medical_history":   "糖尿病、冠心病",
			"medication":        "降糖药、心脏病药物",
			"allergy_history":   "无",
			"contact_name":      "陈明",
			"relationship":      "儿子",
			"contact_phone":     "13800138005",
			"contact_address":   "上海市长宁区XX路XX号",
			"remarks":           "需要定期测血糖",
			"created_at":        checkInDate.AddDate(0, 0, -10),
			"updated_at":        now,
		},
	}

	result, err := collection.InsertMany(ctx, customers)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条客户数据\n", len(ids))
	return ids, nil
}

// 8. 更新床位的客户关联
func updateBedsWithCustomers(ctx context.Context, db *mongo.Database, bedIDs, customerIDs []primitive.ObjectID) error {
	collection := db.Collection("beds")

	now := time.Now()

	// 更新前5个床位为占用状态
	for i := 0; i < 5 && i < len(bedIDs) && i < len(customerIDs); i++ {
		_, err := collection.UpdateOne(
			ctx,
			map[string]interface{}{"_id": bedIDs[i]},
			map[string]interface{}{
				"$set": map[string]interface{}{
					"customer_id": customerIDs[i],
					"status":      "占用",
					"updated_at":  now,
				},
			},
		)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ 成功更新 %d 个床位的客户关联\n", 5)
	return nil
}

// 9. 填充服务项目
func seedServices(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("services")

	now := time.Now()
	services := []interface{}{
		map[string]interface{}{
			"name":        "康复理疗",
			"description": "专业康复师提供物理治疗",
			"category":    "医疗",
			"price":       200.00,
			"unit":        "次",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "定期体检",
			"description": "每月一次全面体检",
			"category":    "医疗",
			"price":       500.00,
			"unit":        "月",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "洗衣服务",
			"description": "专人负责衣物清洗",
			"category":    "生活",
			"price":       300.00,
			"unit":        "月",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "理发服务",
			"description": "专业理发师上门服务",
			"category":    "生活",
			"price":       50.00,
			"unit":        "次",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "文娱活动",
			"description": "组织各类文娱活动",
			"category":    "娱乐",
			"price":       200.00,
			"unit":        "月",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "心理咨询",
			"description": "专业心理咨询师服务",
			"category":    "医疗",
			"price":       300.00,
			"unit":        "次",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
		map[string]interface{}{
			"name":        "陪同就医",
			"description": "护理人员陪同外出就医",
			"category":    "生活",
			"price":       150.00,
			"unit":        "次",
			"status":      "启用",
			"created_at":  now,
			"updated_at":  now,
		},
	}

	result, err := collection.InsertMany(ctx, services)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条服务项目数据\n", len(ids))
	return ids, nil
}

// 10. 填充客户服务
func seedCustomerServices(ctx context.Context, db *mongo.Database, customerIDs, serviceIDs []primitive.ObjectID) error {
	collection := db.Collection("customer_services")

	now := time.Now()
	startDate := now.AddDate(0, -2, 0) // 2个月前开始
	endDate := now.AddDate(0, 1, 0)    // 1个月后结束

	var customerServices []interface{}

	// 为每个客户分配一些服务
	for i, customerID := range customerIDs {
		// 每个客户购买2-3个服务
		numServices := 2 + (i % 2)
		for j := 0; j < numServices && j < len(serviceIDs); j++ {
			serviceIdx := (i + j) % len(serviceIDs)
			customerServices = append(customerServices, map[string]interface{}{
				"customer_id": customerID,
				"service_id":  serviceIDs[serviceIdx],
				"start_date":  startDate,
				"end_date":    endDate,
				"status":      "进行中",
				"created_at":  startDate,
				"updated_at":  now,
			})
		}
	}

	if len(customerServices) > 0 {
		_, err := collection.InsertMany(ctx, customerServices)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ 成功插入 %d 条客户服务数据\n", len(customerServices))
	return nil
}

// 11. 填充登记记录
func seedRecords(ctx context.Context, db *mongo.Database, customerIDs []primitive.ObjectID) error {
	collection := db.Collection("records")

	now := time.Now()
	var records []interface{}

	// 为每个客户创建入住记录
	for i, customerID := range customerIDs {
		checkInTime := now.AddDate(0, -3, -i*5) // 不同时间入住

		// 入住记录
		records = append(records, map[string]interface{}{
			"customer_id": customerID,
			"type":        "入住",
			"start_time":  checkInTime,
			"note":        "正常入住",
			"created_by":  "管理员",
			"created_at":  checkInTime,
		})

		// 部分客户有外出记录
		if i%2 == 0 {
			outTime := checkInTime.AddDate(0, 1, 0)
			returnTime := outTime.Add(4 * time.Hour)

			records = append(records, map[string]interface{}{
				"customer_id": customerID,
				"type":        "外出",
				"start_time":  outTime,
				"end_time":    returnTime,
				"note":        "家属陪同外出",
				"created_by":  "护士",
				"created_at":  outTime,
			})
		}
	}

	if len(records) > 0 {
		_, err := collection.InsertMany(ctx, records)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ 成功插入 %d 条登记记录数据\n", len(records))
	return nil
}

// 12. 填充护理记录
func seedCareRecords(ctx context.Context, db *mongo.Database, customerIDs []primitive.ObjectID) error {
	collection := db.Collection("care_records")

	now := time.Now()
	var careRecords []interface{}

	customerNames := []string{"张老伯", "李奶奶", "王老先生", "刘奶奶", "陈老伯"}

	// 为每个客户创建护理记录
	for i, customerID := range customerIDs {
		var recordItems []map[string]interface{}

		// 创建最近7天的护理记录
		for day := 0; day < 7; day++ {
			careTime := now.AddDate(0, 0, -day)

			// 每天3条护理记录
			recordItems = append(recordItems,
				map[string]interface{}{
					"care_item":      "测量血压",
					"care_time":      careTime.Add(-8 * time.Hour),
					"care_personnel": "李护士",
					"care_result":    "血压正常 120/80",
					"created_at":     careTime.Add(-8 * time.Hour),
				},
				map[string]interface{}{
					"care_item":      "协助用餐",
					"care_time":      careTime.Add(-6 * time.Hour),
					"care_personnel": "王护工",
					"care_result":    "用餐正常",
					"created_at":     careTime.Add(-6 * time.Hour),
				},
				map[string]interface{}{
					"care_item":      "晚间巡视",
					"care_time":      careTime.Add(-2 * time.Hour),
					"care_personnel": "刘护士",
					"care_result":    "睡眠良好",
					"created_at":     careTime.Add(-2 * time.Hour),
				},
			)
		}

		careRecords = append(careRecords, map[string]interface{}{
			"customer_id":   customerID,
			"customer_name": customerNames[i],
			"records":       recordItems,
			"updated_at":    now,
		})
	}

	if len(careRecords) > 0 {
		_, err := collection.InsertMany(ctx, careRecords)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ 成功插入 %d 条护理记录数据\n", len(careRecords))
	return nil
}

// 13. 填充患者数据
func seedPatients(ctx context.Context, db *mongo.Database) ([]primitive.ObjectID, error) {
	collection := db.Collection("patients")

	now := time.Now()
	patients := []interface{}{
		map[string]interface{}{
			"patientId": "P001",
			"name":      "张老伯",
			"password":  "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"email":     "zhanglaobai@example.com",
			"phone":     "13700137001",
			"avatar":    "https://api.dicebear.com/7.x/avataaars/svg?seed=Zhang",
			"createdAt": now,
			"updatedAt": now,
		},
		map[string]interface{}{
			"patientId": "P002",
			"name":      "李奶奶",
			"password":  "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"email":     "linaainai@example.com",
			"phone":     "13700137002",
			"avatar":    "https://api.dicebear.com/7.x/avataaars/svg?seed=Li",
			"createdAt": now,
			"updatedAt": now,
		},
		map[string]interface{}{
			"patientId": "P003",
			"name":      "王老先生",
			"password":  "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"email":     "wanglaoxiansheng@example.com",
			"phone":     "13700137003",
			"avatar":    "https://api.dicebear.com/7.x/avataaars/svg?seed=Wang",
			"createdAt": now,
			"updatedAt": now,
		},
		map[string]interface{}{
			"patientId": "P004",
			"name":      "刘奶奶",
			"password":  "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"email":     "liunainai@example.com",
			"phone":     "13700137004",
			"avatar":    "https://api.dicebear.com/7.x/avataaars/svg?seed=Liu",
			"createdAt": now,
			"updatedAt": now,
		},
		map[string]interface{}{
			"patientId": "P005",
			"name":      "陈老伯",
			"password":  "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH",
			"email":     "chenlaobai@example.com",
			"phone":     "13700137005",
			"avatar":    "https://api.dicebear.com/7.x/avataaars/svg?seed=Chen",
			"createdAt": now,
			"updatedAt": now,
		},
	}

	result, err := collection.InsertMany(ctx, patients)
	if err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Printf("✅ 成功插入 %d 条患者数据\n", len(ids))
	return ids, nil
}

// 14. 填充分析日志(告警数据)
func seedAnalysisLogs(ctx context.Context, db *mongo.Database, patientIDs, bedIDs []primitive.ObjectID) error {
	collection := db.Collection("analysis")

	now := time.Now()
	var analysisLogs []interface{}

	patientIDStrs := []string{"P001", "P002", "P003", "P004", "P005"}
	bedIDStrs := []string{"A101-1", "A101-2", "A102-1", "A201-1", "A201-2"}

	// 创建一些跌倒检测事件
	for i := 0; i < 5; i++ {
		eventTime := now.AddDate(0, 0, -i).Unix()

		analysisLogs = append(analysisLogs, map[string]interface{}{
			"type":        "event",
			"timestamp":   eventTime,
			"event_type":  "fall",
			"level":       "critical",
			"message":     fmt.Sprintf("检测到%s跌倒", patientIDStrs[i%5]),
			"video_url":   fmt.Sprintf("http://82.156.64.69:9000/patient/evidence/fall_%s.mp4", time.Unix(eventTime, 0).Format("20060102_150405")),
			"patient_id":  patientIDStrs[i%5],
			"bed_id":      bedIDStrs[i%5],
			"is_resolved": i < 3, // 前3个已处理
			"created_at":  eventTime,
		})
	}

	// 创建一些情绪异常事件
	for i := 0; i < 3; i++ {
		eventTime := now.AddDate(0, 0, -i-1).Unix()

		analysisLogs = append(analysisLogs, map[string]interface{}{
			"type":        "event",
			"timestamp":   eventTime,
			"event_type":  "emotion_abnormal",
			"level":       "warning",
			"message":     fmt.Sprintf("检测到%s情绪异常", patientIDStrs[i%5]),
			"video_url":   fmt.Sprintf("http://82.156.64.69:9000/patient/evidence/emotion_%s.mp4", time.Unix(eventTime, 0).Format("20060102_150405")),
			"patient_id":  patientIDStrs[i%5],
			"bed_id":      bedIDStrs[i%5],
			"is_resolved": i < 2,
			"created_at":  eventTime,
		})
	}

	// 创建一些会话完成记录
	for i := 0; i < 3; i++ {
		eventTime := now.AddDate(0, 0, -i-2).Unix()

		analysisLogs = append(analysisLogs, map[string]interface{}{
			"type":       "session",
			"timestamp":  time.Unix(eventTime, 0).Format("2006-01-02 15:04:05"),
			"alert_type": "session_completed",
			"confidence": 0.95,
			"device_id":  fmt.Sprintf("camera_%d", i+1),
			"location":   fmt.Sprintf("房间%s", bedIDStrs[i%5][:4]),
			"details": map[string]interface{}{
				"source_video":    fmt.Sprintf("video_%d.mp4", i+1),
				"video_url":       fmt.Sprintf("http://82.156.64.69:9000/patient/evidence/session_%s.mp4", time.Unix(eventTime, 0).Format("20060102_150405")),
				"local_path":      fmt.Sprintf("/tmp/session_%d.mp4", i+1),
				"frames_analysed": 1500 + i*100,
				"timestamp_end":   time.Unix(eventTime+3600, 0).Format("2006-01-02 15:04:05"),
			},
			"created_at": eventTime,
		})
	}

	if len(analysisLogs) > 0 {
		_, err := collection.InsertMany(ctx, analysisLogs)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ 成功插入 %d 条分析日志数据\n", len(analysisLogs))
	return nil
}
