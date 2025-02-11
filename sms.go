package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// متغیرهای پیکربندی (بهتر است از متغیرهای محیطی خوانده شوند)
var (
	accountSid    = os.Getenv("TWILIO_ACCOUNT_SID")
	authToken     = os.Getenv("TWILIO_AUTH_TOKEN")
	twilioNumber  = os.Getenv("TWILIO_PHONE_NUMBER")
	adminNumber   = os.Getenv("ADMIN_PHONE_NUMBER")
	systemURL     = "https://dl.jzac.ir/.com" // آدرس سامانه خود را اینجا قرار دهید
	checkInterval = 5 * time.Minute           // فاصله زمانی بین بررسی‌ها
)

func main() {
	fmt.Println("Monitoring system started...")

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := checkSystemStatus(systemURL); err != nil {
			fmt.Printf("System is down: %v\n", err)
			if err := sendSMS(fmt.Sprintf("System %s is DOWN!", systemURL)); err != nil {
				fmt.Printf("Failed to send SMS: %v\n", err)
			}
		} else {
			fmt.Println("System is UP.")
		}
	}
}

func checkSystemStatus(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	return nil
}

func sendSMS(message string) error {
	twilioURL := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	data := url.Values{}
	data.Set("To", adminNumber)
	data.Set("From", twilioNumber)
	data.Set("Body", message)

	req, err := http.NewRequest("POST", twilioURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("SMS sent successfully.")
		return nil
	}

	return fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
}
