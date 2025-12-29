package ioc

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongodb() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取配置
	cfg := InitViper()
	// 构建连接字符串
	// 格式: mongodb://username:password@host:port
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.Mongodb.Account,
		cfg.Mongodb.Password,
		cfg.Mongodb.Address1,
		cfg.Mongodb.Port1)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
func InitMongoDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("kongdong")
}
