package Models

type ProductModel struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Cost      int8      `json:"cost"`
	Count     int8      `json:"count"`
	Developer Developer `json:"developer"`
	Market    Market    `json:"market"`
}
