package repository

import (
	"context"

	"github.com/michelaquino/golang_api_skeleton/src/models"
)

// UserRepository is a interface that defines methods to user's repository.
type UserRepository interface {

	// Insert is a method that inserts an user into database.
	Insert(ctx context.Context, userToInsert models.UserModel) error
}
