package driver

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/davidddw2017/panzer/proj/ginMvc/config"
)

// query need rows.Close to release db ins
// exec will release automatic
var MysqlDb *gorm.DB // db pool instance
var MysqlDbErr error // db err instance

func init() {
	// get db config
	dbConfig := config.SystemConfig.DBConfig

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		dbConfig.User,
		dbConfig.Passwd,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
		dbConfig.Charset,
	)

	// connect and open db connection
	MysqlDb, MysqlDbErr = gorm.Open("mysql", dbDSN)

	if MysqlDbErr != nil {
		panic("database data source name error: " + MysqlDbErr.Error())
	}

	// max open connections
	MysqlDb.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)

	// max idle connections
	MysqlDb.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)

	// max lifetime of connection if <=0 will forever
	MysqlDb.DB().SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetimeConns))

	// check db connection at once avoid connect failed
	// else error will be reported until db first sql operate
	if MysqlDbErr = MysqlDb.DB().Ping(); nil != MysqlDbErr {
		panic("database connect failed: " + MysqlDbErr.Error())
	}
}
