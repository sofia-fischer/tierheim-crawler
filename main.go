package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"tierheim-crawler/database"
	"tierheim-crawler/models"
	"tierheim-crawler/repositories"
	"time"
)

func getDogs(requestContext *gin.Context) {
	collector := colly.NewCollector()

	//
	//collector.OnError(func(_ *colly.Response, err error) {
	//	log.Println("Something went wrong: ", err)
	//})
	//
	//collector.OnResponse(func(r *colly.Response) {
	//	fmt.Println("Page visited: ", r.Request.URL)
	//})
	//
	collector.OnHTML("main", func(element *colly.HTMLElement) {
		foundDog, err := repositories.FromHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog.FetchedAt = time.Now()
		dog := models.Dog{ShelterIdentifier: foundDog.ShelterIdentifier}
		dog = repositories.UpdateOrCreate(dog, foundDog)

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})

	collector.Visit("https://tierschutzverein-muenchen.de/tiervermittlung/tierheim/hunde/200078")

	//c.IndentedJSON(http.StatusOK, testDogs)
}

func main() {
	database.ConnectDb()

	router := gin.Default()
	router.GET("/dogs", getDogs)

	router.Run("localhost:8080")
}
