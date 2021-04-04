package Models

type Market struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	MaxProducts int64  `json:"max_products"`
}

func (m *Market) GetFields() (string, []interface{}) {
	var ModelValues []interface{}
	dbCols := "name,max_products"

	ModelValues = append(ModelValues, m.Name)
	ModelValues = append(ModelValues, m.MaxProducts)

	return dbCols, ModelValues
}
