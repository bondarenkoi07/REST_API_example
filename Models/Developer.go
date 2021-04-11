package Models

type Developer struct {
	OrgName string `json:"org_name"`
	Section string `json:"section"`
	UserId  *User  `json:"id"`
}

func (d Developer) GetFields() (string, []interface{}) {
	var ModelValues []interface{}
	dbCols := "id,org_name,section"
	id := d.UserId.Id
	ModelValues = append(ModelValues, id)
	ModelValues = append(ModelValues, d.OrgName)
	ModelValues = append(ModelValues, d.Section)

	return dbCols, ModelValues
}
