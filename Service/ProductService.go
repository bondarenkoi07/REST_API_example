package Service

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
)

type ProductService struct {
	dbp *database.Database
}

func NewProductService(dbp *database.Database) *ProductService {
	return &ProductService{dbp: dbp}
}

func (ps ProductService) Create(model Models.Product) error {
	MarketCapacity := model.Market.MaxProducts
	count, err := ps.getProductsCount(model.Market.Id)
	if count == -1 && err != nil {
		if !(err.Error() == "no rows in result set") {
			return err
		}

	} else if err == nil && count == -1 {
		return errors.New("err occurred while getting products marketId")
	}

	if MarketCapacity >= count+int64(model.Count) {

		err = ps.dbp.Create("product", model)
	} else {
		err = errors.New("market's storage is full")
	}

	return err
}

func (ps ProductService) ReadOne(id int64) (*Models.Product, error) {
	row, err := ps.dbp.ReadOne("product", id)

	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, nil
	}

	var DeveloperId int64
	var MarketId int64

	var model Models.Product
	err = (*row).Scan(&model.Id, &model.Name, &model.Cost, &model.Count, &DeveloperId, &MarketId)

	if err != nil {
		return nil, err
	}

	developerService := NewDeveloperService(ps.dbp)

	marketService := NewMarketService(ps.dbp)

	developer, err := developerService.ReadOne(DeveloperId)

	if err != nil {
		return nil, err
	}

	market, err := marketService.ReadOne(MarketId)

	model.Developer = developer
	model.Market = market

	return &model, nil
}

func (ps ProductService) ReadAll() ([]Models.Product, error) {
	rows, err := ps.dbp.ReadAll("product")

	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}
	return ps.fetchProducts(rows)
}

func (ps ProductService) Update(model Models.Product, id int64) error {
	MarketCapacity := model.Market.MaxProducts
	count, err := ps.getProductsCount(model.Market.Id)
	if MarketCapacity > count+int64(model.Count) {
		err = ps.dbp.Update("product", id, model)
	} else {
		err = errors.New(fmt.Sprintf(" storage is full in %s market", model.Market.Name))
	}

	return err
}

func (ps ProductService) DeleteOne(id int64) error {
	err := ps.dbp.DeleteOne("product", id)
	return err
}

func (ps ProductService) DeleteAll() error {
	return ps.dbp.DeleteAll("product")
}

func (ps ProductService) getProductsCount(id int64) (int64, error) {
	row, err := ps.dbp.GetProductsCount(id)
	if err != nil {
		return -1, err
	}

	var (
		Id    int64
		count int64
	)

	err = (*row).Scan(&Id, &count)

	if err != nil {
		return -1, err
	} else if Id != id {
		return -1, errors.New("could not find current market")
	} else {
		return count, nil
	}
}

func (ps *ProductService) Deserialize(data map[string]string, devService DeveloperService, marketService MarketService) (error, Models.Product) {
	var validate = true

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

func (ps ProductService) FilterProductsByMarket(id int64) ([]Models.Product, error) {
	rows, err := ps.dbp.GetMarketProducts(id)
	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	return ps.fetchProducts(rows)
}

func (ps ProductService) FilterProductsByDeveloper(id int64) ([]Models.Product, error) {
	rows, err := ps.dbp.GetDeveloperProducts(id)
	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	return ps.fetchProducts(rows)
}

func (ps ProductService) fetchProducts(rows *pgx.Rows) ([]Models.Product, error) {

	defer (*rows).Close()

	models := make([]Models.Product, 0)

	for (*rows).Next() {
		var iterModel Models.Product
		var DeveloperId int64
		var MarketId int64
		err := (*rows).Scan(&iterModel.Id, &iterModel.Name, &iterModel.Cost, &iterModel.Count, &DeveloperId, &MarketId)
		if err != nil {
			return nil, err
		}

		developerService := NewDeveloperService(ps.dbp)

		marketService := NewMarketService(ps.dbp)

		developer, err := developerService.ReadOne(DeveloperId)

		if err != nil {
			return nil, err
		}

		market, err := marketService.ReadOne(MarketId)

		iterModel.Developer = developer
		iterModel.Market = market

		models = append(models, iterModel)
	}
	return models, nil
}
