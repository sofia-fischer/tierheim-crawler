package repositories

import (
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"tierheim-crawler/database"
	"tierheim-crawler/models"
	"time"
)

const germanTimeLayout = "01.02.2006"

func UpdateOrCreate(existing models.Dog, updating models.Dog) models.Dog {
	dog := existing
	database.DB.FirstOrCreate(&dog, existing)
	database.DB.Model(dog).Updates(updating)

	return dog
}

func FromHtml(element *colly.HTMLElement) (models.Dog, error) {
	foundDog := models.Dog{}
	var err error

	foundDog.Name = element.ChildText("H1")
	foundDog.ShelterIdentifier = element.ChildText("div.text-lg.font-bold")

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

	return foundDog, err
}
