package context

import (
	"fmt"
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var mongoSession *mgo.Session
var onceDatabase sync.Once

// GetMongoSession returns a copy of MongoDB session.
func GetMongoSession() *mgo.Session {
	onceDatabase.Do(func() {
		var err error

		mongoSession, err = getNewMongoSession()
		if err != nil {
			errorMsg := fmt.Sprintf("Error on start database: %s", err.Error())
			panic(errorMsg)
		}
	})

	return mongoSession.Copy()
}

func getNewMongoSession() (*mgo.Session, error) {
	apiConfig := GetAPIConfig()
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs: []string{
			apiConfig.MongoDBConfig.Address,
		},
		Database: apiConfig.MongoDBConfig.DatabaseName,
		Timeout:  apiConfig.MongoDBConfig.Timeout,
		Username: apiConfig.MongoDBConfig.Username,
		Password: apiConfig.MongoDBConfig.Password,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	return session, err
}
