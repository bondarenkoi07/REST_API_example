package Test

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
)

func TestCRDUser(t *testing.T) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}

	dbp, err := database.New(user, password, "localhost", databaseName)

	if err != nil {
		t.Error("Could not create database connection: ", err)
	}

	var TestModel = &Models.User{Login: "foo", Password: "bar"}

	err = dbp.Create("users", TestModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	var rows *pgx.Rows

	rows, err = dbp.ReadAll("users")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	var ValidateModel Models.User

	for (*rows).Next() {
		err = (*rows).Scan(&ValidateModel.Id, &ValidateModel.Login, &ValidateModel.Password)
		if err != nil {
			t.Error("Cannot serialize row: ", ValidateModel)
		}
	}

	if !(ValidateModel.Login == TestModel.Login && ValidateModel.Password == TestModel.Password) {
		t.Error("Test and Validation models are not equal")
	} else {
		var row *pgx.Row
		row, err = dbp.ReadOne("users", ValidateModel.Id)
		if err != nil {
			t.Error("could not read row:", err)
		} else if row == nil {
			t.Error("row is a nil pointer:", ValidateModel)
		} else {

			err = (*row).Scan(&ValidateModel.Id, &ValidateModel.Login, &ValidateModel.Password)

			if err != nil {
				t.Error("Cannot serialize row: ", ValidateModel)
			} else if !(ValidateModel.Login == TestModel.Login && ValidateModel.Password == TestModel.Password) {
				t.Error("Test and Validation models are not equal")
			}
		}

	}

	err = dbp.DeleteOne("users", ValidateModel.Id)

	if err != nil {
		t.Error("could not delete row:", err)
	}
	err = dbp.DeleteAll("users")
	if err != nil {
		t.Error("Could not insert model: ", err)
	}
}

func TestDeveloperOneToOne(t *testing.T) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}

	dbp, err := database.New(user, password, "localhost", databaseName)

	if err != nil {
		t.Error("Could not create database connection: ", err)
	}

	var TestUserModel = &Models.User{Login: "foo", Password: "bar"}

	err = dbp.Create("users", TestUserModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	defer dbp.DeleteAll("users")

	var rows *pgx.Rows

	rows, err = dbp.ReadAll("users")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestUserModel.Id, &TestUserModel.Login, &TestUserModel.Password)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	var TestDeveloperModel = &Models.Developer{OrgName: "roga", Section: "KoPbITA", UserId: TestUserModel}

	err = dbp.Create("developers", TestDeveloperModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}
}

func TestUpdateModel(t *testing.T) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}

	dbp, err := database.New(user, password, "localhost", databaseName)

	if err != nil {
		t.Error("Could not create database connection: ", err)
	}

	var TestUserModel = &Models.User{Login: "foo", Password: "bar"}

	err = dbp.Create("users", TestUserModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	defer dbp.DeleteAll("users")

	var rows *pgx.Rows

	rows, err = dbp.ReadAll("users")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestUserModel.Id, &TestUserModel.Login, &TestUserModel.Password)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	TestUserModel.Password = "foo"

	err = dbp.Update("users", TestUserModel.Id, TestUserModel)

	if err != nil {
		t.Error("update is failed: ", err)
	}
	var ValidateModel = &Models.User{Login: "foo", Password: "foo"}

	var row *pgx.Row
	row, err = dbp.ReadOne("users", TestUserModel.Id)
	if err != nil {
		t.Error("could not read row:", err)
	} else if row == nil {
		t.Error("row is a nil pointer:", ValidateModel)
	} else {

		err = (*row).Scan(&TestUserModel.Id, &TestUserModel.Login, &TestUserModel.Password)

		if err != nil {
			t.Error("Cannot serialize row: ", ValidateModel)
		} else if !(ValidateModel.Login == TestUserModel.Login && ValidateModel.Password == TestUserModel.Password) {
			t.Error("Test and Validation models are not equal")
		}
	}
}

func TestManyToOne(t *testing.T) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}

	dbp, err := database.New(user, password, "localhost", databaseName)

	if err != nil {
		t.Error("Could not create database connection: ", err)
	}

	var TestUserModel = &Models.User{Login: "foo", Password: "bar"}

	err = dbp.Create("users", TestUserModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	var rows *pgx.Rows

	rows, err = dbp.ReadAll("users")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestUserModel.Id, &TestUserModel.Login, &TestUserModel.Password)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	var TestDeveloperModel = &Models.Developer{OrgName: "roga", Section: "KoPbITA", UserId: TestUserModel}

	err = dbp.Create("developers", TestDeveloperModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	var TestMarketModel = &Models.Market{Name: "roga", MaxProducts: 5}

	err = dbp.Create("markets", TestMarketModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	rows, err = dbp.ReadAll("markets")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestMarketModel.Id, &TestMarketModel.Name, &TestMarketModel.MaxProducts)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	for i := 0; i < 5; i++ {
		var TestProductModel = &Models.Product{Name: fmt.Sprintf("name%d", i), Cost: int8(i * 10), Count: int8(i * 5), Developer: TestDeveloperModel, Market: TestMarketModel}

		err = dbp.Create("product", TestProductModel)
	}

	row, err := dbp.GetProductsCount(TestMarketModel.Id)

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	id := 0
	cnt := 0
	err = (*row).Scan(&id, &cnt)
	if err != nil {
		t.Error("Cannot serialize row: ", TestUserModel)
	}

	if cnt != 5 {
		t.Error("Counted not properly enough", TestUserModel)
	}
}
