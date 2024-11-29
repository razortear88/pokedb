package models

type Pokemon struct {
	Name       string   `json:"name" validate:"required"`
	NationalNo string   `json:"nationalno" validate:"required"`
	Type       []string `json:"type" validate:"required"`
	Species    string   `json:"species" validate:"required"`
	Height     float32  `json:"height" validate:"required"`
	Weight     float32  `json:"weight" validate:"required"`
	Image      string   `json:"image" validate:"required"`

	BaseHP int `json:"basehp" validate:"required"`
	MinHP  int `json:"minhp" validate:"required"`
	MaxHP  int `json:"maxhp" validate:"required"`

	BaseAttack int `json:"baseattack" validate:"required"`
	MinAttack  int `json:"minattack" validate:"required"`
	MaxAttack  int `json:"maxattack" validate:"required"`

	BaseDefense int `json:"basedefense" validate:"required"`
	MinDefense  int `json:"mindefense" validate:"required"`
	MaxDefense  int `json:"maxdefense" validate:"required"`

	BaseSPAttack int `json:"basespattack" validate:"required"`
	MinSPAttack  int `json:"minspattack" validate:"required"`
	MaxSPAttack  int `json:"maxspattack" validate:"required"`

	BaseSPDefense int `json:"basespdefense" validate:"required"`
	MinSPDefense  int `json:"minspdefense" validate:"required"`
	MaxSPDefense  int `json:"maxspdefense" validate:"required"`

	BaseSpeed int `json:"basespeed" validate:"required"`
	MinSpeed  int `json:"minspeed" validate:"required"`
	MaxSpeed  int `json:"maxspeed" validate:"required"`

	Total int `json:"total" validate:"required"`

	PrevEvo string `json:"prevevo"`
	NextEvo string `json:"nextevo"`
}
