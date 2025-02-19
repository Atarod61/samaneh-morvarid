package send

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendSMS() error {
	url := "https://panel.asanak.com/webservice/v1rest/sendsms"
	str := "username=xxxxxxxx&password=xxxxxx" +
		"source=98xxxxxxxxxx&destination=98xxxxxxxx&message=سایت خراب شد،سلام" //+ url
	payload := strings.NewReader(str)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %w", err)
	}
	defer res.Body.Close()

	// Replacing ioutil.ReadAll with io.ReadAll
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}

	fmt.Println(res)
	fmt.Println(string(body))
	return nil // در صورت موفقیت، مقدار nil برگردانده می‌شود
}
