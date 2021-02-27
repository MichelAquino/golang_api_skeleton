package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apiMiddleware "github.com/michelaquino/golang_api_skeleton/src/middleware"
	"github.com/michelaquino/golang_api_skeleton/src/models"
	"github.com/michelaquino/golang_api_skeleton/src/repository"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	userRepositoryMock *repository.UserRepositoryMock
	userHandler        *UserHandler
	serverMock         *httptest.Server
	requestLogDataMock models.RequestLogData
)

func setupUserHandlerTest(t *testing.T) {
	userRepositoryMock = &repository.UserRepositoryMock{}
	userHandler = NewUserHandler(userRepositoryMock)

	serverMock = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "response")
		}))

	requestLogDataMock = models.RequestLogData{
		ID:       "RequestIDMock",
		OriginIP: "187.298.0.123",
	}
}

func Test_CreateUser_ShouldReturnStatusInternalServerErrorWhenRepositoryReturnError(t *testing.T) {
	setupUserHandlerTest(t)

	bodyPayload := strings.NewReader(`{"name": "123456", "email": "89978"}`)
	recorder, echoContext := getTestBaseObjects(bodyPayload)

	userRepositoryMock.On("Insert", mock.Anything, mock.Anything).Return(errors.New("Unexpected error"))

	userHandler.CreateUser(echoContext)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func Test_CreateUser_ShouldReturnStatusBadRequestWhenBindPayloadToModel(t *testing.T) {
	setupUserHandlerTest(t)

	invalidBodyPayload := strings.NewReader(`invalid json`)
	recorder, echoContext := getTestBaseObjects(invalidBodyPayload)

	userRepositoryMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

	userHandler.CreateUser(echoContext)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func Test_CreateUser_ShouldReturnStatusCreated(t *testing.T) {
	setupUserHandlerTest(t)

	bodyPayload := strings.NewReader(`{"name": "123456", "email": "89978"}`)
	recorder, echoContext := getTestBaseObjects(bodyPayload)

	userRepositoryMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

	userHandler.CreateUser(echoContext)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func getTestBaseObjects(body io.Reader) (*httptest.ResponseRecorder, echo.Context) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, serverMock.URL, body)
	request.Header.Add("Content-Type", "application/json")

	echoInstance := echo.New()
	echoContext := echoInstance.NewContext(request, recorder)

	echoContext.SetRequest(request)
	echoContext.Set(apiMiddleware.RequestIDKey, requestLogDataMock)

	return recorder, echoContext
}
