package middleware

import (
	"encoding/json" // Para decodificar respuestas JSON
	"fmt"           // Para formatear strings
	"net/http"      // Para realizar solicitudes HTTP
	"os"            // Para acceder a variables de entorno

	"github.com/labstack/echo/v4" // Framework web Echo
)

// AuthMiddleware es un middleware para autenticar las solicitudes
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtiene el token de autorización del encabezado de la solicitud
		token := c.Request().Header.Get("Authorization")

		// Verifica si el token está presente
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token no proporcionado"})
		}

		// Valida el token con el servicio de usuarios
		userID, err := validateTokenWithUserService(token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error al validar el token"})
		}

		// Verifica si el userID es válido
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
		}

		// Almacena el userID en el contexto para su uso posterior
		c.Set("user_id", userID)
		// Continúa con el siguiente manejador
		return next(c)
	}
}

// validateTokenWithUserService valida el token con el servicio de usuarios
func validateTokenWithUserService(token string) (uint, error) {
	// Obtiene la URL del servicio de usuarios desde las variables de entorno
	userServiceURL := os.Getenv("APP_USER_SERVICE_URL")

	// Construye la URL para la validación del token
	url := fmt.Sprintf("%s/validate", userServiceURL)
	// Crea una nueva solicitud GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creando la solicitud: %v", err)
	}

	// Establece el token en el encabezado de autorización
	req.Header.Set("Authorization", token)

	// Crea un cliente HTTP y realiza la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error haciendo la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Verifica si la respuesta es exitosa
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error de autenticación: código de estado %d", resp.StatusCode)
	}

	// Estructura para decodificar la respuesta JSON
	var response struct {
		User struct {
			ID uint `json:"id"`
		} `json:"user"`
	}

	// Decodifica la respuesta JSON
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decodificando la respuesta: %v", err)
	}

	// Verifica si el ID de usuario es válido
	if response.User.ID == 0 {
		return 0, fmt.Errorf("ID de usuario no válido")
	}

	// Devuelve el ID de usuario
	return response.User.ID, nil
}
