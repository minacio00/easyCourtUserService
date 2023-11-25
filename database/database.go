package database

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbCredentials struct {
	client   string
	user     string
	password string
	port     string
	host     string
	database string
	ssl      string
}

func (db *dbCredentials) formatStr() string {
	return "user=" + db.user + " host=" + db.host + " port=" + db.port + " password=" + db.password + " dbname=" + db.database + " sslmode=" + db.ssl
}

var Db *gorm.DB

func Connectdb() {
	sql.Drivers()
	creds := dbCredentials{
		client: "postgresql", user: "postgres",
		password: "postgresql", port: "5432",
		host:     "db",
		database: "clubster",
		ssl:      "disable",
	}
	var err error
	env := os.Getenv("APP_ENV")
	if env == "test" {
		Db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("connected to sqlite")
	} else {
		Db, err = gorm.Open(postgres.Open(creds.formatStr()), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	}

	if err != nil {
		log.Fatal(err.Error())
	}

}
