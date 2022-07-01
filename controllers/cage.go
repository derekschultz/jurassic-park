package controllers

import (
	"fmt"
	"net/http"

	"github.com/derekschultz/jurassic-park/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /cages
// Get all cages
func FindCages(c *gin.Context) {
	var cage []models.Cage
	if err := models.DB.Find(&cage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": cage})
}

// GET /cage/:id
// Get a cage by id
func FindCage(c *gin.Context) {
	var cage []models.Cage
	if err := models.DB.Model(&models.Cage{}).Where("id = ?", c.Param("id")).Preload("Species").Find(&cage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": cage})
}

// GET /cages/:status
// Filter cages by status
func FindCagesStatus(c *gin.Context) {
	var cage []models.Cage
	if err := models.DB.Model(&models.Cage{}).Where("status = ?", c.Param("status")).Find(&cage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": cage})
}

type CreateCageInput struct {
	Capacity    int              `json:"capacity"`
	MaxCapacity int              `json:"maxCapacity"`
	Species     []models.Species `json:"species"`
	Status      string           `json:"status" binding:"required"`
}

// POST /cages
// Create new cage
func CreateCage(c *gin.Context) {
	// Validate input
	var input CreateCageInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	if input.Capacity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, capacity is less than zero!"})
		return
	}
	if input.MaxCapacity < input.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, max capacity is less than capacity!"})
		return
	}

	// Create cage in DB
	cage := models.Cage{Capacity: input.Capacity, MaxCapacity: input.MaxCapacity, Status: input.Status}
	models.DB.Create(&cage)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": cage})
}

type UpdateCageInput struct {
	// need gorm.Model to avoid panic on setting Status
	gorm.Model
	CageID      int              `json:"cageId"`
	Capacity    int              `json:"capacity"`
	MaxCapacity int              `json:"maxCapacity"`
	Species     []models.Species `json:"species" gorm:"foreignKey:CageID"`
	Status      string           `json:"status"`
}

// PATCH /cage/:id
// Update existing cage by id
func UpdateCage(c *gin.Context) {
	// Get cage by ID if it exists
	var cage []models.Cage
	if err := models.DB.Where("id = ?", c.Param("id")).First(&cage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	// Validate input
	var input UpdateCageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	if input.Capacity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, capacity is less than zero!"})
		return
	}

	if input.Capacity != 0 && input.MaxCapacity != 0 {
		if input.Capacity > input.MaxCapacity {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, capacity exceeds max capacity!"})
			return
		}
	}

	if input.Capacity == 0 {
		models.DB.Model(&cage).Where("id = ?", c.Param("id")).Updates(map[string]interface{}{"capacity": input.Capacity})
	}

	// Check if cage has dinosaurs in it, if so, cannot toggle power from 'ACTIVE' to 'DOWN'
	var cageCapacity int
	models.DB.Model(&cage).Where("id = ?", c.Param("id")).Select("capacity").Find(&cageCapacity)
	fmt.Printf("cageCapacity: %v", cageCapacity) // DEBUG
	if input.Status == "DOWN" && cageCapacity > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, cannot power down cage with capacity > 0!"})
		return
	}

	models.DB.Model(&cage).Where("id = ?", c.Param("id")).Updates(input)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": cage})
}

// DELETE /cage/:id
// Delete existing cage by id
func DeleteCage(c *gin.Context) {
	var cage models.Cage
	if err := models.DB.Where("id = ?", c.Param("id")).First(&cage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	models.DB.Delete(&cage)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": true})
}
