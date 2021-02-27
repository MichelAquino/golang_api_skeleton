package mongo

import (
	"context"
	"fmt"
	"time"

	appContext "github.com/michelaquino/golang_api_skeleton/src/context"
	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/michelaquino/golang_api_skeleton/src/metrics"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoDatabaseName = "api-skeleton"
)

var (
	logger = log.GetLogger()
)

func Insert(ctx context.Context, collection string, objectToInsert interface{}) error {
	// Now time for metrics
	now := time.Now()

	session := appContext.GetMongoSession()
	defer session.Close()

	logAction := fmt.Sprintf("Inserting object in collection %s", collection)
	logger.Info(ctx, logAction, "", nil)

	connection := session.DB(mongoDatabaseName).C(collection)
	err := connection.Insert(&objectToInsert)

	// Send metrics to Prometheus
	metrics.MongoDBDurationsSumary.WithLabelValues("Insert").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Insert").Observe(time.Since(now).Seconds())

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Info(ctx, logAction, "object inserted with success", nil)
	return nil
}

func FindOne(ctx context.Context, collection string, query bson.M, object interface{}) error {
	// Now time for metrics
	now := time.Now()

	session := appContext.GetMongoSession()
	defer session.Close()

	logAction := fmt.Sprintf("getting object in collection %s", collection)
	logger.Info(ctx, logAction, "", nil)
	connection := session.DB(mongoDatabaseName).C(collection)

	err := connection.Find(query).One(object)

	metrics.MongoDBDurationsSumary.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Info(ctx, logAction, "object got with success", nil)
	return nil
}

func FindAll(ctx context.Context, collection string, query bson.M) ([]interface{}, error) {
	// Now time for metrics
	now := time.Now()

	session := appContext.GetMongoSession()
	defer session.Close()

	var objectList []interface{}

	logAction := fmt.Sprintf("getting object list in collection %s", collection)
	logger.Info(ctx, logAction, "", nil)

	connection := session.DB(mongoDatabaseName).C(collection)

	err := connection.Find(query).All(&objectList)

	metrics.MongoDBDurationsSumary.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return nil, err
	}

	logger.Info(ctx, logAction, "object list getted with success", nil)
	return objectList, nil
}

func Remove(ctx context.Context, collection string, query bson.M) error {
	// Now time for metrics
	now := time.Now()

	session := appContext.GetMongoSession()
	defer session.Close()

	logAction := fmt.Sprintf("removing object in collection %s", collection)
	logger.Info(ctx, logAction, "", nil)
	connection := session.DB(mongoDatabaseName).C(collection)

	_, err := connection.RemoveAll(query)

	metrics.MongoDBDurationsSumary.WithLabelValues("Remove").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Remove").Observe(time.Since(now).Seconds())

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Info(ctx, logAction, "object removed with success", nil)
	return nil
}

func Update(ctx context.Context, collection string, objectID bson.ObjectId, objectToUpdate interface{}) error {
	// Now time for metrics
	now := time.Now()

	logAction := fmt.Sprintf("updating object in collection  %s", collection)
	logger.Info(ctx, logAction, "", nil)

	session := appContext.GetMongoSession()
	defer session.Close()

	query := bson.M{"_id": bson.ObjectIdHex(objectID.Hex())}
	change := bson.M{"$set": objectToUpdate}

	connection := session.DB(mongoDatabaseName).C(collection)
	err := connection.Update(query, change)

	metrics.MongoDBDurationsSumary.WithLabelValues("Update").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Update").Observe(time.Since(now).Seconds())

	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	logger.Info(ctx, logAction, "object updated with success", nil)
	return nil
}
