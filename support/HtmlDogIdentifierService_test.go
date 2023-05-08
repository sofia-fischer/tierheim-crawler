package support

import (
	"github.com/gocolly/colly"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFromHtmlShow(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html><html>
<main>
    <div class="fluid-container">
        <div class="p-8 pb-12 flex">
            <h1 class="text-4xl leading-h1">HARCOS</h1>
            <div class="text-lg font-bold flex-grow">171066</div>
            <ul class="grid-2">
                <li><strong>Hunde, Holländischer Schäferhund</strong></li>
                <li><strong>Farben:</strong> braun, schwarz</li>
                <li><strong>Geschlecht:</strong> Männlich</li>
                <li><strong>Kastriert:</strong> Nein</li>
                <li><strong>Geburtstag:</strong> 05.07.2020</li>
                <li><strong>Nur für erfahrene Halter</strong></li>
            </ul>
        </div>
        <section class="mt-20">
            <div class="prose">
                <h2 class="text-h2">Informationen zu HARCOS:</h2>
                <p class="text-intro">Beschreibung. Schulterhöhe von 60 cm und wiegt aktuell 30 kg.</p>
                <h3 class="leading-h3">Charakter</h3>
                <p class="leading-relaxed text-black">Charakter Informationen.</p>
                <h3 class="leading-h3">Besonderheiten</h3>
                <p class="leading-relaxed text-black ">Besondere Informationen</p>
            </div>
        </section>
    </div>
</main>
</html>`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	dogIdentifierService := NewDogIdentifierService()
	collector := colly.NewCollector()
	collector.OnHTML("/", func(element *colly.HTMLElement) {
		foundDog, err := dogIdentifierService.fromShowHtml(element)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if foundDog.Name != "HARCOS" {
			t.Errorf("Expected name to be HARCOS, but was %v", foundDog.Name)
		}

		if foundDog.ShelterIdentifier != "171066" {
			t.Errorf("Expected shelter identifier to be 171066, but was %v", foundDog.ShelterIdentifier)
		}

		if foundDog.Breed != "Holländischer Schäferhund" {
			t.Errorf("Expected breed to be Holländischer Schäferhund, but was %v", foundDog.Breed)
		}

		if foundDog.Color != "braun, schwarz" {
			t.Errorf("Expected color to be braun, schwarz, but was %v", foundDog.Color)
		}

		if foundDog.IsMale != true {
			t.Errorf("Expected isMale to be true, but was %v", foundDog.IsMale)
		}

		if foundDog.BornAt.Year() != 2020 {
			t.Errorf("Expected year to be 2015, but was %v", foundDog.BornAt.Year())
		}

		if foundDog.Description != "Informationen zu HARCOS: Beschreibung. Schulterhöhe von 60 cm und wiegt aktuell 30 kg. Charakter Charakter Informationen. Besonderheiten Besondere Informationen" {
			t.Errorf("Expected description to be \"Informationen zu HARCOS: Beschreibung. Schulterhöhe von 60 cm und wiegt aktuell 30 kg. Charakter Charakter Informationen. Besonderheiten Besondere Informationen\", but was \"%v\"", foundDog.Description)
		}

		if foundDog.Weight != 30 {
			t.Errorf("Expected weight to be 30, but was %v", foundDog.Weight)
		}

		if foundDog.Height != 60 {
			t.Errorf("Expected height to be 60, but was %v", foundDog.Height)
		}
	})

	collector.Visit(server.URL)
}

func TestFromHtmlIndex(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html><html>
<div class="tsv-tiervermittlung-animal-name flex justify-center lg:justify-between items-center">
      <h3 class="font-bold lg:flex-grow">AMY</h3>
                  <span class="text-xs hidden lg:block tsv-tiervermittlung-animal-id">200078</span>
    </div>
</html>`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	dogIdentifierService := NewDogIdentifierService()

	collector := colly.NewCollector()
	collector.OnHTML("/", func(element *colly.HTMLElement) {
		foundDog, err := dogIdentifierService.fromIndexHtml(element)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if foundDog.Name != "AMY" {
			t.Errorf("Expected name to be HARCOS, but was %v", foundDog.Name)
		}

		if foundDog.ShelterIdentifier != "200078" {
			t.Errorf("Expected shelter identifier to be 171066, but was %v", foundDog.ShelterIdentifier)
		}
	})

	collector.Visit(server.URL)
}
