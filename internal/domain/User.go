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
type BedResponse struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RoomID       primitive.ObjectID `json:"room_id" bson:"room_id" binding:"required"`
	RoomNumber   string             `json:"room_number"`
	Number       string             `json:"number" bson:"number" binding:"required"`            // 床位号，如 "A101-1"
	Status       string             `json:"status" bson:"status"`                               // 状态："空闲"/"占用"/"维护中"
	CustomerID   primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"` // 当前入住客户ID
	CustomerName string             `json:"customer_name"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
type CreateBed struct {
	RoomId string `json:"room_id"`
	Number string `json:"number"`
	Status string `json:"status"`
}

type CustomerNameID struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
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

	// 健康状况
	HealthLevel    string `json:"health_level" bson:"health_level"`       // 健康等级
	MedicalHistory string `json:"medical_history" bson:"medical_history"` // 既往病史
	Medication     string `json:"medication" bson:"medication"`           // 用药情况
	AllergyHistory string `json:"allergy_history" bson:"allergy_history"` // 过敏史

	// 紧急联系人
	ContactName    string `json:"contact_name" bson:"contact_name"`       // 联系人姓名
	Relationship   string `json:"relationship" bson:"relationship"`       // 关系
	ContactPhone   string `json:"contact_phone" bson:"contact_phone"`     // 联系电话
	ContactAddress string `json:"contact_address" bson:"contact_address"` // 联系地址

	Remarks string `json:"remarks" bson:"remarks"` // 备注说明

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
type CustomerResponse struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"` // 关联注册用户
	Name            string             `json:"name" bson:"name" binding:"required"`
	Age             int                `json:"age" bson:"age" binding:"required"`
	Gender          string             `json:"gender" bson:"gender" binding:"required"`
	Phone           string             `json:"phone" bson:"phone"`
	IDCard          string             `json:"id_card" bson:"id_card"`                   // 身份证号
	BedID           primitive.ObjectID `json:"bed_id,omitempty" bson:"bed_id,omitempty"` // 关联床位
	Bed             string             `json:"bed_name,omitempty" bson:"bed_name,omitempty"`
	DietPlanID      primitive.ObjectID `json:"diet_plan_id,omitempty" bson:"diet_plan_id,omitempty"` // 膳食计划ID
	DietPlan        string             `json:"diet_plan_name,omitempty" bson:"diet_plan_name,omitempty"`
	CareLevelID     primitive.ObjectID `json:"care_level_id,omitempty" bson:"care_level_id,omitempty"` // 护理级别ID
	CareLevel       string             `json:"care_level_name,omitempty" bson:"care_level_name,omitempty"`
	HealthManagerID primitive.ObjectID `json:"health_manager_id,omitempty" bson:"health_manager_id,omitempty"` // 健康管家ID
	HealthManager   string             `json:"health_manager" bson:"health_manager"`                           // 健康管家姓名
	Status          string             `json:"status" bson:"status"`                                           // 状态："入住中"/"已退住"/"外出中"
	CheckInDate     time.Time          `json:"check_in_date,omitempty" bson:"check_in_date,omitempty"`         // 入住日期
	CheckOutDate    time.Time          `json:"check_out_date,omitempty" bson:"check_out_date,omitempty"`       // 退住日期

	// 健康状况
	HealthLevel    string `json:"health_level" bson:"health_level"`       // 健康等级
	MedicalHistory string `json:"medical_history" bson:"medical_history"` // 既往病史
	Medication     string `json:"medication" bson:"medication"`           // 用药情况
	AllergyHistory string `json:"allergy_history" bson:"allergy_history"` // 过敏史

	// 紧急联系人
	ContactName    string `json:"contact_name" bson:"contact_name"`       // 联系人姓名
	Relationship   string `json:"relationship" bson:"relationship"`       // 关系
	ContactPhone   string `json:"contact_phone" bson:"contact_phone"`     // 联系电话
	ContactAddress string `json:"contact_address" bson:"contact_address"` // 联系地址

	Remarks string `json:"remarks" bson:"remarks"` // 备注说明

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CustomerQuery 客户查询条件
type CustomerQuery struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	IDCard          string `json:"id_card"`
	Status          string `json:"status"`
	BedID           string `json:"bed_id"`
	CareLevelID     string `json:"care_level_id"`
	DietPlanID      string `json:"diet_plan_id"`
	HealthManagerID string `json:"health_manager_id"`
	MinAge          int    `json:"min_age"`
	MaxAge          int    `json:"max_age"`
	SearchKey       string `json:"search_key"` // 万能搜索框：匹配姓名或手机号
}

