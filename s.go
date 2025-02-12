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
		"https://dl.jzac.ir.com",
		"https://google.com",
	}
	statusMap = make(map[string]string)
	mu        sync.Mutex

	// SMS Configuration - Replace with your actual credentials
	smsURL         = "https://panel.asanak.com/webservice/v1rest/sendsms"
	smsUsername    = "xxxxxxxxx"  // Replace with your username
	smsPassword    = xxxxxxx"      // Replace with your password
	smsSource      = "xxxxxxxhttp" // Replace with your source number
	smsDestination = "98xxxxxxx" // Replace with your destination number
)

func checkWebsite(url string) {
	for {
		resp, err := http.Get(url)
		mu.Lock()
		if err != nil || resp.StatusCode != http.StatusOK {
			statusMap[url] = "DOWN"
			sendAlert(url, "Website "+url+" is DOWN!") // Include URL in the alert message

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

func sendAlert(url string, message string) {
	fmt.Printf("ALERT: %s is DOWN!\n", url)

	// Construct SMS Payload
	smsMessage := fmt.Sprintf("ALERT: %s - %s", url, message) // Prepend "ALERT:"
	smsPayloadStr := fmt.Sprintf("username=%s&password=%s&source=%s&destination=%s&message=%s",
		smsUsername, smsPassword, smsSource, smsDestination, smsMessage)
	smsPayload := strings.NewReader(smsPayloadStr)

	// Create HTTP Request
	req, err := http.NewRequest("POST", smsURL, smsPayload)
	if err != nil {
		fmt.Println("Error creating SMS request:", err)
		return
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Send SMS
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		return
	}
	defer res.Body.Close()

	// Read Response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading SMS response:", err)
		return
	}

	fmt.Println("SMS Response:", string(body))
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
