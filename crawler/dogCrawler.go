package crawler

import (
	"github.com/gocolly/colly"
)

const baseUrl = "https://tierschutzverein-muenchen.de/tiervermittlung/tierheim/hunde"

func DogFind(identifier string, callBack colly.HTMLCallback) error {
	collector := colly.NewCollector()
	collector.OnHTML("main:not(.modal__content)", callBack)
	return collector.Visit(baseUrl + "/" + identifier)
}

func DogIndex(callBack colly.HTMLCallback) error {
	collector := colly.NewCollector()
	collector.OnHTML("div.tsv-tiervermittlung-animal-name", callBack)
	return collector.Visit(baseUrl)
}
