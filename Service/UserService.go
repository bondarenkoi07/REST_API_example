package Service

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/database"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
)

type UserService struct {
	dbp *database.Database
}

func NewUserService(dbp *database.Database) *UserService {
	return &UserService{dbp: dbp}
}

func (us UserService) Create(model Models.User) error {
	err := us.dbp.Create("users", model)
	return err
}

func (us UserService) ReadOne(id int64) (*Models.User, error) {
	row, err := us.dbp.ReadOne("users", id)

	if err != nil {
		return nil, err
	} else if row == nil {
		return nil, nil
	}

	var model Models.User
	err = (*row).Scan(&model.Id, &model.Login, &model.Password)
	if err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

func (us UserService) ReadAll() ([]Models.User, error) {
	rows, err := us.dbp.ReadAll("users")

	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, nil
	}

	defer (*rows).Close()

	models := make([]Models.User, 0)

	for (*rows).Next() {
		var iterModel Models.User
		err = (*rows).Scan(&iterModel.Id, &iterModel.Login, &iterModel.Password)
		if err != nil {
			return nil, err
		} else {
			models = append(models, iterModel)
		}
	}
	return models, nil
}

func (us UserService) Update(model Models.User, id int64) error {
	err := us.dbp.Update("users", id, model)
	return err
}

func (us UserService) DeleteOne(id int64) error {
	err := us.dbp.DeleteOne("users", id)
	return err
}

func (us UserService) DeleteAll() error {
	return us.dbp.DeleteAll("users")
}

func (us UserService) Deserialize(data map[string]string) (error, Models.User) {
	var validate bool
	id, isSet := data["id"]
	validate = validate && isSet

	Id, err := strconv.Atoi(id)
	if err != nil {
		return err, Models.User{}
	}

	Name, isSet := data["login"]
	validate = validate && isSet

	password, isSet := data["password"]
	validate = validate && isSet

	if validate {
		var user Models.User
		user.Login = Name
		hashPass := sha256.Sum256([]byte(password))
		user.Password = hex.EncodeToString(hashPass[:])

		user.Id = int64(Id)

		return nil, user
	} else {
		return errors.New("wrong JSON"), Models.User{}
	}
}
