package models

import "gorm.io/gorm"

// Cages can have status of 'ACTIVE' or 'DOWN' and capacity >= 0
type Cage struct {
	gorm.Model
	// ID          int `json:"id" gorm:"primary_key"`
	CageID      int `json:"cageId"`
	Capacity    int `json:"capacity" gorm:"default:0"`
	MaxCapacity int `json:"maxCapacity" gorm:"default:4"`
	// Cage should have info regarding what species is being contained
	Species []Species `json:"species" gorm:"foreignKey:CageID"`                // has many association
	Status  string    `json:"status" gorm:"check:status IN ('ACTIVE','DOWN')"` // required
}
