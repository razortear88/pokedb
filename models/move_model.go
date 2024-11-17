package models

type Move struct {
	Name        string `json:"name", validate:"required"`
	Category    string `json:"category", validate:"required"`
	TypeName    string `json:"typename", validate:"required"`
	Power       int    `json:"power", validate:"required"`
	Accuracy    int    `json:"accuracy", validate:"required"`
	PP          int    `json:"pp", validate:"required"`
	MakeContact bool   `json:"makecontact", validate:"required"`
	Effect      string `json:"effect", validate:"required"`
}
