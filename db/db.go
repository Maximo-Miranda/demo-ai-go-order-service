// Package db maneja la conexión y operaciones con la base de datos
package db

import (
	"order-service/config" // Importa el paquete de configuración local
	"order-service/models" // Importa los modelos de datos locales

	"gorm.io/driver/postgres" // Driver de PostgreSQL para GORM
	"gorm.io/gorm"            // ORM (Object-Relational Mapping) para Go
)

// DB es una variable global que contiene la conexión a la base de datos
var DB *gorm.DB

// ConnectDatabase establece la conexión con la base de datos y realiza la migración automática
func ConnectDatabase(config *config.Config) {

	// Intenta abrir una conexión a la base de datos usando la cadena de conexión proporcionada
	database, err := gorm.Open(postgres.Open(config.DBConnectionString), &gorm.Config{})
	if err != nil {
		// Si hay un error al conectar, se detiene la ejecución del programa
		panic("Error al conectar a la base de datos")
	}

	// Realiza la migración automática del modelo Order
	database.AutoMigrate(&models.Order{})

	// Asigna la conexión de base de datos a la variable global DB
	DB = database
}
