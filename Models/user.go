package Models

type User struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) GetFields() (string, []interface{}) {
	var ModelValues []interface{}
	dbCols := "login,password"

	ModelValues = append(ModelValues, u.Login)
	ModelValues = append(ModelValues, u.Password)

	return dbCols, ModelValues
}
