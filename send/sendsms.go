package send

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func SendSMS(message string) error {
	apiURL := "https://panel.asanak.com/webservice/v1rest/sendsms"

	// ساختن داده‌های ارسالی به صورت URL encoded
	data := url.Values{}
	data.Set("username", " 09126123391")
	data.Set("password", "TEST110")
	data.Set("source", "982191099135")
	data.Set("destination", "989105856973")
	data.Set("message", " سایت دان شد")

	// ساختن درخواست HTTP
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return err
	}
	defer resp.Body.Close()

	// خواندن پاسخ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return err
	}

	fmt.Println("SMS API Response:", string(body))
	return nil
}
