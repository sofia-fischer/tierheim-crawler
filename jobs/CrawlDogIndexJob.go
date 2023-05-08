package jobs

import (
	"github.com/gocolly/colly"
	"log"
	"tierheim-crawler/database"
	"tierheim-crawler/models"
	"tierheim-crawler/support"
)

func crawlIndex() {
	// setup database connection
	repository := models.DogRepository{Database: database.ConnectDb()}

	// list of all crawled dogs
	fetchedDogs := make([]models.Dog, 0)

	// crawl index page
	_ = support.DogIndex(func(element *colly.HTMLElement) {
		foundDog, err := support.FromIndexHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = repository.UpdateOrCreate(foundDog)
		fetchedDogs = append(fetchedDogs, foundDog)
	})

	// mark dogs that are not in the list as adopted
	repository.MarkAdoptedWhereNotIn(fetchedDogs)
}
