// main.go
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Atarod61/myproject/send" // مسیر صحیح پکیج send
)

var (
	urlsToCheck = []string{
		"https://exampleadddd.com",
		"https://google.com",
	}
	statusMap = make(map[string]string)
	mu        sync.Mutex
)

func checkWebsite(url string) {
	for {
		resp, err := http.Get(url)
		mu.Lock()
		if err != nil || resp.StatusCode != http.StatusOK {
			statusMap[url] = "DOWN"
			send.SendSMS() // فراخوانی از پکیج send
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

func sendAlert(url string) {
	fmt.Printf("ALERT: %s is DOWN!\n", url)
	send.SendSMS() // فراخوانی از پکیج send
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
