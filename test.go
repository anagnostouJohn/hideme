package main

import (
	"fmt"
	"os"
)

func main() {
	// Get the original username from the SUDO_USER environment variable
	originalUser := os.Getenv("SUDO_USER")
	if originalUser == "" {
		fmt.Println("The command was not run using sudo or the SUDO_USER environment variable is not set.")
	} else {
		fmt.Printf("Original User: %s\n", originalUser)
	}

	// For additional information, let's also print the current user and effective user ID
	currentUser := os.Getenv("USER")
	fmt.Printf("Current User: %s\n", currentUser)

	euid := os.Geteuid()
	fmt.Printf("Effective User ID: %d\n", euid)
}
