// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"

	"github.com/Atarod61/myproject/send" // مسیر صحیح پکیج send
)

var (
	urlsToCheck = []string{
		"https://dl.jzac.ir/",
		"https://google.com",
	}
	db *sql.DB
	mu sync.Mutex
)

// replace with your MySQL connection string
const dsn = "root:21@tcp(127.0.0.1:3308)/fiber"

func initDB() *sql.DB {
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// Adjust your table creation query based on your needs.
	statement, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS website_status (
			url VARCHAR(255) PRIMARY KEY,
			status VARCHAR(50)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func checkWebsite(url string) {
	for {
		resp, err := fiber.Get(url)
		mu.Lock()
		if err != nil || resp.StatusCode() != fiber.StatusOK {
			updateStatus(url, "DOWN")
			send.SendSMS() // فراخوانی از پکیج send
		} else {
			updateStatus(url, "UP")
		}
		mu.Unlock()
		if resp != nil {
			resp.ResBodyClose()
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}

func updateStatus(url string, status string) {
	statement, err := db.Prepare("INSERT INTO website_status (url, status) VALUES (?, ?) ON DUPLICATE KEY UPDATE status = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(url, status, status)
	if err != nil {
		log.Fatal(err)
	}
}

func getStatus(url string) string {
	row := db.QueryRow("SELECT status FROM website_status WHERE url = ?", url)
	var status string
	err := row.Scan(&status)
	if err != nil {
		// Log the error, but return a default value
		log.Printf("Error getting status for %s: %v", url, err)
		return "UNKNOWN"
	}
	return status
}

func statusHandler(c *fiber.Ctx) error {
	html := `<html><head><title>Status Page</title></head><body>
             <h1>Website Status</h1><ul>`
	for _, url := range urlsToCheck {
		status := getStatus(url)
		html += fmt.Sprintf("<li>%s: <strong>%s</strong></li>", url, status)
	}
	html += `</ul></body></html>`
	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

func main() {
	db = initDB()
	defer db.Close()

	for _, url := range urlsToCheck {
		go checkWebsite(url)
	}

	app := fiber.New()
	app.Get("/status", statusHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(app.Listen(":8080"))
}
