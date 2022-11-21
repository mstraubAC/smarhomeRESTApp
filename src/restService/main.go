package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mstraubAC/smarhomeRESTApp/src/restService/accessors"
	"github.com/mstraubAC/smarhomeRESTApp/src/restService/configuration"
	"github.com/mstraubAC/smarhomeRESTApp/src/restService/controllers/aggregates"
	"github.com/mstraubAC/smarhomeRESTApp/src/restService/controllers/locations"
	"github.com/mstraubAC/smarhomeRESTApp/src/restService/middleware"
	"go.uber.org/zap"
)

func main() {
	// logger
	logger, _ := zap.NewDevelopment()
	logger.Debug("Test")

	// configuration
	config, err := configuration.LoadConfig(*logger)
	if err != nil {
		logger.Fatal("Failed to read config: " + err.Error())
		logger.Fatal("Terminating")
		return
	}
	logger.Debug(config.RestListener)

	// setup database
	databaseAccessor := accessors.DatabaseAccessor{Config: &config, Logger: logger}
	databaseAccessor.GetSqlConnection()

	// setup routing middleware
	router := gin.Default()
	router.Use(middleware.ZapLoggingHandler(logger))
	router.Use(middleware.ErrorHandler(logger))

	// register routes
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(200, config)
	})

	// v1 API
	locations.RegisterRoutes(router.Group("/v1"), &config, logger, &databaseAccessor)
	aggregates.RegisterRoutes(router.Group("/v1"), &config, logger, &databaseAccessor)

	// startup router
	router.Run(config.RestListener)
}