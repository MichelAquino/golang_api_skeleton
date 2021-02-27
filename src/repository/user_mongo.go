package repository

import (
	"context"

	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/michelaquino/golang_api_skeleton/src/models"
	"github.com/michelaquino/golang_api_skeleton/src/mongo"
)

var (
	userMongoCollectionName = "user"
	logger                  = log.GetLogger()
)

// UserMongoRepository is a user repository for MongoDB.
type UserMongoRepository struct{}

// Insert is a method that inserts an user into database.
func (u UserMongoRepository) Insert(ctx context.Context, userToInsert models.UserModel) error {
	// Execute Mongo's Insert
	err := mongo.Insert(ctx, userMongoCollectionName, &userToInsert)
	if err != nil {
		logger.Error(ctx, "create user", err.Error(), nil)
		return err
	}

	logger.Info(ctx, "create user", "Success", nil)
	return nil
}
