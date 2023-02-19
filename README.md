# Customer Subscription and Invoice Management System

This project is a basic implementation of a customer subscription and invoice management system. It allows customers to create subscriptions and view their invoices. Customers can also delete subscriptions and invoices, and adjust their charges based on their available credit.

### Prerequisites

To run this application, you will need to have the following installed:
- Golang
- PostgreSQL
- Echo
 
## Getting Started

1. Clone this repository to your local machine:
```bash 
git clone https://github.com/sadegh-msm/CSI-management.git
````
2. Change to the project directory:
```bash 
cd customer-subscription-invoice
````
3. Install dependencies:
```go  
go mod download
````
4. Start the application:
```bash
go run main.go
````
5. The application will now be running on `http://localhost:8080`.

## API Endpoints

### Customers
- `GET /customers/:id`: List all customers
- `POST /customers`: Create a new customer
- `DELETE /customers/:id`: Delete a customer with the given ID
### Subscriptions
- `GET /subscriptions/:id`: List all subscriptions
- `POST /subscriptions`: Create a new subscription
- `DELETE /subscriptions/:id`: Delete a subscription with the given ID
### Invoices
- `GET /invoices`: List all invoices for all customers, or list invoices for a specific customer by providing the customer_id query parameter
- `POST /invoices`: Create a new invoice
- `DELETE /invoices/:id`: Delete an invoice with the given ID
- `PUT /invoices/:id/charge`: Charge an invoice by deducting the amount from the customer's credit
- `GET /invoices/:id`: Returns the invoice with specific id.