package models

import (
	"gorm.io/gorm" // Importa el paquete gorm para el mapeo objeto-relacional
)

// Order representa la estructura de una orden en la base de datos
type Order struct {
	gorm.Model         // Incorpora los campos básicos de gorm (ID, CreatedAt, UpdatedAt, DeletedAt)
	UserID     uint    `json:"user_id"`    // ID del usuario que realizó la orden
	ProductID  uint    `json:"product_id"` // ID del producto ordenado
	Quantity   int     `json:"quantity"`   // Cantidad del producto ordenado
	Total      float64 `json:"total"`      // Total de la orden (precio * cantidad)
}
