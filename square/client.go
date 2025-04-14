package square

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)


// Print values to confirm they are loaded from .env
func init() {

	SquareAccessToken := os.Getenv("SQUARE_ACCESS_TOKEN")
	SquareLocationID := os.Getenv("SQUARE_LOCATION_ID")
	
	fmt.Println("Loaded Square Access Token:", SquareAccessToken)
	fmt.Println("Loaded Square Location ID:", SquareLocationID)
}

// Helper to call Square API
func CallSquareAPI(method, url string, body interface{}) ([]byte, error) {

	SquareAccessToken := os.Getenv("SQUARE_ACCESS_TOKEN")

	client := &http.Client{}

	// Convert Go struct to JSON
	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Print debug info
	fmt.Println("📡 Calling Square API:", url)

	if len(SquareAccessToken) >= 6 {
		fmt.Println("Authorization Header: Bearer", SquareAccessToken[:6]+"...")
	} else {
		fmt.Println("SQUARE_ACCESS_TOKEN is empty or too short")
	}
	

	req.Header.Add("Authorization", "Bearer "+SquareAccessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	responseBody, _ := ioutil.ReadAll(resp.Body)

	// Print status for debug
	fmt.Println("Square API Response Status:", resp.Status)

	return responseBody, nil
}
