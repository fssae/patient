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

type CareRecordDAO struct {
	collection *mongo.Collection
}

func NewCareRecordDAO(db *mongo.Database) *CareRecordDAO {
	return &CareRecordDAO{
		collection: db.Collection("care_records"),
	}
}

func (dao *CareRecordDAO) Create(ctx context.Context, record *domain.CareRecord) error {
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	record.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (dao *CareRecordDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecord, error) {
	var record domain.CareRecord
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (dao *CareRecordDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecord, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}
	opts.SetSort(bson.D{{Key: "care_time", Value: -1}})

	cursor, err := dao.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.CareRecord
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (dao *CareRecordDAO) Update(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	updates["updated_at"] = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

func (dao *CareRecordDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
