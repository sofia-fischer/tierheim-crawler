package support

import (
	"github.com/gocolly/colly"
	"log"
	"tierheim-crawler/models"
)

type CrawlerService struct {
	repository           models.DogRepository
	dogIdentifierService DogIdentifierServiceInterface
}

func NewCrawlerService(repository models.DogRepository) CrawlerService {
	return CrawlerService{repository, NewDogIdentifierService()}
}

func (service *CrawlerService) CrawlIndex() []models.Dog {
	// list of all crawled dogs
	fetchedDogs := make([]models.Dog, 0)
	// list of created dogs
	createdDogs := make([]models.Dog, 0)

	// crawl index page
	_ = service.dogIdentifierService.dogIndex(func(element *colly.HTMLElement) {
		foundDog, err := service.dogIdentifierService.fromIndexHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog = service.repository.UpdateOrCreate(foundDog)
		fetchedDogs = append(fetchedDogs, foundDog)

		if foundDog.CreatedAt.Date() == foundDog.UpdatedAt.Date() {
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

func (service *CrawlerService) CrawlDetails(dog models.Dog) models.Dog {

	_ = service.dogIdentifierService.dogShow(dog.ShelterIdentifier, func(element *colly.HTMLElement) {
		foundDog, err := service.dogIdentifierService.fromShowHtml(element)

		if err != nil {
			log.Println("getDogs:: error while formatting dog", err)
			return
		}

		foundDog.ID = dog.ID

		dog = service.repository.UpdateOrCreate(foundDog)
	})

	return dog
}
