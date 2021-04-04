package Models

type Market struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	MaxProducts int64  `json:"max_products"`
}
