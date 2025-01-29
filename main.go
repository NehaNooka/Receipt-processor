package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var (
	receipts sync.Map
)

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)
	router.Run(":8080")
}

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if !validateReceipt(receipt) {
		c.JSON(400, gin.H{"error": "Invalid receipt data"})
		return
	}

	points, err := calculatePoints(receipt)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid receipt data"})
		return
	}

	id := uuid.New().String()
	receipts.Store(id, points)

	c.JSON(200, gin.H{"id": id})
}

func getPoints(c *gin.Context) {
	id := c.Param("id")
	points, exists := receipts.Load(id)
	if !exists {
		c.JSON(404, gin.H{"error": "Receipt not found"})
		return
	}
	c.JSON(200, gin.H{"points": points})
}

func validateReceipt(receipt Receipt) bool {
	retailerRegex := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerRegex.MatchString(receipt.Retailer) {
		return false
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return false
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return false
	}

	if len(receipt.Items) < 1 {
		return false
	}

	itemDescRegex := regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
	for _, item := range receipt.Items {
		if !itemDescRegex.MatchString(item.ShortDescription) {
			return false
		}
		if !priceRegex.MatchString(item.Price) {
			return false
		}
	}

	if !priceRegex.MatchString(receipt.Total) {
		return false
	}

	return true
}

func calculatePoints(receipt Receipt) (int, error) {
	points := 0

	// Rule 1: Alphanumeric characters in retailer name
	alnumCount := 0
	for _, c := range receipt.Retailer {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			alnumCount++
		}
	}
	points += alnumCount

	// Rule 2: Round dollar amount
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0, err
	}
	totalCents := int(math.Round(total * 100))
	if totalCents%100 == 0 {
		points += 50
	}

	// Rule 3: Multiple of 0.25
	if totalCents%25 == 0 {
		points += 25
	}

	// Rule 4: 5 points per two items
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: Item description length multiple of 3
	for _, item := range receipt.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: Purchase day odd
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return 0, err
	}
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: Time between 2pm-4pm
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return 0, err
	}
	if hour := purchaseTime.Hour(); hour >= 14 && hour < 16 {
		points += 10
	}

	return points, nil
}
