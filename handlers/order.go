package handlers

import (
	"encoding/json"        // Para codificar/decodificar JSON
	"fmt"                  // Para formatear strings
	"net/http"             // Para manejar solicitudes HTTP
	"order-service/db"     // Paquete local para interactuar con la base de datos
	"order-service/models" // Paquete local con los modelos de datos
	"os"                   // Para interactuar con el sistema operativo

	"github.com/labstack/echo/v4" // Framework web Echo
)

// CreateOrder maneja la creación de una nueva orden
func CreateOrder(c echo.Context) error {
	// Obtiene el ID del usuario del contexto
	userID := c.Get("user_id").(uint)
	// Obtiene el token de autorización del encabezado de la solicitud
	token := c.Request().Header.Get("Authorization")

	var order models.Order
	// Vincula los datos de la solicitud a la estructura Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos de solicitud inválidos"})
	}

	// Asigna el ID del usuario a la orden
	order.UserID = userID

	// Valida el producto con el servicio de productos
	productValid, price := validateProductWithProductService(order.ProductID, token)
	if !productValid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Producto no válido"})
	}

	// Calcula el total de la orden
	order.Total = float64(order.Quantity) * price

	// Crea la orden en la base de datos
	if err := db.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear la orden"})
	}

	// Devuelve la orden creada
	return c.JSON(http.StatusCreated, order)
}

// GetUserOrders obtiene todas las órdenes de un usuario
func GetUserOrders(c echo.Context) error {
	// Obtiene el ID del usuario del contexto
	userID := c.Get("user_id").(uint)

	var orders []models.Order
	// Busca todas las órdenes del usuario en la base de datos
	if err := db.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al obtener las órdenes"})
	}

	// Devuelve las órdenes encontradas
	return c.JSON(http.StatusOK, orders)
}

// validateProductWithProductService valida un producto con el servicio de productos
func validateProductWithProductService(productID uint, token string) (bool, float64) {
	// Obtiene la URL del servicio de productos desde las variables de entorno
	productServiceUrl := os.Getenv("APP_PRODUCT_SERVICE_URL")

	// Construye la URL para la solicitud al servicio de productos
	url := fmt.Sprintf("%s/products/%d", productServiceUrl, productID)

	// Crea una nueva solicitud GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0
	}

	// Establece el token de autorización en el encabezado de la solicitud
	req.Header.Set("Authorization", token)

	// Crea un cliente HTTP y realiza la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	// Verifica si la respuesta es exitosa
	if resp.StatusCode != http.StatusOK {
		return false, 0
	}

	// Estructura para decodificar la respuesta JSON del producto
	var product struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	// Decodifica la respuesta JSON
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return false, 0
	}

	// Verifica si el producto es válido y tiene cantidad disponible
	if product.ID == 0 || product.Quantity <= 0 {
		return false, 0
	}

	// Devuelve true y el precio del producto si es válido
	return true, product.Price
}
