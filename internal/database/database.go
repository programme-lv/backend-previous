package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
}

func NewDatabase(dbConnString string) *Database {
	gorm, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Database{Gorm: gorm}
}
