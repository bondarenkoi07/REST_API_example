package Models

type ProductModel struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	Cost      int8       `json:"cost"`
	Count     int8       `json:"count"`
	Developer *Developer `json:"developer"`
	Market    *Market    `json:"market"`
}

func (p *ProductModel) GetFields() (string, []interface{}) {
	var ModelValues []interface{}

	dbCols := "name,cost,count,developerId,marketId"

	developerId := p.Developer.UserId.Id
	marketId := p.Market.Id

	ModelValues = append(ModelValues, p.Name)
	ModelValues = append(ModelValues, p.Cost)
	ModelValues = append(ModelValues, p.Count)
	ModelValues = append(ModelValues, developerId)
	ModelValues = append(ModelValues, marketId)

	return dbCols, ModelValues
}
