package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tierheim-crawler/crawler"
	"tierheim-crawler/database"
	"tierheim-crawler/repositories"
)

func index(requestContext *gin.Context) {
	repository := repositories.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}

	crawler.DogFind("200078", func(element *colly.HTMLElement) {
		foundDog, err := repositories.FromHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = repository.UpdateOrCreate(foundDog)

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})
}

func show(requestContext *gin.Context) {
	repository := repositories.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}
	id := requestContext.Param("id")

	crawler.DogFind(id, func(element *colly.HTMLElement) {
		foundDog, err := repositories.FromHtml(element)

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
