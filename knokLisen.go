package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

// PortMonitor monitors a list of ports for incoming connection attempts

func main() {
	// Define the ports to monitor
	// ports := []int{7666}
	// Define the timeout duration
	timeout := 20 * time.Second
	wg.Add(1)
	// Start monitoring the ports
	go listenPort(7666, timeout)
	wg.Wait()
	// Block forever (or you can use a different way to manage program lifecycle)
	// select {}
}

// func PortMonitor(ports []int, timeout time.Duration) {
// 	for _, port := range ports {
// 		go listenPort(port, timeout)
// 	}

// }

// listenPort listens for incoming connection attempts on a specific port
func listenPort(port int, timeout time.Duration) {
	// Format the address
	for {
		address := fmt.Sprintf("0.0.0.0:%d", port)
		// Start listening on the port
		listener, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Printf("Error listening on port %d: %s\n", port, err)
			return
		}
		// defer

		fmt.Printf("Monitoring port %d...\n", port)

		// Set a timeout for the listener
		listener.(*net.TCPListener).SetDeadline(time.Now().Add(timeout))

		// Accept the connection
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				fmt.Println("END")
				listener.Close()
				time.Sleep(10 * time.Second)
				continue // Timeout occurred, continue listening
			}
			fmt.Printf("Error accepting connection on port %d: %s\n", port, err)
			return
		}

		// Log the connection attempt
		fmt.Printf("Connection attempt on port %d from %s\n", port, conn.RemoteAddr())

		// Close the connection immediately
		listener.Close()
		conn.Close()
	}
}
