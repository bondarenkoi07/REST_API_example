package Service

import (
	"app/REST_API_example/database"
	"os"
)

type Service struct {
	DeveloperService *DeveloperService
	UserService      *UserService
	ProductService   *ProductService
	MarketService    *MarketService
}

func NewService() Service {

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	connection := "db"
	print("test: ", user, " ", password, " ", databaseName)
	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}

	dbp, err := database.New(user, password, connection, databaseName)
	if err != nil {
		connection = "localhost"
		dbp, err = database.New(user, password, connection, databaseName)
		if err != nil {
			panic(err)
		}
	}

	var service Service
	service.ProductService = NewProductService(dbp)
	service.DeveloperService = NewDeveloperService(dbp)
	service.MarketService = NewMarketService(dbp)
	service.UserService = NewUserService(dbp)
	return service
}
