package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekschultz/jurassic-park/controllers"
	"github.com/derekschultz/jurassic-park/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func TestFindCages(t *testing.T) {
	router := setupRouter()
	models.ConnectDatabase()

	router.GET("/cages", controllers.FindCages)
	req, err := http.NewRequest("GET", "/cages", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
