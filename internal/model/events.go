package model

type Event struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	User  string `json:"user"`
	Team  string `json:"team"`
	Role  string `json:"role"`
}
