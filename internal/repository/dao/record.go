package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordDAO struct {
	collection *mongo.Collection
}

func NewRecordDAO(db *mongo.Database) *RecordDAO {
	return &RecordDAO{
		collection: db.Collection("records"),
	}
}

// Create 创建登记记录
func (dao *RecordDAO) Create(ctx context.Context, record *domain.Record) error {
	record.CreatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	record.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找登记记录
func (dao *RecordDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Record, error) {
	var record domain.Record
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// FindByCustomerID 根据客户ID查找登记记录
func (dao *RecordDAO) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.Record, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Record
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindList 查询登记记录列表
func (dao *RecordDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Record, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}
	opts.SetSort(bson.D{{Key: "start_time", Value: -1}})

	cursor, err := dao.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Record
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新登记记录
func (dao *RecordDAO) Update(ctx context.Context, record *domain.Record) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": record.ID},
		bson.M{"$set": record},
	)
	return err
}

