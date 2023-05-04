package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"net/http"
	"tierheim-crawler/database"
	"time"
)

type dogModel struct {
	ID                string    `json:"id"`
	ShelterIdentifier string    `json:"identifier"`
	Name              string    `json:"name"`
	Breed             string    `json:"breed"`
	Color             string    `json:"color"`
	IsMale            bool      `json:"is_male"`
	Description       string    `json:"description"`
	BornAt            time.Time `json:"born_at"`
	AdoptedAt         time.Time `json:"adopted_at"`
	LastUpdatedAt     time.Time `json:"last_updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

//
//var testDogs = []dogModel{
//	{ID: "1", ShelterIdentifier: "345", Name: "Alo", Breed: "Husky"},
//	{ID: "2", ShelterIdentifier: "346", Name: "Bunny", Breed: "Poodle"},
//	{ID: "3", ShelterIdentifier: "348", Name: "Charly", Breed: "Corgi"},
//}

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
		foundDog, err := FromHtml(element)

		if err != nil {
			requestContext.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		requestContext.IndentedJSON(http.StatusOK, foundDog)
	})
	//
	//collector.OnScraped(func(r *colly.Response) {
	//	requestContext.IndentedJSON(http.StatusOK, r)
	//})

	collector.Visit("https://tierschutzverein-muenchen.de/tiervermittlung/tierheim/hunde/200078")

	//c.IndentedJSON(http.StatusOK, testDogs)
}

func main() {
	database.ConnectDb()

	router := gin.Default()
	router.GET("/dogs", getDogs)

	router.Run("localhost:8080")
}
