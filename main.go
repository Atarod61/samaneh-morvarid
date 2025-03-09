// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Atarod61/samaneh-morvarid/send"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	urlsToCheck = []string{
		"https://https://dl.jzac.ir/",
		"https://google.com",
	}
	statusMap = make(map[string]string)
	mu        sync.Mutex
	db        *gorm.DB
)

type WebsiteStatus struct {
	ID        uint   `gorm:"primaryKey"`  // اصلاح این خط
	URL       string `gorm:"uniqueIndex"` // اصلاح این خط
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func checkWebsite(url string) {
	for {
		resp, err := http.Get(url)
		mu.Lock()
		if err != nil || resp.StatusCode != http.StatusOK {
			statusMap[url] = "DOWN"
			// Send SMS if status changes to DOWN
			if updateStatusInDB(url, "DOWN") {
				sendAlert(url)
			}
		} else {
			statusMap[url] = "UP"
			updateStatusInDB(url, "UP")
		}
		mu.Unlock()
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}

func updateStatusInDB(url, status string) bool {
	var websiteStatus WebsiteStatus
	result := db.Where("url = ?", url).First(&websiteStatus)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new record
			websiteStatus = WebsiteStatus{URL: url, Status: status}
			db.Create(&websiteStatus)
			return true // Status changed
		} else {
			log.Println("Error finding record:", result.Error)
			return false
		}
	}

	if websiteStatus.Status != status {
		// Update existing record if status has changed
		websiteStatus.Status = status
		db.Save(&websiteStatus)
		return true // Status changed
	}
	return false // Status not changed
}

func sendAlert(url string) {
	fmt.Printf("ALERT: %s is DOWN!\n", url)
	// Send SMS here
	err := send.SendSMS(fmt.Sprintf("ALERT: %s is DOWN!", url)) // Corrected line
	if err != nil {
		log.Println("Error sending SMS:", err)
	}
}

func statusHandler(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<html><head><title>Status Page</title></head><body>")
	htmlBuilder.WriteString("<h1>Website Status</h1><ul>")

	for url, status := range statusMap {
		htmlBuilder.WriteString(fmt.Sprintf("<li>%s: <strong>%s</strong></li>", url, status))
	}

	htmlBuilder.WriteString("</ul></body></html>")

	// اصلاح به کد مناسب
	err := c.SendString(htmlBuilder.String())
	if err != nil {
		log.Println("Error sending response:", err)
		return err
	}
	return nil // اضافه کردن این خط
}

func main() {
	// Database connection
	dsn := "root:21@tcp(127.0.0.1:3308)/fiber?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate
	db.AutoMigrate(&WebsiteStatus{})

	// Initialize Fiber app
	app := fiber.New()

	// Start website monitoring in goroutines
	for _, url := range urlsToCheck {
		go checkWebsite(url)
	}

	// Define status route
	app.Get("/status", statusHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
