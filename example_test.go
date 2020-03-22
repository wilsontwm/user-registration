package user

import (
	"reflect"
	"testing"
)

func contains(array interface{}, ele interface{}) bool {
	arr := reflect.ValueOf(array)

	if arr.Kind() != reflect.Slice {
		panic("Invalid slice type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == ele {
			return true
		}
	}

	return false
}

// Test connection with MySQL database with configuration
func TestConnectMysqlDatabase(t *testing.T) {
	tableName := "tests"
	// Mysql
	dbConfig := DBConfig{
		Driver:   Mysql,
		Username: "admin",
		Password: "password",
		Host:     "localhost",
		DBName:   "test",
	}

	Initialize(dbConfig)
	Config(TableName(tableName), MigrateDatabase())

	var ret []string
	db := getDB()
	err := db.Raw("show tables").Pluck("Tables_in_mysql", &ret).Error

	if err != nil {
		t.Error(err)
	}

	if ok := contains(ret, tableName); !ok {
		t.Errorf("The table '%v' is not created.", tableName)
	}
}

// Test connection with Postgres database with configuration
func TestConnectPostgresDatabase(t *testing.T) {
	tableName := "tests"
	// Postgres
	dbConfig := DBConfig{
		Driver:   Postgres,
		Username: "postgres",
		Password: "password",
		Host:     "localhost",
		Port:     "5432",
		DBName:   "test",
	}

	Initialize(dbConfig)
	Config(TableName(tableName), MigrateDatabase())

	var ret []string
	db := getDB()
	err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = 'public' ORDER BY table_type, table_name").Pluck("table_name", &ret).Error

	if err != nil {
		t.Error(err)
	}

	if ok := contains(ret, tableName); !ok {
		t.Errorf("The table '%v' is not created.", tableName)
	}
}
