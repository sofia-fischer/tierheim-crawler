package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
	"log"
	"net/http"
	"tierheim-crawler/database"
	"tierheim-crawler/repositories"
)

func getDogs(requestContext *gin.Context) {
	collector := colly.NewCollector()
	repository := repositories.DogRepository{Database: requestContext.MustGet("database").(*gorm.DB)}

	collector.OnHTML("main:not(.modal__content)", func(element *colly.HTMLElement) {
		foundDog, err := repositories.FromHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = repository.UpdateOrCreate(foundDog)

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})

	collector.Visit("https://tierschutzverein-muenchen.de/tiervermittlung/tierheim/hunde/200078")
}

func main() {
	db := database.ConnectDb()
	router := gin.Default()

	// Add database context to request context
	router.Use(addDatabaseContextMiddleware(db))

	router.GET("/dogs", getDogs)

	router.Run("localhost:8080")
}

func addDatabaseContextMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("database", db)
	}
}
