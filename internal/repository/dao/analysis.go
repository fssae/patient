package dao

import (
	"classroom-analysis/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type AnalysisDAO struct {
	BaseDAO[domain.AnalysisLog]
}

func NewAnalysisDAO(db *mongo.Database) *AnalysisDAO {
	return &AnalysisDAO{
		BaseDAO: NewBaseDAO[domain.AnalysisLog](db, "analysis"),
	}
}
