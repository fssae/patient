package ioc

import (
	"classroom-analysis/internal/repository"
	"classroom-analysis/internal/repository/dao"
	"classroom-analysis/internal/service"
	"classroom-analysis/internal/web"

	"github.com/google/wire"
)

var MinimalSet = wire.NewSet(
	InitMongodb,
	InitMongoDatabase,
	InitRedis,
	InitLogger,
	NewLogger,
	NewMinioClient,

	// DAO层

	dao.NewFileDao,
	dao.ProvideFileDaoInterface,
	dao.NewPatientDAO,
	dao.NewApiDao,
	dao.NewUserDAO,
	dao.NewHealthManagerDAO,
	dao.NewRoomDAO,
	dao.NewBedDAO,
	dao.NewCareLevelDAO,
	dao.NewDietPlanDAO,
	dao.NewCustomerDAO,
	dao.NewRecordDAO,
	dao.NewServiceDAO,
	dao.NewCustomerServiceDAO,

	// Repository层

	repository.NewFileRepository,
	repository.ProvideFileInterface,
	repository.NewPatientRepository,
	repository.NewUserRepository,
	repository.NewHealthManagerRepository,
	repository.NewRoomRepository,
	repository.NewBedRepository,
	repository.NewCareLevelRepository,
	repository.NewDietPlanRepository,
	repository.NewCustomerRepository,
	repository.NewRecordRepository,
	repository.NewServiceRepository,
	repository.NewCustomerServiceRepository,

	// Service层
	service.NewFileService,
	service.NewPatientService,
	service.NewUserService,
	service.NewHealthManagerService,
	service.NewRoomService,
	service.NewBedService,
	service.NewCareLevelService,
	service.NewDietPlanService,
	service.NewCustomerService,
	service.NewRecordService,
	service.NewServiceService,

	// Kafka相关
	//InitKafkaWriter,
	// Web层
	web.NewFileHandler,
	web.NewPatientHandler,
	web.NewUserHandler,
	web.NewHealthManagerHandler,
	web.NewRoomHandler,
	web.NewBedHandler,
	web.NewCareLevelHandler,
	web.NewDietPlanHandler,
	web.NewCustomerHandler,
	web.NewRecordHandler,
	web.NewServiceHandler,

	// Gin引擎
	InitGin,
)
