package mysql

import (
	"database/sql"
	"fmt"

	"bitbucket.com/bbd/unzip-parse/entities"
	"bitbucket.com/bbd/unzip-parse/usecases"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	DataSourceName string `json:"dsn"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Database       string `json:"database"`
}

func OpenDB(mysqlUser, mysqlPassword, mysqlPath, dbName string) {
	dbDataSource := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlPath, dbName)
	initDB, err := sql.Open("mysql", dbDataSource)
	usecases.CheckError(err)

	entities.DB = initDB
}
