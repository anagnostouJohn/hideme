package knock

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var ports = []int{7666, 8666, 6666}

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
func sendData(host string, port int, delay time.Duration) {
	// for _, port := range ports {
	err := knock(host, ports[port], 5*time.Second)
	if err != nil {
		fmt.Printf("Failed to knock on port %d: %s\n", port, err)
	} else {
		fmt.Printf("Knocked on port %d successfully\n", port)
	}

	// Wait for the specified delay before the next knock
	time.Sleep(delay)
	// }
}

func SendIAmAlive(host string, port int, delay time.Duration) {
	for {
		fmt.Println("Send i am alive", host, port, delay)
		knock(host, ports[port], delay)
		time.Sleep(2 * time.Second)
	}

}

func SendKnock(data, host string) {
	data = "1-833-2END" //TODO remove
	// The host to knock on
	host = "192.168.23.61" // TODO Replace with the actual IP or hostname

	// The sequence of ports to knock
	// ports := []int{7666, 8666}

	// Delay between knocks
	delay := 500 * time.Millisecond
	x := []string{}
	// go SendIAmAlive(host, 2, delay)

	// Convert string to bytes
	byteArray := []byte(name)

	fmt.Println("Bit representation of 'John':")

	// Print each byte as bits
	for _, b := range byteArray {
		x = append(x, byteToBits(b))
		// fmt.Printf("Character '%c' -> Bits: %s\n", name[i], byteToBits(b))
	}
	fmt.Println(x)
	for _, j := range x {
		for _, k := range j {
			// char := rune(k)
			char := fmt.Sprintf("%c", k)
			// fmt.Println(char)
			intValue, err := strconv.Atoi(char)
			if err != nil {
				fmt.Println(err)
			}
			sendData(host, intValue, delay)
		}
	}

	// Send the knocks

	// select {}
}

func byteToBits(b byte) string {
	var bits string
	for i := 7; i >= 0; i-- { // Iterate through each bit from MSB to LSB
		if b&(1<<i) != 0 {
			bits += "1"
		} else {
			bits += "0"
		}
	}
	return bits
}
