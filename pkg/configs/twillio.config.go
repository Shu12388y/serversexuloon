package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func SendWhatsAppMessage(toNumber string, contentVariable string) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	accountSID := os.Getenv("ACCOUNTSID")
	authToken := os.Getenv("AUTHTOKEN")
	fromNumber := os.Getenv("FROMNUMBER")
	contentSID := os.Getenv("CONTENTSID")

	// Prepare the form data
	form := url.Values{}
	form.Add("To", "whatsapp:"+toNumber)
	form.Add("From", fromNumber)
	form.Add("ContentSid", contentSID)
	form.Add("ContentVariables", fmt.Sprintf(`{"1":"%s"}`, contentVariable))

	// Create the HTTP request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID),
		strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Response:", string(body))

	return nil
}


func SendWhatsAppMessageLogin(toNumber string, contentVariable string) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	accountSID := os.Getenv("ACCOUNTSID")
	authToken := os.Getenv("AUTHTOKEN")
	fromNumber := os.Getenv("FROMNUMBER")
	contentSID := os.Getenv("CONTENTSID")

	// Prepare the form data
	form := url.Values{}
	form.Add("To", "whatsapp:"+toNumber)
	form.Add("From", fromNumber)
	form.Add("ContentSid", contentSID)
	form.Add("ContentVariables", fmt.Sprintf(`{"1":"%s"}`, contentVariable))

	// Create the HTTP request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID),
		strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// Set headers
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Response:", string(body))

	return nil
}

