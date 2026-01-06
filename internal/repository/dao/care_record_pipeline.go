package dao

import (
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CareRecordPipeline 护理记录聚合管道构建器
// 封装复杂的聚合逻辑，提供清晰的链式调用 API
type CareRecordPipeline struct {
	stages qmgo.Pipeline
}

// NewCareRecordPipeline 创建聚合管道构建器
func NewCareRecordPipeline() *CareRecordPipeline {
	return &CareRecordPipeline{
		stages: qmgo.Pipeline{},
	}
}

// MatchCustomer 匹配客户ID
func (p *CareRecordPipeline) MatchCustomer(customerID primitive.ObjectID) *CareRecordPipeline {
	p.stages = append(p.stages, bson.D{{Key: "$match", Value: bson.M{"customer_id": customerID}}})
	return p
}

// UnwindRecords 展开 records 数组
func (p *CareRecordPipeline) UnwindRecords() *CareRecordPipeline {
	p.stages = append(p.stages, bson.D{{Key: "$unwind", Value: "$records"}})
	return p
}

// FilterByTimeRange 按时间范围过滤
func (p *CareRecordPipeline) FilterByTimeRange(startTime, endTime *time.Time) *CareRecordPipeline {
	if startTime == nil && endTime == nil {
		return p
	}

	timeFilter := bson.M{}
	switch {
	case startTime != nil && endTime != nil:
		timeFilter["records.care_time"] = bson.M{"$gte": *startTime, "$lt": *endTime}
	case startTime != nil:
		timeFilter["records.care_time"] = bson.M{"$gte": *startTime}
	default:
		timeFilter["records.care_time"] = bson.M{"$lt": *endTime}
	}

	p.stages = append(p.stages, bson.D{{Key: "$match", Value: timeFilter}})
	return p
}

// SortByTimeDesc 按护理时间倒序
func (p *CareRecordPipeline) SortByTimeDesc() *CareRecordPipeline {
	p.stages = append(p.stages, bson.D{{Key: "$sort", Value: bson.D{{Key: "records.care_time", Value: -1}}}})
	return p
}

// Paginate 分页
func (p *CareRecordPipeline) Paginate(skip, limit int64) *CareRecordPipeline {
	if skip > 0 {
		p.stages = append(p.stages, bson.D{{Key: "$skip", Value: skip}})
	}
	if limit > 0 {
		p.stages = append(p.stages, bson.D{{Key: "$limit", Value: limit}})
	}
	return p
}

// GroupBack 重新组合成完整文档
func (p *CareRecordPipeline) GroupBack() *CareRecordPipeline {
	p.stages = append(p.stages, bson.D{{Key: "$group", Value: bson.M{
		"_id":           "$_id",
		"customer_id":   bson.M{"$first": "$customer_id"},
		"customer_name": bson.M{"$first": "$customer_name"},
		"updated_at":    bson.M{"$first": "$updated_at"},
		"records":       bson.M{"$push": "$records"},
	}}})
	return p
}

// Count 添加计数阶段
func (p *CareRecordPipeline) Count() *CareRecordPipeline {
	p.stages = append(p.stages, bson.D{{Key: "$count", Value: "total"}})
	return p
}

// Build 构建最终的 Pipeline
func (p *CareRecordPipeline) Build() qmgo.Pipeline {
	return p.stages
}

// Clone 克隆当前管道（用于创建计数管道）
func (p *CareRecordPipeline) Clone() *CareRecordPipeline {
	cloned := make(qmgo.Pipeline, len(p.stages))
	copy(cloned, p.stages)
	return &CareRecordPipeline{stages: cloned}
}
