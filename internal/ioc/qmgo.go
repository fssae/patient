package ioc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qiniu/qmgo"
)

// InitQmgoClient 初始化 Qmgo 客户端
func InitQmgoClient() *qmgo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := InitViper()
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.Mongodb.Account,
		cfg.Mongodb.Password,
		cfg.Mongodb.Address1,
		cfg.Mongodb.Port1)

	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		log.Fatal("Qmgo connection failed: ", err)
	}

	// Ping 测试连接
	if err = client.Ping(10); err != nil {
		log.Fatal("Qmgo ping failed: ", err)
	}

	return client
}

// InitQmgoDatabase 初始化 Qmgo 数据库连接
func InitQmgoDatabase(client *qmgo.Client) *qmgo.Database {
	return client.Database("kongdong")
}

// InitQmgoCareRecordCollection 初始化 CareRecord 集合
func InitQmgoCareRecordCollection(db *qmgo.Database) *qmgo.Collection {
	return db.Collection("care_records")
}
