package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	databaseName   = "tokens"
	collectionName = "blacklistedTokens"
)

type Tenant struct {
	ID          uint   `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	TrialPeriod bool   `json:"periodo_teste"`
}

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
var MongoClient *mongo.Client

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
		t := &Tenant{}
		Db.AutoMigrate(&t)
	} else {
		Db, err = gorm.Open(postgres.Open(creds.formatStr()), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	connectMongoDb()
}

func connectMongoDb() {
	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoUser == "" || mongoPassword == "" {
		log.Fatal("missing MONGO_INITDB_ROOT_USERNAME or MONGO_INITDB_ROOT_PASSWORD")
	}

	mongoUri := fmt.Sprintf("mongodb://%s:%s@mongo:27017/?authSource=admin", mongoUser, mongoPassword)
	clientOptions := options.Client().ApplyURI(mongoUri)

	var err error
	MongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = MongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	createTokenTTLIndex()
}

func createTokenTTLIndex() {
	collection := MongoClient.Database(databaseName).Collection(collectionName)

	// Specify the index model
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expirationTime", Value: 1}},    // Index on the expirationTime field
		Options: options.Index().SetExpireAfterSeconds(86400), // 1 day
	}

	// Create the TTL index
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
}
