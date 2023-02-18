package handlers

import (
	"abrnoc_ch/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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

		err = rows.Scan(&subscription.ID, &subscription.Name, &subscription.Price, &subscription.IsActive, &subscription.Period, &subscription.CustomerID)
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
	query := "INSERT INTO invoices (ID, start_time, end_time, amount, subscription_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	now := time.Now()
	err := DB.QueryRow(query, invoice.ID, now, invoice.EndTime, invoice.Amount, invoice.SubscriptionID).Scan(&invoice.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, invoice)
}

func CreateSubscription(c echo.Context) error {
	subscription := new(models.Subscription)
	if err := c.Bind(subscription); err != nil {
		return err
	}
	query := "INSERT INTO subscriptions (id, name, price, is_active, period, customer_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := DB.QueryRow(query, subscription.ID, subscription.Name, subscription.Price, subscription.IsActive, subscription.Period, subscription.CustomerID).Scan(&subscription.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, subscription)
}

func CalculateInvoicePrice(c echo.Context) error {
	customerID, err := strconv.Atoi(c.FormValue("customer_id"))
	if err != nil {
		return err
	}
	query := "SELECT SUM(amount) FROM subscriptions WHERE customer_id = $1 AND is_active = true"
	var totalPrice float64
	err = DB.QueryRow(query, customerID).Scan(&totalPrice)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"customer_id": customerID,
		"total_price": totalPrice,
	})
}
