package controllers

import (
	"fmt"
	"net/http"

	"github.com/derekschultz/jurassic-park/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /species
// Get all species
func FindSpecies(c *gin.Context) {
	var species []models.Species
	if err := models.DB.Find(&species).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": species})
}

type CreateSpeciesInput struct {
	Name     string `json:"name" binding:"required"`
	CageID   int    `json:"cageId"`
	Diet     string `json:"diet" binding:"required"`
	Quantity int    `json:"quantity"`
}

// GET /species/:name
// Filter species by name
func FindSpeciesByName(c *gin.Context) {
	var species []models.Species
	if err := models.DB.Where("name = ?", c.Param("name")).Find(&species).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": species})
}

// POST /species
// Create new species
func CreateSpecies(c *gin.Context) {
	// Validate input
	var input CreateSpeciesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	// Create species in DB
	species := models.Species{Name: input.Name, Diet: input.Diet, Quantity: input.Quantity}

	models.DB.Create(&species)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": species})
}

type UpdateSpeciesInput struct {
	gorm.Model
	Name     string `json:"name"`
	CageID   int    `json:"cageId"`
	Diet     string `json:"diet"`
	Quantity int    `json:"quantity"`
}

type CageStatus struct {
	CageID int    `json:"cageId"`
	Status string `json:"status"`
}

// PATCH /species/:id
// Update existing species
func UpdateSpecies(c *gin.Context) {
	// Get species by ID if it exists
	var species []models.Species
	if err := models.DB.Where("id = ?", c.Param("id")).First(&species).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	// Validate input
	var input UpdateSpeciesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	// FIXME: allow user to set cageId to zero value
	var cageStatus CageStatus
	if input.CageID != 0 {
		// Check status of cage
		if err := models.DB.Model(&models.Cage{}).Where("id = ?", input.CageID).Find(&cageStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return
		}
		fmt.Printf("cageID: %v, cageStatus: %v", input.CageID, cageStatus.Status)
		if cageStatus.Status == "DOWN" {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, cannot move to cage with DOWN status!"})
			return
		}
	}

	models.DB.Model(&species).Where("id = ?", c.Param("id")).Updates(input)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": species})
}

// DELETE /species/:id
// Delete existing species
func DeleteSpecies(c *gin.Context) {
	var species []models.Species
	if err := models.DB.Where("id = ?", c.Param("id")).First(&species).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
		return
	}

	models.DB.Delete(&species)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": true})
}
