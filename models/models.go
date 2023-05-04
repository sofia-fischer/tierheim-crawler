package models

import "time"

type Dog struct {
	ID                string    `json:"id"`
	ShelterIdentifier string    `json:"identifier"`
	Name              string    `json:"name"`
	Breed             string    `json:"breed"`
	Color             string    `json:"color"`
	IsMale            bool      `json:"is_male" gorm:"default:false;not null"`
	Description       string    `json:"description"`
	BornAt            time.Time `json:"born_at"`
	AdoptedAt         time.Time `json:"adopted_at"`
	LastUpdatedAt     time.Time `json:"last_updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}
