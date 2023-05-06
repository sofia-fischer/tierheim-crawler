package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Dog struct {
	gorm.Model

	ID                string    `gorm:"type:varchar(20);column:id;next:uuid"`
	ShelterIdentifier string    `json:"identifier" gorm:"uniqueIndex"`
	Name              string    `json:"name"`
	Breed             string    `json:"breed"`
	Color             string    `json:"color"`
	IsMale            bool      `json:"is_male"`
	Weight            int       `json:"weight"`
	Height            int       `json:"height"`
	Description       string    `json:"description"`
	BornAt            time.Time `json:"born_at"`
	AdoptedAt         time.Time `json:"adopted_at"`
	FetchedAt         time.Time `json:"fetched_at_at"`
}

type DogRepository struct {
	Database *gorm.DB
}

func (repository DogRepository) UpdateOrCreate(dogData Dog) Dog {

	var existingDogs []Dog
	repository.Database.Limit(1).Where("shelter_identifier = ?", dogData.ShelterIdentifier).Find(&existingDogs)

	if len(existingDogs) == 0 {
		dogData.ID = uuid.NewString()
		repository.Database.Create(&dogData)
		return dogData
	}

	repository.Database.Model(existingDogs[0]).Updates(dogData)

	return dogData
}
