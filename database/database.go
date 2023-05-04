package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"tierheim-crawler/models"
)

type Database struct {
	Db *gorm.DB
}

var DB Database

func ConnectDb() {
	// data source name string to connect to postgres
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UT",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// connect to postgres database (or handle error)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("ConnectDb:: Failed to connect to database. \n", err)
	}

	log.Println("ConnectDb:: connected to database")
	database.Logger = logger.Default.LogMode(logger.Info)

	log.Println("ConnectDb:: running migrations")
	err = database.AutoMigrate(&models.Dog{})

	if err != nil {
		log.Fatal("ConnectDb:: Failed to run migrations. \n", err)
	}

	// set global database variable
	DB = Database{
		Db: database,
	}
}
