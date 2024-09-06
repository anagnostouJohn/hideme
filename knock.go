package main

import (
	"fmt"
	"net"
	"time"
)

// knock sends a TCP knock to a specific port on a host
func knock(host string, port int, timeout time.Duration) error {
	// Create a TCP address from the host and port
	address := fmt.Sprintf("%s:%d", host, port)

	// Dial the address with a specific timeout
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}

	// Close the connection immediately
	conn.Close()
	return nil
}

// sendKnocks sends a sequence of knocks to a host
func sendKnocks(host string, ports []int, delay time.Duration) {
	for _, port := range ports {
		err := knock(host, port, 5*time.Second)
		if err != nil {
			fmt.Printf("Failed to knock on port %d: %s\n", port, err)
		} else {
			fmt.Printf("Knocked on port %d successfully\n", port)
		}

		// Wait for the specified delay before the next knock
		time.Sleep(delay)
	}
}

func main() {

	// The host to knock on
	host := "192.168.23.61" // Replace with the actual IP or hostname

	// The sequence of ports to knock
	ports := []int{7666}

	// Delay between knocks
	delay := 1 * time.Second

	// Send the knocks
	sendKnocks(host, ports, delay)
	// select {}
}
