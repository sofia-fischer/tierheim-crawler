package repositories

import (
	"errors"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"tierheim-crawler/models"
	"time"
)

type DogRepository struct {
	Database *gorm.DB
}

const germanTimeLayout = "02.01.2006"

func (repository DogRepository) UpdateOrCreate(dogData models.Dog) models.Dog {

	var existingDogs []models.Dog
	repository.Database.Limit(1).Where("shelter_identifier = ?", dogData.ShelterIdentifier).Find(&existingDogs)

	if len(existingDogs) == 0 {
		dogData.ID = uuid.NewString()
		repository.Database.Create(&dogData)
		return dogData
	}

	repository.Database.Model(existingDogs[0]).Updates(dogData)

	return dogData
}

func FromHtml(element *colly.HTMLElement) (models.Dog, error) {
	foundDog := models.Dog{}
	var err error

	foundDog.Name = element.ChildText("H1")
	foundDog.ShelterIdentifier = element.ChildText("div.text-lg.font-bold")

	if foundDog.ShelterIdentifier == "" {
		return foundDog, errors.New("FromHtml::could not find shelter identifier")
	}

	element.ForEach("li", func(i int, lineElement *colly.HTMLElement) {
		switch {
		case strings.Contains(lineElement.Text, "Hunde, "):
			foundDog.Breed = strings.Replace(lineElement.Text, "Hunde, ", "", 1)
		case strings.Contains(lineElement.Text, "Geschlecht: Weiblich"):
			foundDog.IsMale = false
		case strings.Contains(lineElement.Text, "Geschlecht: Männlich"):
			foundDog.IsMale = true
		case strings.Contains(lineElement.Text, "Geburtstag: "):
			birthday := strings.Replace(lineElement.Text, "Geburtstag: ", "", 1)
			foundDog.BornAt, err = time.Parse(germanTimeLayout, birthday)
		case strings.Contains(lineElement.Text, "Farben: "):
			foundDog.Color = strings.Replace(lineElement.Text, "Farben: ", "", 1)
		}
	})

	text := element.ChildText("div.prose")
	spaceRegex := regexp.MustCompile(`(\s|\n)+`)
	foundDog.Description = spaceRegex.ReplaceAllString(text, " ")

	digitOnlyRegex := regexp.MustCompile(`\d+`)
	weightRegex := regexp.MustCompile(` (\d+) kg(\s|\.)`)
	weight := strings.TrimSpace(weightRegex.FindString(foundDog.Description))
	foundDog.Weight, _ = strconv.Atoi(digitOnlyRegex.FindString(weight))
	heightRegex := regexp.MustCompile(` (\d+) cm(\s|\.)`)
	height := strings.TrimSpace(heightRegex.FindString(foundDog.Description))
	foundDog.Height, _ = strconv.Atoi(digitOnlyRegex.FindString(height))
	foundDog.FetchedAt = time.Now()

	return foundDog, err
}
