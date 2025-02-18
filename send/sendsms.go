package send

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendSMS() {
	url := "https://panel.asanak.com/webservice/v1rest/sendsms"
	str := "username=xxxxxxxx&password=xxxxxx" +
		"source=98xxxxxxxx&destination=98xxxxxxx&message=سایت خراب شد،سلام" + url
	payload := strings.NewReader(str)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	// Replacing ioutil.ReadAll with io.ReadAll
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(res)
	fmt.Println(string(body))
}
