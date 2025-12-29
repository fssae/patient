package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户（注册用户）
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Phone     string             `json:"phone" bson:"phone" binding:"required"`       // 手机号码，唯一标识
	Password  string             `json:"password" bson:"password" binding:"required"` // 登录密码，加密存储
	Name      string             `json:"name" bson:"name" binding:"required"`         // 真实姓名
	Age       int                `json:"age" bson:"age" binding:"required"`           // 年龄
	Gender    string             `json:"gender" bson:"gender" binding:"required"`     // 性别："男"/"女"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// Room 房间
type Room struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Number      string             `json:"number" bson:"number" binding:"required"` // 房间号，如 "A101"
	Floor       int                `json:"floor" bson:"floor"`                      // 楼层
	Type        string             `json:"type" bson:"type"`                        // 房间类型："单人间"/"双人间"/"多人间"
	Capacity    int                `json:"capacity" bson:"capacity"`                // 床位容量
	Status      string             `json:"status" bson:"status"`                    // 状态："可用"/"已满"/"维护中"
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Bed 床位
type Bed struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RoomID     primitive.ObjectID `json:"room_id" bson:"room_id" binding:"required"`
	Number     string             `json:"number" bson:"number" binding:"required"`            // 床位号，如 "A101-1"
	Status     string             `json:"status" bson:"status"`                               // 状态："空闲"/"占用"/"维护中"
	CustomerID primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"` // 当前入住客户ID
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
type CreateBedRequest struct {
	// 基本信息
	Name    string `json:"name" bson:"name" binding:"required"`
	Age     int    `json:"age" bson:"age" binding:"required"`
	Gender  string `json:"gender" bson:"gender" binding:"required"`
	IDCard  string `json:"id_card" bson:"id_card" binding:"required"`
	Phone   string `json:"phone" bson:"phone" binding:"required"`
	Address string `json:"address" bson:"address"`

	// 关联信息 (使用 primitive.ObjectID 对应 MongoDB 的 $oid)
	BedID primitive.ObjectID `json:"bed_id" bson:"bed_id"`

	// 业务等级
	CareLevel   string `json:"care_level" bson:"care_level" binding:"required"`
	HealthLevel string `json:"health_level" bson:"health_level" binding:"required"`
	DietType    string `json:"diet_type" bson:"diet_type" binding:"required"`

	// 状态与时间
	Status string `json:"status" bson:"status" binding:"required"`
	// 数据库通常存 time.Time，这里接收 int64 时间戳，写入 DAO 时需要转换
	CheckInTime int64 `json:"check_in_time" bson:"check_in_time"`

	// 其他补充
	EmergencyPhone string `json:"emergency_phone" bson:"emergency_phone"`
	Medication     string `json:"medication" bson:"medication"`
	Allergies      string `json:"allergies" bson:"allergies"`

	// 自动生成的字段通常在 DAO 层处理，不一定非要放在 Request 结构体里
	CreatedAt time.Time `json:"-" bson:"created_at"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
}

// Customer 客户（入住老人）
type Customer struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"` // 关联注册用户
	Name            string             `json:"name" bson:"name" binding:"required"`
	Age             int                `json:"age" bson:"age" binding:"required"`
	Gender          string             `json:"gender" bson:"gender" binding:"required"`
	Phone           string             `json:"phone" bson:"phone"`
	IDCard          string             `json:"id_card" bson:"id_card"`                                         // 身份证号
	BedID           primitive.ObjectID `json:"bed_id,omitempty" bson:"bed_id,omitempty"`                       // 关联床位
	DietPlanID      primitive.ObjectID `json:"diet_plan_id,omitempty" bson:"diet_plan_id,omitempty"`           // 膳食计划ID
	CareLevelID     primitive.ObjectID `json:"care_level_id,omitempty" bson:"care_level_id,omitempty"`         // 护理级别ID
	HealthManagerID primitive.ObjectID `json:"health_manager_id,omitempty" bson:"health_manager_id,omitempty"` // 健康管家ID
	HealthManager   string             `json:"health_manager" bson:"health_manager"`                           // 健康管家姓名
	Status          string             `json:"status" bson:"status"`                                           // 状态："入住中"/"已退住"/"外出中"
	CheckInDate     time.Time          `json:"check_in_date,omitempty" bson:"check_in_date,omitempty"`         // 入住日期
	CheckOutDate    time.Time          `json:"check_out_date,omitempty" bson:"check_out_date,omitempty"`       // 退住日期
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// DietPlan 膳食计划
type DietPlan struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"` // 如"低糖餐"、"流质餐"
	Description string             `json:"description" bson:"description"`
	WeekMenu    []WeekDayMenu      `json:"week_menu" bson:"week_menu"` // 每周菜单
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// WeekDayMenu 每日菜单
type WeekDayMenu struct {
	Day       string `json:"day" bson:"day"`                         // "周一"、"周二"等
	Breakfast string `json:"breakfast" bson:"breakfast"`             // 早餐
	Lunch     string `json:"lunch" bson:"lunch"`                     // 午餐
	Dinner    string `json:"dinner" bson:"dinner"`                   // 晚餐
	Snack     string `json:"snack,omitempty" bson:"snack,omitempty"` // 加餐
}

// CareLevel 护理级别
type CareLevel struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"` // 如"一级护理"、"二级护理"
	Level       int                `json:"level" bson:"level"`                  // 级别数字：1、2、3等
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"` // 护理内容描述
	Price       float64            `json:"price" bson:"price"`     // 护理费用
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Record 登记记录（入住/退住/外出）
type Record struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id" bson:"customer_id" binding:"required"`
	Type       string             `json:"type" bson:"type" binding:"required"` // "入住"/"退住"/"外出"
	StartTime  time.Time          `json:"start_time" bson:"start_time" binding:"required"`
	EndTime    time.Time          `json:"end_time,omitempty" bson:"end_time,omitempty"` // 退住时间或外出返回时间
	Note       string             `json:"note" bson:"note"`
	CreatedBy  string             `json:"created_by" bson:"created_by"` // 操作人
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

// Service 服务项目
type Service struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"` // 服务名称
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"` // 服务类别："医疗"/"生活"/"娱乐"等
	Price       float64            `json:"price" bson:"price"`       // 服务价格
	Unit        string             `json:"unit" bson:"unit"`         // 计价单位："次"/"月"/"年"
	Status      string             `json:"status" bson:"status"`     // 状态："启用"/"停用"
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CustomerService 客户购买的服务
type CustomerService struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID  primitive.ObjectID `json:"customer_id" bson:"customer_id" binding:"required"`
	ServiceID   primitive.ObjectID `json:"service_id" bson:"service_id" binding:"required"`
	ServiceName string             `json:"service_name" bson:"service_name"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	EndDate     time.Time          `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Status      string             `json:"status" bson:"status"` // "进行中"/"已结束"/"已取消"
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ServiceAttention 服务关注（服务对象设置）
type ServiceAttention struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id" bson:"customer_id" binding:"required"`
	ServiceID  primitive.ObjectID `json:"service_id" bson:"service_id" binding:"required"`
	Priority   int                `json:"priority" bson:"priority"` // 优先级：1-高、2-中、3-低
	Note       string             `json:"note" bson:"note"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// HealthManager 健康管家
type HealthManager struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Phone     string             `json:"phone" bson:"phone"`
	Email     string             `json:"email" bson:"email"`
	Specialty string             `json:"specialty" bson:"specialty"` // 专长
	Status    string             `json:"status" bson:"status"`       // "在职"/"离职"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
