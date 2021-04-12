package Service

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"errors"
	"strconv"
)

type MarketService struct {
	dbp *database.Database
}

func NewMarketService(dbp *database.Database) *MarketService {
	return &MarketService{dbp: dbp}
}

func (ms MarketService) Create(model Models.Market) error {
	err := ms.dbp.Create("markets", model)
	return err
}

func (ms MarketService) ReadOne(id int64) (*Models.Market, error) {
	row, err := ms.dbp.ReadOne("markets", id)

	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, nil
	}

	var model Models.Market
	err = (*row).Scan(&model.Id, &model.Name, &model.MaxProducts)
	if err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

func (ms MarketService) ReadAll() ([]Models.Market, error) {
	rows, err := ms.dbp.ReadAll("markets")

	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	defer (*rows).Close()

	models := make([]Models.Market, 0)

	for (*rows).Next() {
		var iterModel Models.Market
		err = (*rows).Scan(&iterModel.Id, &iterModel.Name, &iterModel.MaxProducts)
		if err != nil {
			return nil, err
		} else {
			models = append(models, iterModel)
		}
	}
	return models, nil
}

func (ms MarketService) Update(model Models.Market, id int64) error {
	err := ms.dbp.Update("markets", id, model)
	return err
}

func (ms MarketService) DeleteOne(id int64) error {
	err := ms.dbp.DeleteOne("markets", id)
	return err
}

func (ms MarketService) DeleteAll() error {
	return ms.dbp.DeleteAll("markets")
}

func (ms MarketService) Deserialize(data map[string]string) (Models.Market, error) {
	var validate = true

	Name, isSet := data["name"]
	validate = validate && isSet

	max, isSet := data["max_products"]
	validate = validate && isSet

	Max, err := strconv.Atoi(max)
	if err != nil {
		return Models.Market{}, err
	}

	if validate {
		var market Models.Market
		market.Name = Name
		market.MaxProducts = int64(Max)

		return market, nil
	} else {
		return Models.Market{}, errors.New("wrong JSON")
	}
}
