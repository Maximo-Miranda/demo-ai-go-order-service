package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token no proporcionado"})
		}

		userID, err := validateTokenWithUserService(token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al validar el token"})
		}

		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}

func validateTokenWithUserService(token string) (uint, error) {

	userServiceURL := os.Getenv("APP_USER_SERVICE_URL")

	url := fmt.Sprintf("%s/validate", userServiceURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando la solicitud: %v", err)
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error haciendo la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error de autenticación: código de estado %d", resp.StatusCode)
	}

	var response struct {
		User struct {
			ID uint `json:"id"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decodificando la respuesta: %v", err)
	}

	if response.User.ID == 0 {
		return 0, fmt.Errorf("ID de usuario no válido")
	}

	return response.User.ID, nil
}
