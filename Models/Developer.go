package Models

type Developer struct {
	Id      int64  `json:"id"`
	OrgName string `json:"org_name"`
	Section string `json:"section"`
	UserId  User   `json:"user"`
}
