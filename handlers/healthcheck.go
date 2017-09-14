package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/context"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/middleware"
	"github.com/michelaquino/golang_api_skeleton/models"
)

// Healthcheck is a method that respond only WORKING
func Healthcheck(echoContext echo.Context) error {
	healthcheckHandlerLog := context.GetLogger()
	requestLogData := echoContext.Get(apiMiddleware.RequestIDKey).(models.RequestLogData)

	healthcheckHandlerLog.Info("Handlers", "Healthcheck", requestLogData.ID, requestLogData.OriginIP, "Verify Healthcheck", "success", "")

	testMysql()

	return echoContext.String(http.StatusOK, "WORKING")
}

func testMysql() {
	log := context.GetLogger()

	err := context.CheckMysqlConn()
	if err != nil {
		log.Error("Main", "testMysql", "", "", "", "Error", err.Error())
		return
	}

	log.Info("Main", "testMysql", "", "", "", "Success", "")
}
