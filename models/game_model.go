package models

type Game struct {
	Name       string `json:"name" validate:"required"`
	Generation int    `json:"generation" validate:"required"`
	Cover      string `json:"cover" validate:"required"`
}
