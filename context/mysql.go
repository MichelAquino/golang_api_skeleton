package context

import (
	"fmt"
	"sync"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var onceMysqlDatabase sync.Once
var mysqlSession *sql.DB

func GetMysqlConnection() *sql.DB {
	apiConfig := GetAPIConfig()

	onceMysqlDatabase.Do(func() {
		var err error

		mysqlCredentials := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			apiConfig.MySQLConfig.Username,
			apiConfig.MySQLConfig.Password,
			apiConfig.MySQLConfig.Host,
			apiConfig.MySQLConfig.Port,
			apiConfig.MySQLConfig.Database)

		mysqlSession, err = sql.Open("mysql", mysqlCredentials)
		if err != nil {
			errorMsg := fmt.Sprintf("Error on start database: %s", err.Error())
			panic(errorMsg)
		}

		mysqlSession.SetMaxOpenConns(10)
		mysqlSession.SetMaxIdleConns(20)
		mysqlSession.SetConnMaxLifetime(time.Second * 60)
	})

	return mysqlSession
}

func CheckMysqlConn() error {
	conn := GetMysqlConnection()
	err := conn.Ping()
	return err
}
