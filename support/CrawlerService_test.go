package support

import (
	"github.com/gocolly/colly"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"tierheim-crawler/models"
	"time"
)

type MockDogIdentifierService struct {
	mock.Mock
}

func (mock *MockDogIdentifierService) dogIndex(callBack colly.HTMLCallback) error {
	callBack(&colly.HTMLElement{})
	return nil
}

func (mock *MockDogIdentifierService) dogShow(identifier string, callBack colly.HTMLCallback) error {
	callBack(&colly.HTMLElement{})
	return nil
}

func (mock *MockDogIdentifierService) fromIndexHtml(element *colly.HTMLElement) (models.Dog, error) {
	return models.Dog{ShelterIdentifier: "3", ID: "unique-id-3", Name: "Green"}, nil
}

func (mock *MockDogIdentifierService) fromShowHtml(element *colly.HTMLElement) (models.Dog, error) {
	bornAt, _ := time.Parse(models.GermanTimeLayout, "01.01.2020")

	return models.Dog{
		ShelterIdentifier: "1234",
		Name:              "Test",
		Breed:             "Husky",
		Color:             "Black",
		Weight:            2,
		Height:            3,
		IsMale:            true,
		Description:       "So testy",
		BornAt:            bornAt,
		FetchedAt:         time.Now(),
	}, nil
}

func TestCrawlDetails(t *testing.T) {
	mockedService := new(MockDogIdentifierService)
	db := connectTestingDb()

	service := CrawlerService{
		repository:           models.DogRepository{db},
		dogIdentifierService: mockedService,
	}

	dog := service.CrawlDetails(models.Dog{ShelterIdentifier: "1", Name: "Blue"})
	db.Create(&dog)
	crawledDog := service.CrawlDetails(dog)
	if crawledDog.Name != "Test" {
		t.Errorf("Expected dog name to be 'Test', got %s", dog.Name)
	}
	if crawledDog.ID != dog.ID {
		t.Errorf("Expected dog id to match the inital dog id, got %s", dog.ID)
	}
}

func TestCrawlIndex(t *testing.T) {
	mockedService := new(MockDogIdentifierService)
	db := connectTestingDb()

	service := CrawlerService{
		repository:           models.DogRepository{db},
		dogIdentifierService: mockedService,
	}

	// unadopted Dog that will be marked as adopted
	db.Create(&models.Dog{ShelterIdentifier: "1", ID: "unique-id", Name: "Blue"})
	// adopted Dog that will be not touched
	db.Create(&models.Dog{ShelterIdentifier: "2", ID: "unique-id-2", Name: "Red Herring", AdoptedAt: time.Now()})
	// unadopted Dog that will be fetched
	db.Create(&models.Dog{ShelterIdentifier: "3", ID: "unique-id-3", Name: "Green"})

	dogs := service.CrawlIndex()

	if len(dogs) != 1 {
		t.Errorf("Expected 1 dog, got %d", len(dogs))
	}

	if dogs[0].ShelterIdentifier != "3" {
		t.Errorf("Expected first dog to be '3', got %s", dogs[0].Name)
	}

	var adopted models.Dog
	db.Where("shelter_identifier = ?", "1").First(&adopted)
	if adopted.AdoptedAt.Year() != time.Now().Year() {
		t.Errorf("Expected adopted dog to be marked as adopted, got %s", adopted.AdoptedAt)
	}
}

func connectTestingDb() *gorm.DB {
	database, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	_ = database.AutoMigrate(&models.Dog{})

	return database
}
