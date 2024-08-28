package db

import (
	"order-service/config"
	"order-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(config *config.Config) {

	database, err := gorm.Open(postgres.Open(config.DBConnectionString), &gorm.Config{})
	if err != nil {
		panic("Error al conectar a la base de datos")
	}

	database.AutoMigrate(&models.Order{})

	DB = database
}
