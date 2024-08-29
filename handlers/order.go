package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/db"
	"order-service/models"
	"os"

	"github.com/labstack/echo/v4"
)

func CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	token := c.Request().Header.Get("Authorization")

	var order models.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos de solicitud inválidos"})
	}

	order.UserID = userID

	// Validar el producto con el servicio de productos
	productValid, price := validateProductWithProductService(order.ProductID, token)
	if !productValid {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Producto no válido"})
	}

	order.Total = float64(order.Quantity) * price

	if err := db.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al crear la orden"})
	}

	return c.JSON(http.StatusCreated, order)
}

func GetUserOrders(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var orders []models.Order
	if err := db.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al obtener las órdenes"})
	}

	return c.JSON(http.StatusOK, orders)
}

func validateProductWithProductService(productID uint, token string) (bool, float64) {

	productServiceUrl := os.Getenv("APP_PRODUCT_SERVICE_URL")

	url := fmt.Sprintf("%s/products/%d", productServiceUrl, productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, 0
	}

	var product struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return false, 0
	}

	if product.ID == 0 || product.Quantity <= 0 {
		return false, 0
	}

	return true, product.Price
}
