package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDAO struct {
	collection *mongo.Collection
}

func NewUserDAO(db *mongo.Database) *UserDAO {
	return &UserDAO{
		collection: db.Collection("users"),
	}
}

// Create 创建用户
func (dao *UserDAO) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByPhone 根据手机号查找用户
func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	err := dao.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindById 根据ID查找用户
func (dao *UserDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (dao *UserDAO) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

