package jobs

import (
	"tierheim-crawler/database"
	"tierheim-crawler/models"
)

func crawlIndex() {
	// setup database connection
	repository := models.DogRepository{Database: database.ConnectDb()}

	// list of all crawled dogs
	fetchedDogs := make([]models.Dog, 0)

	// crawl index page

	// mark dogs that are not in the list as adopted
	repository.MarkAdoptedWhereNotIn(fetchedDogs)
}
