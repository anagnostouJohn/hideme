package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var ZeroPort = make(chan bool)
var OnePort = make(chan bool)
var LivePort = make(chan bool)
var finalIncome = ""

// chan OnePort bool
// PortMonitor monitors a list of ports for incoming connection attempts

func main() {
	// Define the ports to monitor
	ports := []int{7666, 8666, 6666}
	// Define the timeout duration
	timeout := 20 * time.Second
	wg.Add(1)
	go CollectPorts(ZeroPort, OnePort, LivePort)
	// lisen := "tcp"
	// if lisen == "tcp" {
	for i, j := range ports {
		wg.Add(1)
		// Start monitoring the ports
		if i == 0 {
			go listenTcpPort(j, timeout, ZeroPort)
		} else if i == 1 {
			go listenTcpPort(j, timeout, OnePort)
		} else if i == 2 {
			fmt.Println("STARTA MAMAMAMAMA")
			go listenTcpPort(j, timeout, LivePort)
		}

	}
	// } else {

	// }
	wg.Wait()
	// Block forever (or you can use a different way to manage program lifecycle)
	// select {}
}

func listenUdpPort(port int, timeout time.Duration, PortSelect chan bool) {
	address := ":8080" // Change this to the address and port you want to listen to
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening on", address)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buf[:n]))
	}
}

func listenTcpPort(port int, timeout time.Duration, PortSelect chan bool) {
	// Format the address
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Error listening on port %d: %s\n", port, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Monitoring port %d...\n", port)
	for {
		// Set a timeout for the listener
		listener.(*net.TCPListener).SetDeadline(time.Now().Add(timeout))
		// Accept the connection
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				fmt.Println("END")
				//
				continue // Timeout occurred, continue listening
			}
			fmt.Printf("Error accepting connection on port %d: %s\n", port, err)
			return
		}
		// Log the connection attempt
		fmt.Printf("Connection attempt on port %d from %s\n", port, conn.RemoteAddr())
		PortSelect <- true
		conn.Close()
	}
}

func CollectPorts(ZPort, OPort, LPort chan bool) {
	fmt.Println("HELLO THERE")
	for {
		select {
		case <-ZPort:
			CollectBytes(false)
			// fmt.Println(Zero, "ZERO")
		case <-OPort:
			CollectBytes(true)
			// fmt.Println(One, "ONE")
		case <-LPort:
			fmt.Println("ALIVE")
		default:
			time.Sleep(100 * time.Millisecond)
		}

	}
	// wg.Done()
}

func CollectBytes(zeroOrOne bool) {
	// fmt.Println(len(finalIncome))
	if zeroOrOne {
		finalIncome += "1"
	} else {
		finalIncome += "0"
	}
	// fmt.Println(finalIncome)
	// finalIncome = "test"
	if len(finalIncome) == 8 {
		chatString := []int{}
		for _, char := range finalIncome {

			c := fmt.Sprintf("%c", char)
			num, err := strconv.Atoi(c)

			if err != nil {
				fmt.Println("Conversion error:", err)
			}
			chatString = append(chatString, num)

		}
		x := bitsToString(chatString)
		fmt.Println("->>>>>>>>>>> ", x, " <<<<<<<<<<<<<")
		finalIncome = ""
	}

}

func bitsToString(bits []int) string {
	if len(bits)%8 != 0 {
		panic("The number of bits is not a multiple of 8")
	}

	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bits); i += 8 {
		byteStr := ""
		for j := 0; j < 8; j++ {
			byteStr += strconv.Itoa(bits[i+j])
		}
		b, _ := strconv.ParseUint(byteStr, 2, 8)
		bytes[i/8] = byte(b)
	}

	return string(bytes)
}
