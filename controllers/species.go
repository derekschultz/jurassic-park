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
	gorm.Model
	CageID   int    `json:"cageId"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
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
		coexist := matchSpecies(c, input.CageID)

		if !coexist {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "The provided species cannot coexist in the same cage!"})
			return
		}

		// Check status of cage
		if err := models.DB.Model(&models.Cage{}).Where("id = ?", input.CageID).Find(&cageStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return
		}
		fmt.Printf("cageID: %v, cageStatus: %v\n", input.CageID, cageStatus.Status) // DEBUG
		if cageStatus.Status == "DOWN" {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Oops, cannot move to cage with DOWN status!"})
			return
		}
	}

	models.DB.Model(&species).Where("id = ?", c.Param("id")).Updates(input)

	if input.Quantity != 0 {
		// Find capacity of cage_id in species table
		if err := models.DB.Model(&species).Where("id = ?", c.Param("id")).Select("cage_id, SUM(quantity) as capacity").Find(&cageStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return
		}

		fmt.Printf("cage_id: %v, capacity: %v\n", cageStatus.CageID, cageStatus.Capacity) // DEBUG
		fmt.Println("\nHOLD ONTO YOUR BUTTS!")

		// Update capacity field in cages table
		if err := models.DB.Model(&models.Cage{}).Where("id = ?", cageStatus.CageID).Update("capacity", cageStatus.Capacity).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return
		}
	}

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

// matchSpecies takes in a *gin.Context with an int for the cage ID
func matchSpecies(c *gin.Context, cageId int) bool {
	var cur, new UpdateSpeciesInput

	if c.Param("id") != "" {
		// Get the diet of species in species.cage_id
		if err := models.DB.Model(&models.Species{}).Where("id = ?", c.Param("id")).Select("name, diet").Find(&cur).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return false
		}
		// Get the diet of species in the updated cageID
		if err := models.DB.Model(&models.Species{}).Where("id = ?", cageId).Select("name, diet").Find(&new).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": err.Error()})
			return false
		}

		fmt.Printf("current cageID: %v, species: %v, diet: %v\n", cur.CageID, cur.Name, cur.Diet) // DEBUG
		fmt.Printf("new cageID: %v, species: %v, diet: %v\n", cageId, new.Name, new.Diet)         // DEBUG

		// Carnivores cannot exist in a cage with any other species
		if cur.Diet != new.Diet {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Cannot mix species with different diets!"})
			return false
		} else if cur.Diet == "Carnivore" {
			// Only allow if species name matches
			if cur.Name == new.Name {
				fmt.Printf("cur.Name %v matches new.Name %v\n", cur.Name, new.Name) // DEBUG
				return true
			} else {
				fmt.Printf("cur.Name %v does NOT match new.Name %v\n", cur.Name, new.Name) // DEBUG
				c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Different species of Carnivores cannot exist in the same cage!"})
				return false
			}
		}
	} else {
		// No ID provided
		c.JSON(http.StatusBadRequest, gin.H{"statusCode": 400, "error": "Please provide a cageID!"})
		return false
	}

	return true
}
