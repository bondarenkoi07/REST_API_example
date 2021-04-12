package Service

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"errors"
	"strconv"
)

type DeveloperService struct {
	dbp *database.Database
}

func NewDeveloperService(dbp *database.Database) *DeveloperService {
	return &DeveloperService{dbp: dbp}
}

func (ds DeveloperService) Create(model Models.Developer) error {
	err := ds.dbp.Create("developers", model)
	return err
}

func (ds DeveloperService) ReadOne(id int64) (*Models.Developer, error) {
	row, err := ds.dbp.ReadOne("developers", id)

	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, nil
	}

	var model Models.Developer
	err = (*row).Scan(&id, &model.OrgName, &model.Section)

	if err != nil {
		return nil, err
	}

	userService := NewUserService(ds.dbp)

	user, err := userService.ReadOne(id)
	if err != nil {
		return nil, err
	}

	model.UserId = user

	return &model, nil
}

func (ds DeveloperService) ReadAll() ([]Models.Developer, error) {
	rows, err := ds.dbp.ReadAll("developers")

	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	defer (*rows).Close()

	models := make([]Models.Developer, 0)

	for (*rows).Next() {
		var iterModel Models.Developer
		var id int64
		err = (*rows).Scan(&id, &iterModel.OrgName, &iterModel.Section)
		if err != nil {
			return nil, err
		}

		userService := NewUserService(ds.dbp)

		user, err := userService.ReadOne(id)
		if err != nil {
			return nil, err
		}

		iterModel.UserId = user

		models = append(models, iterModel)
	}
	return models, nil
}

func (ds DeveloperService) Update(model Models.Developer, id int64) error {
	err := ds.dbp.Update("developers", id, model)
	return err
}

func (ds DeveloperService) DeleteOne(id int64) error {
	err := ds.dbp.DeleteOne("developers", id)
	return err
}

func (ds DeveloperService) DeleteAll() error {
	return ds.dbp.DeleteAll("developers")
}

func (ds *DeveloperService) Deserialize(data map[string]string, Service UserService) (error, Models.Developer) {
	var validate = true
	id, isSet := data["id"]
	validate = validate && isSet

	Id, err := strconv.Atoi(id)
	if err != nil {
		return err, Models.Developer{}
	}

	orgName, isSet := data["org_name"]
	validate = validate && isSet

	section, isSet := data["section"]
	validate = validate && isSet

	if validate {
		var developer Models.Developer
		developer.OrgName = orgName
		developer.Section = section

		developer.UserId, err = Service.ReadOne(int64(Id))
		if err != nil {
			return err, Models.Developer{}
		}

		return nil, developer
	} else {
		return errors.New("wrong JSON"), Models.Developer{}
	}
}
