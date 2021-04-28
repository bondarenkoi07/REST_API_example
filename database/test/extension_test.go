package test

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"os"
	"testing"
)

var (
	dbp *database.Database

	testModels = []struct {
		model database.IterableModel
		table string
	}{
		{Models.User{Login: "foo", Password: "bar"}, "users"},
		{Models.User{Login: "foo2", Password: "bar2"}, "users"},
	}
)

func init() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	connection := "db"
	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}
	var err error
	dbp, err = database.New(user, password, connection, databaseName)
	if err != nil {
		connection = "localhost"
		dbp, err = database.New(user, password, connection, databaseName)
		if err != nil {
			panic(err)
		}
	}
	dbp.DeleteAll("users")
	dbp.DeleteAll("markets")
}

func TestSave(t *testing.T) {
	var IDs = make([]int64, 0)
	for _, value := range testModels {
		id, err := dbp.Save(value.table, value.model)
		if err != nil {
			t.Error(err)
		}
		IDs = append(IDs, id)
	}
	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "users")
	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "markets")
	for i := 0; i < len(IDs)-1; i++ {
		if !(IDs[i+1]-IDs[i] == 1) {
			t.Error("Id do not make up sequence")
		}
	}
}
