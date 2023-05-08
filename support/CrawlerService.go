package support

import (
	"github.com/gocolly/colly"
	"log"
	"tierheim-crawler/models"
)

type CrawlerService struct {
	repository models.DogRepository
}

func NewCrawlerService(repository models.DogRepository) CrawlerService {
	return CrawlerService{repository: repository}
}

func (service CrawlerService) CrawlIndex() []models.Dog {
	// list of all crawled dogs
	fetchedDogs := make([]models.Dog, 0)
	// list of created dogs
	createdDogs := make([]models.Dog, 0)

	// crawl index page
	_ = DogIndex(func(element *colly.HTMLElement) {
		foundDog, err := FromIndexHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = service.repository.UpdateOrCreate(foundDog)
		fetchedDogs = append(fetchedDogs, foundDog)

		if foundDog.CreatedAt == foundDog.UpdatedAt {
			createdDogs = append(createdDogs, foundDog)
		}
	})

	// mark dogs that are not in the list as adopted
	service.repository.MarkAdoptedWhereNotIn(fetchedDogs)

	// crawl details for new dogs
	for _, dog := range createdDogs {
		dog = service.CrawlDetails(dog)
	}

	return fetchedDogs
}

func (service CrawlerService) CrawlDetails(dog models.Dog) models.Dog {
	_ = DogShow(dog.ShelterIdentifier, func(element *colly.HTMLElement) {
		foundDog, err := FromShowHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog.ID = dog.ID

		foundDog = service.repository.UpdateOrCreate(foundDog)
	})

	return dog
}
