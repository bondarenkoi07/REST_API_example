package Service

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"errors"
	"strconv"
)

type ProductService struct {
	dbp *database.Database
}

func NewProductService(dbp *database.Database) *ProductService {
	return &ProductService{dbp: dbp}
}

func (ps ProductService) Create(model Models.Product) error {
	err := ps.dbp.Create("products", model)
	return err
}

func (ps ProductService) ReadOne(id int64) (*Models.Product, error) {
	row, err := ps.dbp.ReadOne("products", id)

	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, nil
	}

	var model Models.Product
	err = (*row).Scan(&model.Name, &model.Cost, &model.Count, &model.Developer.UserId.Id, &model.Market.Id)
	if err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

func (ps ProductService) ReadAll() ([]Models.Product, error) {
	rows, err := ps.dbp.ReadAll("products")

	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	defer (*rows).Close()

	models := make([]Models.Product, 0)

	for (*rows).Next() {
		var iterModel Models.Product
		err = (*rows).Scan(&iterModel.Name, &iterModel.Cost, &iterModel.Count, &iterModel.Developer.UserId.Id, &iterModel.Market.Id)
		if err != nil {
			return nil, err
		} else {
			models = append(models, iterModel)
		}
	}
	return models, nil
}

func (ps ProductService) Update(model Models.Product, id int64) error {
	err := ps.dbp.Update("products", id, model)
	return err
}

func (ps ProductService) DeleteOne(id int64) error {
	err := ps.dbp.DeleteOne("products", id)
	return err
}

func (ps ProductService) DeleteAll() error {
	return ps.dbp.DeleteAll("products")
}

func (ps *ProductService) Deserialize(data map[string]string, devService DeveloperService, marketService MarketService) (error, Models.Product) {
	var validate bool

	name, isSet := data["name"]
	validate = validate && isSet

	cost, isSet := data["cost"]
	validate = validate && isSet

	Cost, err := strconv.Atoi(cost)
	if err != nil {
		return err, Models.Product{}
	}

	count, isSet := data["count"]
	validate = validate && isSet

	Count, err := strconv.Atoi(count)
	if err != nil {
		return err, Models.Product{}
	}

	developerId, isSet := data["developer"]
	validate = validate && isSet

	DeveloperId, err := strconv.Atoi(developerId)
	if err != nil {
		return err, Models.Product{}
	}

	marketId, isSet := data["market"]
	validate = validate && isSet

	MarketId, err := strconv.Atoi(marketId)
	if err != nil {
		return err, Models.Product{}
	}

	if validate {
		var product Models.Product

		product.Name = name
		product.Cost = int8(Cost)
		product.Count = int8(Count)

		product.Developer, err = devService.ReadOne(int64(DeveloperId))
		if err != nil {
			return err, Models.Product{}
		}

		product.Market, err = marketService.ReadOne(int64(MarketId))
		if err != nil {
			return err, Models.Product{}
		}

		return nil, product
	} else {
		return errors.New("wrong JSON"), Models.Product{}
	}
}
