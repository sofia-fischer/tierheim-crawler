package support

import (
	"errors"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strconv"
	"strings"
	"tierheim-crawler/models"
	"time"
)

type DogIdentifierService struct {
	baseUrl    string
	showQuery  string
	indexQuery string
}

type DogIdentifierServiceInterface interface {
	dogShow(identifier string, callBack colly.HTMLCallback) error
	dogIndex(callBack colly.HTMLCallback) error
	fromShowHtml(element *colly.HTMLElement) (models.Dog, error)
	fromIndexHtml(element *colly.HTMLElement) (models.Dog, error)
}

func NewDogIdentifierService() *DogIdentifierService {
	return &DogIdentifierService{
		baseUrl:    "https://tierschutzverein-muenchen.de/tiervermittlung/tierheim/hunde",
		showQuery:  "main:not(.modal__content)",
		indexQuery: "div.tsv-tiervermittlung-animal-name",
	}
}

func (service *DogIdentifierService) dogShow(identifier string, callBack colly.HTMLCallback) error {
	collector := colly.NewCollector()
	collector.OnHTML(service.showQuery, callBack)
	err := collector.Visit(service.baseUrl + "/" + identifier)
	collector.Wait()

	return err
}

func (service *DogIdentifierService) dogIndex(callBack colly.HTMLCallback) error {
	collector := colly.NewCollector()
	collector.OnHTML(service.indexQuery, callBack)
	err := collector.Visit(service.baseUrl)
	collector.Wait()

	return err
}

func (service *DogIdentifierService) fromShowHtml(element *colly.HTMLElement) (models.Dog, error) {
	foundDog := models.Dog{}

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
			var err error
			foundDog.BornAt, err = time.Parse(models.GermanTimeLayout, birthday)

			if err != nil {
				log.Println("getDogs:: error while formatting dog", err)
			}
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

	return foundDog, nil
}

func (service *DogIdentifierService) fromIndexHtml(element *colly.HTMLElement) (models.Dog, error) {
	foundDog := models.Dog{}

	foundDog.Name = element.ChildText("H3")
	foundDog.ShelterIdentifier = element.ChildText(".tsv-tiervermittlung-animal-id")

	if foundDog.ShelterIdentifier == "" {
		return foundDog, errors.New("FromHtml::could not find shelter identifier")
	}

	foundDog.FetchedAt = time.Now()

	return foundDog, nil
}
