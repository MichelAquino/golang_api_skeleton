package repository

import "github.com/michelaquino/golang_api_skeleton/models"

// UserRepository is a interface that defines methods to user's repository.
type UserRepository interface {

	// Insert is a method that inserts an user into database.
	Insert(requestLogData models.RequestLogData, userToInsert models.UserModel) error
}
