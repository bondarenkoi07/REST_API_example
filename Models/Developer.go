package Models

import "fmt"

type Developer struct {
	OrgName string `json:"org_name"`
	Section string `json:"section"`
	UserId  *User  `json:"user"`
}

func (d *Developer) GetFields() (string, []interface{}) {
	var ModelValues []interface{}
	dbCols := "id,org_name,section"
	id := d.UserId.Id
	fmt.Println("inserted user id", id)
	ModelValues = append(ModelValues, id)
	ModelValues = append(ModelValues, d.OrgName)
	ModelValues = append(ModelValues, d.Section)

	return dbCols, ModelValues
}
