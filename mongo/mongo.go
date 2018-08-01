package mongo

import (
	"fmt"
	"time"

	"github.com/michelaquino/golang_api_skeleton/context"
	"github.com/michelaquino/golang_api_skeleton/metrics"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoDatabaseName = "api-skeleton"
)

func Insert(collection string, objectToInsert interface{}) error {
	// Now time for metrics
	now := time.Now()

	session := context.GetMongoSession()
	defer session.Close()

	log := context.GetLogger()
	log.Info("Mongo", "Insert", "", "", "", fmt.Sprintf("Inserting object in collection %s", collection), "")

	connection := session.DB(mongoDatabaseName).C(collection)
	err := connection.Insert(&objectToInsert)

	// Send metrics to Prometheus
	metrics.MongoDBDurationsSumary.WithLabelValues("Insert").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Insert").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("Mongo", "Insert", "", "", "", "Error on insert object", err.Error())
		return err
	}

	log.Info("Mongo", "Insert", "", "", "", "Object inserted with success", "")
	return nil
}

func FindOne(collection string, query bson.M, object interface{}) error {
	// Now time for metrics
	now := time.Now()

	session := context.GetMongoSession()
	defer session.Close()

	log := context.GetLogger()

	log.Info("Mongo", "Find", "", "", "", fmt.Sprintf("Getting object in collection %s", collection), "")
	connection := session.DB(mongoDatabaseName).C(collection)

	err := connection.Find(query).One(object)

	metrics.MongoDBDurationsSumary.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindOne").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("Mongo", "Find", "", "", "", "Error on getting object", err.Error())
		return err
	}

	log.Info("Mongo", "Find", "", "", "", "Object getted with success", "")
	return nil
}

func FindAll(collection string, query bson.M) ([]interface{}, error) {
	// Now time for metrics
	now := time.Now()

	session := context.GetMongoSession()
	defer session.Close()

	log := context.GetLogger()
	var objectList []interface{}

	log.Info("Mongo", "FindAll", "", "", "", fmt.Sprintf("Getting object list in collection %s", collection), "")
	connection := session.DB(mongoDatabaseName).C(collection)

	err := connection.Find(query).All(&objectList)

	metrics.MongoDBDurationsSumary.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("FindAll").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("Mongo", "FindAll", "", "", "", "Error on getting object list", err.Error())
		return nil, err
	}

	log.Info("Mongo", "FindAll", "", "", "", "Object list getted with success", "")
	return objectList, nil
}

func Remove(collection string, query bson.M) error {
	// Now time for metrics
	now := time.Now()

	session := context.GetMongoSession()
	defer session.Close()

	log := context.GetLogger()

	log.Info("Mongo", "Remove", "", "", "", fmt.Sprintf("Removing object in collection %s", collection), "")
	connection := session.DB(mongoDatabaseName).C(collection)

	_, err := connection.RemoveAll(query)

	metrics.MongoDBDurationsSumary.WithLabelValues("Remove").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Remove").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("Mongo", "Remove", "", "", "", "Error on remove object", err.Error())
		return err
	}

	log.Info("Mongo", "Remove", "", "", "", "Object removed with success", "")
	return nil
}

func Update(collection string, objectID bson.ObjectId, objectToUpdate interface{}) error {
	// Now time for metrics
	now := time.Now()

	log := context.GetLogger()
	log.Info("Mongo", "Update", "", "", "", fmt.Sprintf("Updating object in collection %s", collection), "")

	session := context.GetMongoSession()
	defer session.Close()

	query := bson.M{"_id": bson.ObjectIdHex(objectID.Hex())}
	change := bson.M{"$set": objectToUpdate}

	connection := session.DB(mongoDatabaseName).C(collection)
	err := connection.Update(query, change)

	metrics.MongoDBDurationsSumary.WithLabelValues("Update").Observe(time.Since(now).Seconds())
	metrics.MongoDBDurationsHistogram.WithLabelValues("Update").Observe(time.Since(now).Seconds())

	if err != nil {
		log.Error("Mongo", "Update", "", "", "", "Error on update object", err.Error())
		return err
	}

	log.Info("Mongo", "Update", "", "", "", "Object updated with success", "")
	return nil
}
