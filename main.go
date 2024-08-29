package main

import (
	"order-service/config"
	"order-service/db"
	"order-service/handlers"
	"order-service/middleware"
	"time"

	"github.com/labstack/echo/v4"
	echo_middlewares "github.com/labstack/echo/v4/middleware"
)

func main() {
	conf := config.Config{}
	conf.LoadConfig()

	if conf.APPEnv != "development" {
		time.Sleep(time.Second * 5)
	}
	db.ConnectDatabase(&conf)

	e := echo.New()
	e.Use(echo_middlewares.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Order service is running!")
	})

	// Grupo de rutas protegidas
	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware)

	// Rutas de órdenes
	protected.POST("/orders", handlers.CreateOrder)
	protected.GET("/orders", handlers.GetUserOrders)

	e.Logger.Fatal(e.Start(":8082"))
}
