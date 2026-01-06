package ioc

import (
    "classroom-analysis/internal/mq"
    "classroom-analysis/internal/repository"
    "classroom-analysis/internal/repository/dao"
    "classroom-analysis/internal/service"
    "classroom-analysis/internal/web"
    "classroom-analysis/internal/ws"

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
    dao.NewCareRecordDAO,
    dao.NewAnalysisDAO,

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
    repository.NewCareRecordRepository,
    repository.NewAnalysisRepository,

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
    service.NewCareRecordService,
    service.NewAnalysisService,
    service.NewAlertService,
    service.NewStatsService,
    service.NewJWTService,
    service.NewJWTBlacklistService,

    // MQ
    mq.NewAlertConsumer,
    // WS
    ws.NewWebSocketManager,

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
    web.NewCareRecordHandler,
    web.NewAnalysisHandler,
    web.NewStatsHandler,
    web.NewAlertHandler,
    web.NewNotificationHandler,
    web.NewCompatHandler,
    web.NewAuthHandler,

    // Middleware
    InitJWTMiddleware,

    // Gin引擎
    InitGin,
)
