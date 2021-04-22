package test

import (
	"app/REST_API_example/Models"
	_ "app/REST_API_example/Models"
	"app/REST_API_example/Service"
	"app/REST_API_example/database"
	"os"
	"testing"
)

var (
	dbp *database.Database

	user      = Models.User{Login: "foo14", Password: "bar14"}
	developer = Models.Developer{OrgName: "foo23", Section: "bar12"}
	market1   = Models.Market{
		Name:        "techs2",
		MaxProducts: 10,
	}
	market2 = Models.Market{
		Name:        "technologies3",
		MaxProducts: 10,
	}
	product1 = Models.Product{
		Name:  "phone",
		Cost:  10,
		Count: 6,
	}
	product2 = Models.Product{
		Name:  "tablet",
		Cost:  10,
		Count: 6,
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

}

func TestGetMarketProducts(t *testing.T) {
	UserService := Service.NewUserService(dbp)
	MarketService := Service.NewMarketService(dbp)
	ProductService := Service.NewProductService(dbp)
	DeveloperService := Service.NewDeveloperService(dbp)

	err := UserService.DeleteAll()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	err = MarketService.DeleteAll()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	id, err := UserService.Save(user)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	user.Id = id

	developer.UserId = &user

	err = DeveloperService.Create(developer)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	id, err = MarketService.Save(market1)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	market1.Id = id

	id, err = MarketService.Save(market2)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	market2.Id = id

	product1.Market = &market1
	product1.Developer = &developer

	err = ProductService.Create(product1)

	if err != nil {
		t.Error(err)
	}

	product2.Market = &market2
	product2.Developer = &developer

	err = ProductService.Create(product2)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	validateProducts, err := ProductService.FilterProductsByMarket(market1.Id)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else if len(validateProducts) != 1 {
		t.Errorf("Wrong number of products: should be 1, got %d", len(validateProducts))
	} else {
		if !(validateProducts[0].Count == product1.Count && validateProducts[0].Cost == product1.Cost && validateProducts[0].Name == product1.Name) {
			t.Error("Wrong Model fetched")
		}
	}

	validateProducts, err = ProductService.FilterProductsByDeveloper(developer.UserId.Id)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else if len(validateProducts) != 2 {
		t.Errorf("Wrong number of products: should be 2, got %d", len(validateProducts))
	} else {
		validation := validateProducts[0].Count == product1.Count
		validation = validation && validateProducts[0].Cost == product1.Cost
		validation = validation && validateProducts[0].Name == product1.Name
		validation = validation && validateProducts[1].Cost == product2.Cost
		validation = validation && validateProducts[1].Name == product2.Name
		validation = validation && validateProducts[1].Count == product1.Count
		if !(validation) {
			t.Error("Wrong Model fetched")
		}
	}

}
