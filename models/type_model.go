package models

type Type struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
}
