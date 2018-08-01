package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/context"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/middleware"
	"github.com/michelaquino/golang_api_skeleton/models"
	"github.com/michelaquino/golang_api_skeleton/repository"
)

// UserHandler is a struct that stores an userRepository.
type UserHandler struct {
	userRepository repository.UserRepository
}

// NewUserHandler returns a new pointer of user's struct.
func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

// CreateUser is a handler that creates a new user into database.
func (h UserHandler) CreateUser(echoContext echo.Context) error {
	userHandlerLog := context.GetLogger()
	requestLogData := echoContext.Get(apiMiddleware.RequestIDKey).(models.RequestLogData)

	userModel := models.UserModel{}
	if err := echoContext.Bind(&userModel); err != nil {
		userHandlerLog.Error("UserHandler", "CreateUser", requestLogData.ID, requestLogData.OriginIP, "Bind payload to model", "Error", err.Error())
		return echoContext.NoContent(http.StatusBadRequest)
	}

	if err := h.userRepository.Insert(requestLogData, userModel); err != nil {
		userHandlerLog.Error("UserHandler", "CreateUser", requestLogData.ID, requestLogData.OriginIP, "Create user", "Error", err.Error())
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	userHandlerLog.Info("UserHandler", "CreateUser", requestLogData.ID, requestLogData.OriginIP, "Create user", "Success", "")
	return echoContext.NoContent(http.StatusCreated)
}
