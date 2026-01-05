package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalysisRepository interface {
	Save(ctx context.Context, log *domain.AnalysisLog) error
	GetList(ctx context.Context, page, size int64) ([]*domain.AnalysisLog, int64, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.AnalysisLog, error)
	UpdateResolveStatus(ctx context.Context, id primitive.ObjectID, resolved bool) error
}

type analysisRepository struct {
	dao *dao.AnalysisDAO
}

func NewAnalysisRepository(dao *dao.AnalysisDAO) AnalysisRepository {
	return &analysisRepository{dao: dao}
}

func (r *analysisRepository) Save(ctx context.Context, log *domain.AnalysisLog) error {
	_, err := r.dao.InsertOne(ctx, log)
	return err
}

func (r *analysisRepository) GetList(ctx context.Context, page, size int64) ([]*domain.AnalysisLog, int64, error) {
	skip := (page - 1) * size
	// Sort by created_at descending
	sort := bson.D{{Key: "created_at", Value: -1}}
	return r.dao.FindList(ctx, bson.M{}, skip, size, sort)
}

func (r *analysisRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.AnalysisLog, error) {
	return r.dao.FindOne(ctx, bson.M{"_id": id})
}

func (r *analysisRepository) UpdateResolveStatus(ctx context.Context, id primitive.ObjectID, resolved bool) error {
	_, err := r.dao.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_resolved": resolved}})
	return err
}
