package main

import (
	"github.com/gin-gonic/gin"

	"github.com/derekschultz/jurassic-park/controllers"
	"github.com/derekschultz/jurassic-park/models"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	r.GET("/species", controllers.FindSpecies)
	r.GET("/species/:name", controllers.FindSpeciesByName)
	r.POST("/species", controllers.CreateSpecies)
	r.PATCH("/species/:id", controllers.UpdateSpecies)
	r.DELETE("/species/:id", controllers.DeleteSpecies)

	r.GET("/cages", controllers.FindCages)
	r.GET("/cage/:id", controllers.FindCage)
	r.GET("/cages/:status", controllers.FindCagesStatus)
	r.POST("/cages", controllers.CreateCage)
	r.PATCH("/cage/:id", controllers.UpdateCage)
	r.DELETE("/cage/:id", controllers.DeleteCage)

	r.Run(":8080")
}
