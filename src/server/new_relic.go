package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	apiMiddleware "github.com/michelaquino/golang_api_skeleton/src/middleware"
	newrelic "github.com/newrelic/go-agent"
	"github.com/spf13/viper"
)

// configureNewRelic is the method that enables New Relic.
func configureNewRelic(ctx context.Context, echoInstance *echo.Echo) {
	newRelicEnable := viper.GetBool("new_relic.is.enabled")
	if !newRelicEnable {
		return
	}

	newRelicApp, err := createNewRelicApp(ctx)
	if err != nil {
		logger.Error(ctx, "enabling New Relic", err.Error(), nil)
	}

	echoInstance.Use(apiMiddleware.NewRelicMiddleware(newRelicApp))
	logger.Info(ctx, "enabling New Relic", "success", nil)
}

// createNewRelicApp is the method that creates New Relic configuration.
func createNewRelicApp(ctx context.Context) (newrelic.Application, error) {
	licenseKeyEnvVar := viper.GetString("new_relic.licence.key")

	config := newrelic.NewConfig("My Awesome API", licenseKeyEnvVar)
	proxyURL, err := url.Parse(viper.GetString("new_relic.proxy.url"))
	if err != nil {
		logger.Error(ctx, "parse proxy url from env var", err.Error(), nil)
	}

	config.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	newRelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		logger.Error(ctx, "create New Relic APP ", err.Error(), nil)
		return nil, err
	}

	return newRelicApp, nil
}
