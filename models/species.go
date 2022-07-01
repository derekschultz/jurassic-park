package models

import "gorm.io/gorm"

// Species has Name, e.g. Tyrannosaurus and Diet, e.g. carnivore
type Species struct {
	gorm.Model
	// ID     int    `json:"id" gorm:"primary_key"`
	SpeciesID uint
	Name      string `json:"name"`
	Diet      string `json:"diet"`
	CageID    int    `json:"cageId"`
	Quantity  int    `json:"quantity" gorm:"default:1"`
}
