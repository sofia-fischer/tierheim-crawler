package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"tierheim-crawler/models"
)

func ConnectDb() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("dog_crawler.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("ConnectDb:: Failed to connect to database. \n", err)
	}

	log.Println("ConnectDb:: connected to database")
	log.Println("ConnectDb:: running migrations")
	err = database.AutoMigrate(&models.Dog{})

	if err != nil {
		log.Fatal("ConnectDb:: Failed to run migrations. \n", err)
	}

	return database
}
