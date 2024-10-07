package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {

	// f, err := os.OpenFile("c", os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	fmt.Println("errr", err)
	// }

	// f.Close()
	// os.ReadFile("c")
	// fa.Close()
	fileInfo, err := os.Stat("c")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get the underlying data of the file, specific to Unix
	stat := fileInfo.Sys().(*syscall.Stat_t)

	// Convert access time from seconds and nanoseconds to time.Time
	atime := time.Unix(stat.Atim.Sec, stat.Atim.Nsec)

	// Print the last access time
	fmt.Printf("Last access time: %v\n", atime)
}

// 	// Webhook URL (replace with your actual webhook URL)
// 	webhookURL := "https://discordapp.com/api/webhooks/1291024909111267483/D-Q_Q0A6n3orgR4j5PYYo8QfNENspDy9AWa9o_szwWzE_tfNUzMSetUIDVHIhMyvJ_2-"

// 	// Data to send to the webhook
// 	payload := map[string]string{
// 		"content": "Hello, this is a test message!",
// 	}

// 	// Convert the payload to JSON
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Fatalf("Error marshaling payload: %v", err)
// 	}

// 	// Create a new POST request with the JSON payload
// 	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		log.Fatalf("Error creating POST request: %v", err)
// 	}

// 	// Set the content-type header to application/json
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the request using http.DefaultClient
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatalf("Error sending POST request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status
// 	// if resp.StatusCode != http.StatusOK {
// 	// 	if resp.StatusCode != 204 {
// 	// 		log.Fatalf("Received non-OK response: %s", resp.Status)
// 	// 	}

// 	// }

// 	fmt.Println("POST request successful!")
// }
