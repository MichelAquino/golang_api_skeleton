package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/michelaquino/golang_api_skeleton/src/log"
	"github.com/michelaquino/golang_api_skeleton/src/metrics"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongoClient  *mongo.Client
	onceDatabase sync.Once
	logger       = log.GetLogger()
)

// GetMongoClient return a copy of mongodb session
func GetMongoClient(context context.Context) (*mongo.Client, error) {
	onceDatabase.Do(func() {
		var err error
		mongoClient, err = newMongoClient(context)
		if err != nil {
			logger.Error(context, "Getting new mongo connection", err.Error(), nil)
		}
	})

	logger.Debug(context, "Getting mongo connection", "Success", nil)
	return mongoClient, nil
}

// NewMongoClient return a new mongo client
func newMongoClient(ctx context.Context) (*mongo.Client, error) {
	logger.Debug(ctx, "Getting new mongo connection", "", nil)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		logger.Error(ctx, "Connect", err.Error(), nil)
		return nil, err
	}

	contextTimeout, _ := context.WithTimeout(ctx, 5*time.Second)
	if err := client.Ping(contextTimeout, readpref.Primary()); err != nil {
		logger.Error(ctx, "MongoDB Ping", err.Error(), nil)
		panic(err)
	}

	logger.Debug(ctx, "Getting new mongo connection", "Success", nil)
	return client, nil
}

// Insert a new object on database
func Insert(ctx context.Context, collectionName string, objectToInsert interface{}) error {
	// Now time for metrics
	now := time.Now()

	mongoClient, err := GetMongoClient(ctx)
	if err != nil {
		logger.Error(ctx, "Getting mongo client", err.Error(), nil)
		return err
	}

	logAction := fmt.Sprintf("Inserting object in collection %s", collectionName)
	logger.Info(ctx, logAction, "", nil)

	collection := mongoClient.Database(viper.GetString("mongo.database.name")).Collection(collectionName)
	if _, err := collection.InsertOne(ctx, objectToInsert); err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	// Send metrics to Prometheus
	metrics.MongoDBDurationsSumary.WithLabelValues("Insert").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Insert").Observe(time.Since(now).Seconds())

	logger.Info(ctx, logAction, "object inserted with success", nil)
	return nil
}

func FindOne(ctx context.Context, collectionName string, query bson.M, object interface{}) error {
	// Now time for metrics
	now := time.Now()

	mongoClient, err := GetMongoClient(ctx)
	if err != nil {
		logger.Error(ctx, "Getting mongo client", err.Error(), nil)
		return err
	}

	logAction := fmt.Sprintf("getting object in collection %s", collectionName)
	logger.Info(ctx, logAction, "", nil)
	collection := mongoClient.Database(viper.GetString("mongo.database.name")).Collection(collectionName)

	result := collection.FindOne(ctx, query)
	if err := result.Decode(&object); err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	metrics.MongoDBDurationsSumary.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())

	logger.Info(ctx, logAction, "object got with success", nil)
	return nil
}

func FindAll(ctx context.Context, collectionName string, query bson.M) ([]interface{}, error) {
	// Now time for metrics
	now := time.Now()

	mongoClient, err := GetMongoClient(ctx)
	if err != nil {
		logger.Error(ctx, "Getting mongo client", err.Error(), nil)
		return nil, err
	}

	logAction := fmt.Sprintf("getting object list in collection %s", collectionName)
	logger.Info(ctx, logAction, "", nil)

	collection := mongoClient.Database(viper.GetString("mongo.database.name")).Collection(collectionName)

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return nil, err
	}

	var objectList []interface{}
	for cursor.Next(ctx) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			logger.Error(ctx, logAction, err.Error(), nil)
			return nil, err
		}

		objectList = append(objectList, result)
	}

	if err := cursor.Err(); err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return nil, err
	}

	metrics.MongoDBDurationsSumary.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())

	logger.Info(ctx, logAction, "object list getted with success", nil)
	return objectList, nil
}

func Remove(ctx context.Context, collectionName string, query bson.M) error {
	// Now time for metrics
	now := time.Now()

	mongoClient, err := GetMongoClient(ctx)
	if err != nil {
		logger.Error(ctx, "Getting mongo client", err.Error(), nil)
		return err
	}

	logAction := fmt.Sprintf("removing object in collection %s", collectionName)
	logger.Info(ctx, logAction, "", nil)

	collection := mongoClient.Database(viper.GetString("mongo.database.name")).Collection(collectionName)
	if _, err := collection.DeleteOne(ctx, query); err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	metrics.MongoDBDurationsSumary.WithLabelValues("Remove").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Remove").Observe(time.Since(now).Seconds())

	logger.Info(ctx, logAction, "object removed with success", nil)
	return nil
}

func Update(ctx context.Context, collectionName string, objectID bson.ObjectId, objectToUpdate interface{}) error {
	// Now time for metrics
	now := time.Now()

	logAction := fmt.Sprintf("updating object in collection  %s", collectionName)
	logger.Info(ctx, logAction, "", nil)

	mongoClient, err := GetMongoClient(ctx)
	if err != nil {
		logger.Error(ctx, "Getting mongo client", err.Error(), nil)
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(objectID.Hex())}
	change := bson.M{"$set": objectToUpdate}

	collection := mongoClient.Database(viper.GetString("mongo.database.name")).Collection(collectionName)
	if _, err := collection.UpdateOne(ctx, query, change); err != nil {
		logger.Error(ctx, logAction, err.Error(), nil)
		return err
	}

	metrics.MongoDBDurationsSumary.WithLabelValues("Update").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Update").Observe(time.Since(now).Seconds())

	logger.Info(ctx, logAction, "object updated with success", nil)
	return nil
}
