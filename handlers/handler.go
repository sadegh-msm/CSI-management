package handlers

import (
	"abrnoc_ch/models"
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

var db *sql.DB

func SetDB(DB *sql.DB) {
	db = DB
}

func CreateCustomerHandler(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	credit, _ := strconv.ParseFloat(c.FormValue("credit"), 64)

	res, err := db.Exec("INSERT INTO customers (name, email, credit) VALUES ($1, $2, $3)", name, email, credit)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create customer"})
	}

	id, _ := res.LastInsertId()
	return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "name": name, "email": email, "credit": credit})
}

func GetCustomerHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var customer models.Customer
	row := db.QueryRow("SELECT * FROM customers WHERE id=$1", id)
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Credit, &customer.CreatedAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
	}

	return c.JSON(http.StatusOK, customer)
}

func DeleteCustomerHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := db.Exec("DELETE FROM customers WHERE id=$1", id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete customer"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
	}

	return c.NoContent(http.StatusOK)
}

func CreateSubscriptionHandler(c echo.Context) error {
	customerID, _ := strconv.Atoi(c.FormValue("customer_id"))
	plan := c.FormValue("plan")
	duration, _ := strconv.Atoi(c.FormValue("duration"))
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)

	expiresAt := time.Now().Add(time.Duration(duration) * time.Minute)

	res, err := db.Exec("INSERT INTO subscriptions (customer_id, plan, duration, price, expires_at) VALUES ($1, $2, $3, $4, $5)", customerID, plan, duration, price, expiresAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create subscription"})
	}

	id, _ := res.LastInsertId()
	return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "customer_id": customerID, "plan": plan, "duration": duration, "price": price, "expires_at": expiresAt})
}

func GetSubscriptionHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var subscription models.Subscription
	row := db.QueryRow("SELECT * FROM subscriptions WHERE id=$1", id)
	err := row.Scan(&subscription.ID, &subscription.CustomerID, &subscription.Plan, &subscription.Duration, &subscription.Price, &subscription.ExpiresAt, &subscription.CreatedAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Subscription not found"})
	}

	return c.JSON(http.StatusOK, subscription)
}

func DeleteSubscriptionHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := db.Exec("DELETE FROM subscriptions WHERE id=$1", id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete subscription"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Subscription not found"})
	}

	return c.NoContent(http.StatusOK)
}

func GetInvoiceHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	row := db.QueryRow("SELECT * FROM invoices WHERE id=$1", id)
	err := row.Scan(&invoice.ID, &invoice.CustomerID, &invoice.SubscriptionID, &invoice.Amount, &invoice.CreatedAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invoice not found"})
	}

	return c.JSON(http.StatusOK, invoice)
}

func ChargeInvoiceHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	row := db.QueryRow("SELECT * FROM invoices WHERE id=$1", id)
	err := row.Scan(&invoice.ID, &invoice.CustomerID, &invoice.SubscriptionID, &invoice.Amount, &invoice.CreatedAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invoice not found"})
	}

	var customer models.Customer
	row = db.QueryRow("SELECT * FROM customers WHERE id=$1", invoice.CustomerID)
	err = row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Credit, &customer.CreatedAt)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to charge invoice"})
	}

	if customer.Credit < invoice.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient credit"})
	}

	newCredit := customer.Credit - invoice.Amount
	_, err = db.Exec("UPDATE customers SET credit=$1 WHERE id=$2", newCredit, customer.ID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to charge invoice"})
	}

	return c.NoContent(http.StatusOK)
}

func DeleteInvoiceHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := db.Exec("DELETE FROM invoices WHERE id=$1", id)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete invoice"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invoice not found"})
	}

	return c.NoContent(http.StatusOK)
}

func ListInvoicesHandler(c echo.Context) error {
	customerID, _ := strconv.Atoi(c.QueryParam("customer_id"))

	var invoices []models.Invoice
	var rows *sql.Rows
	var err error
	if customerID == 0 {
		rows, err = db.Query("SELECT * FROM invoices")
	} else {
		rows, err = db.Query("SELECT * FROM invoices WHERE customer_id=$1", customerID)
	}
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list invoices"})
	}
	defer rows.Close()

	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(&invoice.ID, &invoice.CustomerID, &invoice.SubscriptionID, &invoice.Amount, &invoice.CreatedAt)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list invoices"})
		}

		invoices = append(invoices, invoice)
	}

	return c.JSON(http.StatusOK, invoices)
}

func CreateInvoiceHandler(c echo.Context) error {
	var invoice models.Invoice

	// Bind the request body to the invoice struct
	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Retrieve the subscription from the database
	var subscription models.Subscription
	if err := db.QueryRow("SELECT id, customer_id, price FROM subscriptions WHERE id = ?", invoice.SubscriptionID).Scan(&subscription.ID, &subscription.CustomerID, &subscription.Price); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid subscription ID"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve subscription"})
	}

	// Retrieve the customer from the database
	var customer models.Customer
	if err := db.QueryRow("SELECT id, name, email, credit FROM customers WHERE id = ?", subscription.CustomerID).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Credit); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid customer ID"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve customer"})
	}

	// Calculate the amount for the invoice
	invoice.Amount = subscription.Price

	// Create the invoice
	result, err := db.Exec("INSERT INTO invoices (subscription_id, amount) VALUES (?, ?)", invoice.SubscriptionID, invoice.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create invoice"})
	}

	// Get the ID of the created invoice
	invoiceID, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to retrieve invoice ID"})
	}
	invoice.ID = int(invoiceID)

	return c.JSON(http.StatusCreated, invoice)
}
