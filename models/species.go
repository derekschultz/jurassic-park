package models

import "gorm.io/gorm"

// Species has Name, e.g. Tyrannosaurus and Diet, e.g. carnivore
type Species struct {
	gorm.Model
	// ID     int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	CageID   int    `json:"cageId"`
	Diet     string `json:"diet"`
	Quantity int    `json:"quantity" gorm:"default:1"`
}
