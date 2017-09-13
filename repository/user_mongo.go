package repository

import (
	"github.com/michelaquino/golang_api_skeleton/context"
	"github.com/michelaquino/golang_api_skeleton/metrics"
	"github.com/michelaquino/golang_api_skeleton/models"

	"time"
)

// UserMongoRepository is a user repository for MongoDB
type UserMongoRepository struct{}

// Insert is a method that insert a user on database
func (u UserMongoRepository) Insert(requestLogData models.RequestLogData, userToInsert models.UserModel) error {
	log := context.GetLogger()

	dbSession := context.GetMongoSession()
	defer dbSession.Close()

	connection := dbSession.DB("api").C("user")

	// Now time to metrics
	now := time.Now()

	// Execute the insert
	err := connection.Insert(&userToInsert)

	// Send metrics to prometheus
	metrics.MongoDBDurationsSumary.WithLabelValues("insert").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("insert").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("UserMongoRepository", "Create", requestLogData.ID, requestLogData.OriginIP, "Create user", "Error", err.Error())
		return err
	}

	log.Info("UserMongoRepository", "Create", requestLogData.ID, requestLogData.OriginIP, "Create user", "Success", "")
	return nil
}
