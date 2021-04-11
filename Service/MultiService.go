package Service

import "app/REST_API_example/database"

type Service struct {
	DeveloperService *DeveloperService
	UserService      *UserService
	ProductService   *ProductService
	MarketService    *MarketService
}

func NewService(dbp *database.Database) Service {
	var service Service
	service.ProductService = NewProductService(dbp)
	service.DeveloperService = NewDeveloperService(dbp)
	service.MarketService = NewMarketService(dbp)
	service.UserService = NewUserService(dbp)
	return service
}
