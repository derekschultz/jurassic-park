package controllers

import (
	"net/http"

	"github.com/derekschultz/jurassic-park/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /species
// Get all species
func FindSpecies(c *gin.Context) {
	var species []models.Species
	models.DB.Find(&species)

	c.JSON(http.StatusOK, gin.H{"statusCode": c.Writer.Status(), "data": species})
}

type CreateSpeciesInput struct {
	Name     string `json:"name" binding:"required"`
	Diet     string `json:"diet" binding:"required"`
	Quantity int    `json:"quantity"`
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
	Diet     string `json:"diet"`
	Quantity int    `json:"quantity"`
}

// UPDATE /species/:id
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
