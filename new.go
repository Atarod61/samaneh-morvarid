package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	urlsToCheck = []string{
		"https://exampleadddd.com",
		"https://google.com",
	}
	statusMap = make(map[string]string)
	mu        sync.Mutex
)

// تابع ارسال پیامک
func sendSMS(alerturl string) {
	url := "https://panel.asanak.com/webservice/v1rest/sendsms"
	str := "username=xxxxxxxxxxx&password=xxxxxx" +
		"source=98xxxxxxxxx&destination=98xxxxxxxx&message=سلام مشکل پیش آمد: " + alerturl + " is DOWN!"
	payload := strings.NewReader(str)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Log request details
	fmt.Println("Sending SMS request to:", url)
	fmt.Println("Payload:", str)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Log response details
	fmt.Println("Response status:", res.Status)
	fmt.Println("Response body:", string(body))
}

func checkWebsite(url string) {
	for {
		resp, err := http.Get(url)
		mu.Lock()
		if err != nil || resp.StatusCode != http.StatusOK {
			statusMap[url] = "DOWN"
			sendAlert(url) // Call sendAlert
		} else {
			statusMap[url] = "UP"
		}
		mu.Unlock()
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(30 * time.Second) // Check every 30 seconds
	}
}

// اصلاح تابع sendAlert برای فراخوانی sendSMS
func sendAlert(url string) {
	fmt.Printf("ALERT: %s is DOWN!\n", url)
	sendSMS(url) // Send SMS alert
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<html><head><title>Status Page</title></head><body>")
	fmt.Fprintln(w, "<h1>Website Status</h1><ul>")
	for url, status := range statusMap {
		fmt.Fprintf(w, "<li>%s: <strong>%s</strong></li>", url, status)
	}
	fmt.Fprintln(w, "</ul></body></html>")
}

func main() {
	for _, url := range urlsToCheck {
		go checkWebsite(url)
	}
	http.HandleFunc("/status", statusHandler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
