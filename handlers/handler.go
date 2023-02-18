package handlers

import (
	"abrnoc_ch/models"
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
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
	query := "SELECT SUM(price) FROM subscriptions WHERE customer_id = $1 AND is_active = true"
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

// Calculate subscription charges for a 10-minute period for each customer
func CalculateSubscriptionCharges() {
	// Get the current time and the time 10 minutes ago
	now := time.Now()
	tenMinutesAgo := now.Add(-10 * time.Minute)

	// Query the database for active subscriptions that started before the last 10 minutes
	query := "SELECT id, name, price, period, customer_id FROM subscriptions WHERE is_active = true AND start_date <= $1 AND end_date > $2"
	rows, err := DB.Query(query, now, tenMinutesAgo)
	if err != nil {
		log.Println("Error querying subscriptions:", err)
		return
	}
	defer rows.Close()

	// Iterate over the rows and calculate the charges for each subscription
	for rows.Next() {
		var id, customerID int
		var startDate, endDate time.Time
		var amount float64
		err = rows.Scan(&id, &customerID, &startDate, &endDate, &amount)
		if err != nil {
			log.Println("Error scanning subscription row:", err)
			continue
		}
		duration := now.Sub(startDate)
		if duration > 10*time.Minute {
			duration = 10 * time.Minute
		}
		charge := amount * duration.Minutes() / 60.0
		log.Printf("Charge $%.2f for subscription %d of customer %d\n", charge, id, customerID)
		// TODO: Create an invoice for the customer with the calculated charge
	}
	if err = rows.Err(); err != nil {
		log.Println("Error iterating over subscription rows:", err)
		return
	}
}

// Handler to close a customer's subscription and generate an invoice
func CloseSubscription(c echo.Context) error {
	// Get the customer ID and subscription ID from the request
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)
	if err != nil {
		log.Println("Error parsing customer ID:", err)
		return c.String(http.StatusBadRequest, "Invalid customer ID")
	}
	subscriptionID, err := strconv.ParseInt(c.Param("subscription_id"), 10, 64)
	if err != nil {
		log.Println("Error parsing subscription ID:", err)
		return c.String(http.StatusBadRequest, "Invalid subscription ID")
	}

	// Check if the subscription exists and belongs to the customer
	subscription := models.Subscription{}
	err = DB.QueryRow("SELECT id, name, price, is_active, period, customer_id FROM subscriptions WHERE id = $1 AND customer_id = $2", subscriptionID, customerID).Scan(&subscription.ID, &subscription.Name, &subscription.Price, &subscription.IsActive, &subscription.Period, &subscription.CustomerID)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "Subscription not found")
	} else if err != nil {
		log.Println("Error querying subscription:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Update the subscription to mark it as inactive and set the end date to the current time
	end := time.Now()
	_, err = DB.Exec("UPDATE subscriptions SET is_active = false, end_date = $1 WHERE id = $2", end, subscriptionID)
	if err != nil {
		log.Println("Error updating subscription:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Calculate the charge for the subscription based on the duration
	duration := end.Sub(time.Unix(int64(subscription.Period), 0))
	charge := subscription.Price * duration.Minutes() / 60.0

	// Create an invoice for the customer with the calculated charge
	invoice := models.Invoice{
		ID:             int64(uuid.New().ID()),
		StartTime:      time.Unix(int64(subscription.Period), 0),
		EndTime:        end,
		Amount:         charge,
		SubscriptionID: subscription.ID,
	}
	_, err = DB.Exec("INSERT INTO invoices (id, start_time, end_time, amount, subscription_id) VALUES ($1, $2, $3, $4, $5)", invoice.ID, invoice.StartTime, invoice.EndTime, invoice.Amount, invoice.SubscriptionID)
	if err != nil {
		log.Println("Error creating invoice:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, invoice)
}
