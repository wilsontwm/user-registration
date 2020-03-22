package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	Mysql    = "mysql"
	Postgres = "postgres"
	Sqlite3  = "sqlite3"
	Mssql    = "mssql"
)

// The configuration for the database connection
type DBConfig struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	Path     string
}

var db *gorm.DB
var dbDriver string
var dbUri string
var tableName = "users"
var isUserActivationRequired = true

// Setup the configuration to the DB connection
func Initialize(dbConfig DBConfig) {
	dbDriver = dbConfig.Driver
	switch dbDriver {
	case Mysql:
		dbUri = fmt.Sprintf("%v:%v@(%v)/%v?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.DBName)
	case Postgres:
		dbUri = fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.DBName, dbConfig.Password)
		// case Sqlite3:
		// 	dbUri = fmt.Sprintf(dbConfig.Path)
		// case Mssql:
		// 	dbUri = fmt.Sprintf("sqlserver://%v:%v@%v:%v?database=%v", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	}

	// Attempt to set up connection
	db := getDB()
	defer db.Close()
}

// Optional configuration to the database / model
func Config(opts ...func()) {
	for _, opt := range opts {
		opt()
	}
}

// Configuration: Set the table name
func TableName(name string) func() {
	return func() {
		tableName = name
	}
}

// Configuration: Set the flag of user activation
func UserActivation(on bool) func() {
	return func() {
		isUserActivationRequired = on
	}
}

// Configuration: Run migration scripts on the database
func MigrateDatabase() func() {
	return migrateDatabase
}

// Run migration scripts on the database
func migrateDatabase() {
	db := getDB()
	defer db.Close()

	db.Debug().CreateTable(&User{})
}

// Get the connection to the DB
func getDB() *gorm.DB {
	db, err := gorm.Open(dbDriver, dbUri)
	if err != nil {
		panic(err)
	}

	return db
}