// 入住老人
type ElderlyRegisterRequest struct {
	// 基本信息
	Name        string `json:"name"`                       // 姓名
	Gender      string `json:"gender" binding:"required"`  // 性别 (男/女)
	Age         int    `json:"age" binding:"required"`     // 年龄
	IDCard      string `json:"id_card" binding:"required"` // 身份证号
	PhoneNumber string `json:"phone_number"`               // 联系电话
	HomeAddress string `json:"home_address"`               // 家庭住址

	// 健康状况记录
	HealthLevel    string `json:"health_level" binding:"required"` // 健康等级 (健康/较好/一般/较差)
	MedicalHistory string `json:"medical_history"`                 // 既往病史
	Medication     string `json:"medication"`                      // 用药情况
	AllergyHistory string `json:"allergy_history"`                 // 过敏史

	// 紧急联系人信息
	ContactName    string `json:"contact_name" binding:"required"`  // 联系人姓名
	Relationship   string `json:"relationship" binding:"required"`  // 关系
	ContactPhone   string `json:"contact_phone" binding:"required"` // 联系电话
	ContactAddress string `json:"contact_address"`                  // 联系地址

	// 入住安排详情
	CheckInDate  string `json:"check_in_date"`                    // 入住日期 (建议格式 "2025-01-01")
	BedID        string `json:"bed_id" binding:"required"`        // 床位分配
	NursingLevel string `json:"nursing_level" binding:"required"` // 护理级别
	DietaryType  string `json:"dietary_type" binding:"required"`  // 膳食类型

	Remarks string `json:"remarks"` // 备注说明
}

// DietPlanNameID 膳食计划名称ID
type DietPlanNameID struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
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
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID         primitive.ObjectID `json:"customer_id" bson:"customer_id" binding:"required"`
	Type               string             `json:"type" bson:"type" binding:"required"`             // "入住"/"退住"/"外出"
	StartTime          time.Time          `json:"start_time" bson:"start_time" binding:"optional"` // 入住时间或外出时间
	EndTime            time.Time          `json:"end_time,omitempty" bson:"end_time,omitempty"`    // 退住时间或外出返回时间
	Note               string             `json:"note" bson:"note"`
	CreatedBy          string             `json:"created_by" bson:"created_by"` // 操作人
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	EmergencyContact   string             `json:"emergency_contact" bson:"emergency_contact"`                           // 紧急联系人姓名
	Destination        string             `json:"destination,omitempty" bson:"destination,omitempty"`                   // 外出目的地
	Escort             string             `json:"escort,omitempty" bson:"escort,omitempty"`                             // 陪护人姓名
	ExpectedReturnTime time.Time          `json:"expected_return_time,omitempty" bson:"expected_return_time,omitempty"` // 预计归来时间
	Remark             string             `json:"remark,omitempty" bson:"remark,omitempty"`                             // 备注说明
}

// Record 更新记录请求体
type UpdateRecordRequest struct {
	ID                 string `json:"id"`
	CustomerID         string `json:"customerId"`         // 客户ID
	ElderID            string `json:"elderId"`            // 老人姓名
	EmergencyContact   string `json:"emergencyContact"`   // 紧急联系电话
	ExpectedReturnTime string `json:"expectedReturnTime"` // 预计返回时间
	Destination        string `json:"destination"`        // 目的地
	Escort             string `json:"escort"`             // 陪同人员
	Remark             string `json:"remark"`             // 备注
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

// CustomerServiceResponse 客户购买的服务（包含关联信息）
type CustomerServiceResponse struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	CustomerID   primitive.ObjectID `json:"customer_id"`
	CustomerName string             `json:"customer_name"` // 客户姓名
	ServiceID    primitive.ObjectID `json:"service_id"`
	ServiceName  string             `json:"service_name"`        // 服务名称
	ServiceDesc  string             `json:"service_description"` // 服务描述
	Category     string             `json:"category"`            // 服务类别
	Price        float64            `json:"price"`               // 服务价格
	Unit         string             `json:"unit"`                // 计价单位
	StartDate    time.Time          `json:"start_date"`
	EndDate      time.Time          `json:"end_date,omitempty"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
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

type HealthManagerQuery struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Specialty string `json:"specialty"`
}

// 请求结构体定义，用于 Swagger 文档生成

type SetHealthManagerRequest struct {
	ManagerID   string `json:"manager_id" binding:"required"`
	ManagerName string `json:"manager_name" binding:"required"`
}

type SetBedRequest struct {
	BedID string `json:"bed_id" binding:"required"`
}

type SetDietPlanRequest struct {
	DietPlanID string `json:"diet_plan_id" binding:"required"`
}

type SetCareLevelRequest struct {
	CareLevelID string `json:"care_level_id" binding:"required"`
}

type CheckOutRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Note       string `json:"note"`
	CreatedBy  string `json:"created_by" binding:"required"`
}

type OutgoingRequest struct {
	CustomerID         string `json:"elder_id"`
	Note               string `json:"note"`
	CustomerName       string `json:"customer_name"`
	CreatedBy          string `json:"created_by"`
	Destination        string `json:"destination" binding:"required"`
	EmergencyContact   string `json:"emergencycontact" binding:"required"`
	Escort             string `json:"escort" binding:"required"`
	ExpectedReturnTime string `json:"expectedreturntime" binding:"required"`
	OutTime            string `json:"outTime" binding:"required"`
	Remark             string `json:"remark"`
}

type ReturnRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Note       string `json:"note"`
	CreatedBy  string `json:"created_by" binding:"required"`
}

type EndServiceRequest struct {
	EndDate string `json:"end_date" binding:"required"`
}

type PurchaseServiceRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	ServiceID  string `json:"service_id" binding:"required"`
	StartDate  string `json:"start_date" binding:"required"`
}
