package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"fmt"

	"github.com/labstack/echo"
	"github.com/michelaquino/golang_api_skeleton/context"
	"github.com/michelaquino/golang_api_skeleton/handlers"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/middleware"
	"github.com/michelaquino/golang_api_skeleton/repository"
	newrelic "github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	context.GetMongoSession()
}

func main() {
	logger := context.GetLogger()

	echoInstance := echo.New()

	// Configure New Relic
	configureNewRelic(echoInstance)

	// Middlewares
	echoInstance.Use(apiMiddleware.RequestLogDataMiddleware())

	// Configure routes
	configureAllRoutes(echoInstance)

	testMysql()
	logger.Info("Main", "main", "", "", "start app", "success", "Started at port 8888!")
	echoInstance.Logger.Fatal(echoInstance.Start(":8888"))
}

func configureAllRoutes(echoInstance *echo.Echo) {
	// Metrics by Prometheus
	configureMetrics(echoInstance)

	// Healthcheck
	configureHealthcheckRoute(echoInstance)

	// User routes
	configureUserRoutes(echoInstance)
}

func configureHealthcheckRoute(echoInstance *echo.Echo) {
	echoInstance.GET("/healthcheck", handlers.Healthcheck)
}

func configureUserRoutes(echoInstance *echo.Echo) {
	userRepository := new(repository.UserMongoRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	userGroup := echoInstance.Group("/user")
	userGroup.POST("", userHandler.CreateUser)
}

// configureMetrics is the method that configure Prometheus metrics route
func configureMetrics(echoInstance *echo.Echo) {
	echoInstance.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

// configureNewRelic is the method that enable the new relic
func configureNewRelic(echoInstance *echo.Echo) {
	logger := context.GetLogger()

	newRelicEnvVar := os.Getenv("ENABLE_NEW_RELIC")
	newRelicEnable, err := strconv.ParseBool(newRelicEnvVar)

	if err == nil && newRelicEnable {
		if newRelicApp, err := createNewRelicApp(); err == nil {
			logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic ENABLED")
			echoInstance.Use(apiMiddleware.NewRelicMiddleware(newRelicApp))
		} else {
			logger.Error("Main", "configureNewRelic", "", "", "Enabling New Relic", "Error", err.Error())
		}

		return
	}

	if err != nil {
		logger.Error("Main", "configureNewRelic", "", "", "Enabling New Relic", "Error", err.Error())
	}

	logger.Info("Main", "configureNewRelic", "", "", "Enabling New Relic", "Success", "New Relic DISABLED")
}

// createNewRelicApp is the method that create new relic config
func createNewRelicApp() (newrelic.Application, error) {
	log := context.GetLogger()
	licenseKeyEnvVar := os.Getenv("NEW_RELIC_LICENSE_KEY")

	config := newrelic.NewConfig("My Awesome API", licenseKeyEnvVar)
	proxyURL, err := url.Parse(os.Getenv("NEW_RELIC_PROXY_URL"))
	if err != nil {
		log.Error("Main", "createNewRelicApp", "", "", "Parse proxy url from env var", "Error", err.Error())
	}

	config.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	newRelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		log.Error("Main", "createNewRelicApp", "", "", "Create New Relic APP ", "Error", err.Error())
		return nil, err
	}

	return newRelicApp, nil
}

func testMysql() {
	config := context.GetAPIConfig()

	mysqlConnectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/hello",
		config.MySQLConfig.User,
		config.MySQLConfig.Password,
		config.MySQLConfig.URL)

	db, err := sql.Open("mysql", mysqlConnectionString)
	if err != nil {
		fmt.Println("[MYSQL] Error on connect to database: ", err.Error())
		return
		// log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("[MYSQL] Error on PING: ", err.Error())
		return
		// do something here
	}

	fmt.Println("[MYSQL] Connected on MySQL With success")
	defer db.Close()
}
