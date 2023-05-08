package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"tierheim-crawler/database"
	"tierheim-crawler/models"
	"tierheim-crawler/support"
)

func index(requestContext *gin.Context) {
	repository := models.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}
	service := support.NewCrawlerService(repository)
	dogs := service.CrawlIndex()

	requestContext.IndentedJSON(http.StatusOK, dogs)
}

func show(requestContext *gin.Context) {
	repository := models.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}
	dog := repository.FirstOrCreateByShelterId(requestContext.Param("id"))
	service := support.NewCrawlerService(repository)
	updatedDog := service.CrawlDetails(dog)
	requestContext.IndentedJSON(http.StatusOK, updatedDog)
}

func main() {
	db := database.ConnectDb()
	router := gin.Default()

	// Add database context to request context
	router.Use(addDatabaseContextMiddleware(db))

	router.GET("/dogs", index)
	router.GET("/dogs/:id", show)

	router.Run("localhost:8080")
}

func addDatabaseContextMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("database", db)
	}
}
