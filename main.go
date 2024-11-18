package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"savannah.go/sms"
	"savannah.go/auth"
	"context"
	"golang.org/x/oauth2"
)

var DB *sql.DB

type Customer struct {
	Id    int    `json:"cust_id"`
	Code  string `json:"code"`
	Name  string `json:"full_name"`
	Phone string `json:"phone"`
}

type Order struct {
	Id             int       `json:"order_id"`
	Item           string    `json:"item"`
	Time           time.Time `json:"time"`
	Amount         float64   `json:"amount"`
	RelatedCustomer int      `json:"cust_id"`
}

func init() {
	if os.Getenv("RENDER") != "true" { 
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
}


func main() {
	auth.InitOIDC() // Ensure OIDC is initialized before usage

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s DBname=%s sslmode=disable",
		"127.0.0.1",
		"5432",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		"savannah",
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connection Successful!")

	router := gin.Default()
	fmt.Printf("%T\n", auth.GetOAuth2Config())

	
	router.GET("/login", func(c *gin.Context) {
		authURL := auth.GetOAuth2Config().AuthCodeURL("state-string", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	})
	
	
	router.GET("/oauth/callback", func(c *gin.Context) {
		ctx := context.Background()
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code missing"})
			return
		}
	
		// Exchange the code for a token
		oauth2Token, err := auth.GetOAuth2Config().Exchange(ctx, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
			return
		}
	
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID token missing"})
			return
		}
	
		// Validate the ID token
		idToken, err := auth.GetVerifier().Verify(ctx, rawIDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}
	
		// Extract user claims from the ID token
		var claims struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := idToken.Claims(&claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			return
		}
	
		// Return the tokens in JSON
		c.JSON(http.StatusOK, gin.H{
			"access_token": oauth2Token.AccessToken,
			"id_token":     rawIDToken,
			"expires_in": time.Until(oauth2Token.Expiry).Seconds(),
			"token_type":   "Bearer",
			"user": gin.H{
				"email": claims.Email,
				"name":  claims.Name,
			},
		})
	})
	
	router.GET("/customers", auth.OIDCAuthMiddleware(), GetCustomers)
	router.POST("/customers", auth.OIDCAuthMiddleware(), PostCustomers)
	router.GET("/orders", auth.OIDCAuthMiddleware(), GetOrders)
	router.POST("/orders", auth.OIDCAuthMiddleware(), PostOrders)

	router.Run("localhost:8080")
}

func GetCustomers(c *gin.Context) {
	rows, err := DB.Query("SELECT cust_id, code, full_name, phone FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query customers"})
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.Id, &customer.Code, &customer.Name, &customer.Phone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan customer"})
			return
		}
		customers = append(customers, customer)
	}

	c.IndentedJSON(http.StatusOK, customers)
}

func PostCustomers(c *gin.Context) {
	var newCustomer Customer

	if err := c.BindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := DB.QueryRow("INSERT INTO customers (code, full_name, phone) VALUES ($1, $2, $3) RETURNING cust_id", newCustomer.Code, newCustomer.Name, newCustomer.Phone).Scan(&newCustomer.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert customer"})
		return
	}

	c.JSON(http.StatusCreated, newCustomer)
}

func GetOrders(c *gin.Context) {
	rows, err := DB.Query("SELECT order_id, item, amount, time FROM orders")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.Id, &order.Item, &order.Amount, &order.Time); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}
		orders = append(orders, order)
	}

	c.IndentedJSON(http.StatusOK, orders)
}

func PostOrders(c *gin.Context) {
	var newOrder Order

	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := DB.QueryRow("INSERT INTO orders (item, time, amount, cust_id) VALUES ($1, $2, $3, $4) RETURNING order_id", newOrder.Item, newOrder.Time, newOrder.Amount, newOrder.RelatedCustomer).Scan(&newOrder.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order"})
		return
	}

	// recipient := "+254728333926"
	recipients := []string{"+254728333926"}
	message := fmt.Sprintf("Order #%d for %s successfully created. Amount: KES %.2f", newOrder.Id, newOrder.Item, newOrder.Amount)
	senderID := "SAVANNAH INF"
	sms.SendSMS(recipients, message, senderID)

	c.JSON(http.StatusCreated, newOrder)
}
