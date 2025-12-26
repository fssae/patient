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

	// Repository层

	repository.NewFileRepository,
	repository.ProvideFileInterface,
	repository.NewPatientRepository,

	// Service层
	service.NewFileService,
	service.NewPatientService,

	// Kafka相关
	InitKafkaWriter,
	// Web层 - 添加FileHandler
	web.NewFileHandler,
	web.NewPatientHandler,

	// Gin引擎
	InitGin,
)
