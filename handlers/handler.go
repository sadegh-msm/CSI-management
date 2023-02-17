package handlers

import (
	"abrnoc_ch/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func CreateCustomer(c echo.Context) error {
	customer := new(models.Customer)

	if err := c.Bind(customer); err != nil {
		return err
	}

	query := "INSERT INTO customers (ID, username, credit) VALUES ($1, $2, $3) RETURNING id"

	err := DB.QueryRow(query, customer.ID, customer.Username, customer.Credit).Scan(&customer.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, customer)
}

func GetCustomer(c echo.Context) error {
	id := c.Param("id")
	var customer models.Customer

	query := "SELECT * FROM customers WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&customer.ID, &customer.Username, &customer.Credit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)
}

func UpdateCustomer(c echo.Context) error {
	id := c.Param("id")
	customer := new(models.Customer)

	if err := c.Bind(customer); err != nil {
		return err
	}

	query := "UPDATE customers SET ID = $1, username = $2, credit = $3 WHERE id = $6"
	_, err := DB.Exec(query, customer.Username, customer.Credit, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)
}

func GetActiveSubscriptions(c echo.Context) error {
	id := c.Param("id")

	var subscriptions []models.Subscription

	query := "SELECT * FROM subscriptions WHERE customer_id = $1 AND is_active = true"
	rows, err := DB.Query(query, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subscription models.Subscription

		err = rows.Scan(&subscription.ID, &subscription.CustomerID, &subscription.StartDate, &subscription.EndDate, &subscription.Amount, &subscription.CreatedAt, &subscription.UpdatedAt, &subscription.IsActive)
		if err != nil {
			return err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return c.JSON(http.StatusOK, subscriptions)
}

func CreateInvoice(c echo.Context) error {
	invoice := new(models.Invoice)

	if err := c.Bind(invoice); err != nil {
		return err
	}
	query := "INSERT INTO invoices (customer_id, subscription_id, amount, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	now := time.Now()
	err := DB.QueryRow(query, invoice.CustomerID, invoice.SubscriptionID, invoice.Amount, now).Scan(&invoice.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, invoice)
}
