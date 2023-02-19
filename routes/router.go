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

	e.POST("/customers", handlers.CreateCustomerHandler)
	e.GET("/customers/:id", handlers.GetCustomerHandler)
	e.DELETE("/customers/:id", handlers.DeleteCustomerHandler)

	e.POST("/subscriptions", handlers.CreateSubscriptionHandler)
	e.GET("/subscriptions/:id", handlers.GetSubscriptionHandler)
	e.DELETE("/subscriptions/:id", handlers.DeleteSubscriptionHandler)

	e.GET("/invoices/:id", handlers.GetInvoiceHandler)
	e.PUT("/invoices/:id/charge", handlers.ChargeInvoiceHandler)
	e.DELETE("/invoices/:id", handlers.DeleteInvoiceHandler)
	e.GET("/invoices", handlers.ListInvoicesHandler)

	return e
}
