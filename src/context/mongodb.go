package context

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
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
	mongoURL := viper.GetString("mongo.url")
	mongoPort := viper.GetInt("mongo.port")
	mongoAddress := fmt.Sprintf("%s:%d", mongoURL, mongoPort)

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs: []string{
			mongoAddress,
		},
		Database: viper.GetString("mongo.database.name"),
		Timeout:  viper.GetDuration("mongo.timeout"),
		Username: viper.GetString("mongo.username"),
		Password: viper.GetString("mongo.password"),
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	return session, err
}
