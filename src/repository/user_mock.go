package repository

import (
	"github.com/michelaquino/golang_api_skeleton/src/models"
	testifyMock "github.com/stretchr/testify/mock"
)

// UserRepositoryMock is a mock that implements UserRepository interface.
type UserRepositoryMock struct {
	testifyMock.Mock
}

// Insert is a method that inserts an user into the database.
func (mock *UserRepositoryMock) Insert(requestLogData models.RequestLogData, userToInsert models.UserModel) error {
	args := mock.Called(requestLogData, userToInsert)
	return args.Error(0)
}
