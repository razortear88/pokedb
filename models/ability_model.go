package models

type Ability struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
