package routes

import (
	"abrnoc_ch/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Router creating new router and add middlewares and routes and return echo object
func Router() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/customers", handlers.CreateCustomer)
	e.GET("/customers/:id", handlers.GetCustomer)
	e.PUT("/customers/:id", handlers.UpdateCustomer)
	e.GET("/customers/:id/subscriptions", handlers.GetActiveSubscriptions)
	e.POST("/invoices", handlers.CreateInvoice)

	return e
}
