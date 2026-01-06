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

func (dao *CareRecordDAO) Create(ctx context.Context, record *domain.CareRecords) error {
	now := time.Now()
	record.UpdatedAt = now
	for i := range record.Records {
		if record.Records[i].CreatedAt.IsZero() {
			record.Records[i].CreatedAt = now
		}
	}
	result, err := dao.collection.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	record.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (dao *CareRecordDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecords, error) {
	var record domain.CareRecords
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (dao *CareRecordDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecords, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}
	opts.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := dao.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.CareRecords
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

func (dao *CareRecordDAO) AppendRecord(ctx context.Context, id primitive.ObjectID, item domain.RecordItems) error {
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$push": bson.M{"records": item},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// FindRecordsByCustomerId 根据 CustomerId 查询护理记录，按 care_time 过滤
// 返回完整的 CareRecords 结构（records 只包含符合条件的记录）和符合条件的记录总数
func (dao *CareRecordDAO) FindRecordsByCustomerId(ctx context.Context, customerID primitive.ObjectID, startTime, endTime *time.Time, skip, limit int64) (*domain.CareRecords, int64, error) {
	// 首先获取文档基本信息
	var baseRecord domain.CareRecords
	err := dao.collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&baseRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	// 构建聚合管道来过滤和分页 records
	pipeline := mongo.Pipeline{}

	// 1. 匹配 customer_id
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"customer_id": customerID}}})

	// 2. 展开 records 数组
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: "$records"}})

	// 3. 按 care_time 过滤
	if startTime != nil || endTime != nil {
		matchFilter := bson.M{}
		if startTime != nil && endTime != nil {
			matchFilter["records.care_time"] = bson.M{"$gte": *startTime, "$lt": *endTime}
		} else if startTime != nil {
			matchFilter["records.care_time"] = bson.M{"$gte": *startTime}
		} else if endTime != nil {
			matchFilter["records.care_time"] = bson.M{"$lt": *endTime}
		}
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchFilter}})
	}

	// 4. 按 care_time 倒序排序
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "records.care_time", Value: -1}}}})

	// 复制 pipeline 用于计数（在分页之前）
	countPipeline := make(mongo.Pipeline, len(pipeline))
	copy(countPipeline, pipeline)
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	// 5. 分页
	if skip > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})
	}

	// 6. 重新组合成完整文档结构
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id":           "$_id",
		"customer_id":   bson.M{"$first": "$customer_id"},
		"customer_name": bson.M{"$first": "$customer_name"},
		"updated_at":    bson.M{"$first": "$updated_at"},
		"records":       bson.M{"$push": "$records"},
	}}})

	// 执行聚合查询
	cursor, err := dao.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []domain.CareRecords
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	// 获取总数
	var total int64 = 0
	countCursor, err := dao.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// 如果没有符合条件的记录，返回基本信息但 records 为空
	if len(results) == 0 {
		baseRecord.Records = []domain.RecordItems{}
		return &baseRecord, total, nil
	}

	return &results[0], total, nil
}
