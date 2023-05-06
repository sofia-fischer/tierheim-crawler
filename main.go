package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tierheim-crawler/database"
	"tierheim-crawler/models"
	"tierheim-crawler/support"
)

func index(requestContext *gin.Context) {
	repository := models.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}

	_ = support.DogIndex(func(element *colly.HTMLElement) {
		foundDog, err := support.FromIndexHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = repository.UpdateOrCreate(foundDog)

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})
}

func show(requestContext *gin.Context) {
	repository := models.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}
	id := requestContext.Param("id")

	_ = support.DogShow(id, func(element *colly.HTMLElement) {
		foundDog, err := support.FromShowHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = repository.UpdateOrCreate(foundDog)

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})
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
