package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	vars "test/VARS"
	"time"

	"github.com/BurntSushi/toml"
)

var wg sync.WaitGroup
var ZeroPort = make(chan bool)
var OnePort = make(chan bool)
var LivePort = make(chan bool)
var finalIncome = ""
var conf vars.Config
var records [][]string
var result = ""

// chan OnePort bool
// PortMonitor monitors a list of ports for incoming connection attempts

func init() {

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	if false { //<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<,//TODO Remove False
		os.Remove("config.toml")
	}

	file, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// Read all the records from the CSV
	records, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	timeout := 20 * time.Second
	wg.Add(1)
	go CollectPorts(ZeroPort, OnePort, LivePort)

	wg.Add(3)

	go listenTcpPort(conf.Flags.KnockData[0], timeout, ZeroPort)
	go listenTcpPort(conf.Flags.KnockData[1], timeout, OnePort)
	go listenTcpPort(conf.Flags.KnockAlive, timeout, LivePort)
	wg.Wait()
	// Block forever (or you can use a different way to manage program lifecycle)
	// select {}
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
		// fmt.Println("->>>>>>>>>>> ", x, " <<<<<<<<<<<<<")
		result += x
		fmt.Println(result, "<<<<<")
		if strings.HasSuffix(result, "!") {
			str := strings.TrimSuffix(result, "!")
			splResult := strings.Split(str, "-")
			ip, err := strconv.Atoi(splResult[0])
			CheckError("ip Place Convert", err)
			user, err := strconv.Atoi(splResult[1])
			CheckError("user Place Convert", err)
			pass, err := strconv.Atoi(splResult[2])
			CheckError("pass Place Convert", err)

			port, err := strconv.Atoi(splResult[3])
			CheckError("port Place Convert", err)

			fmt.Println(records[ip][0], records[user][1], records[pass][2], records[port][3])
			result = ""
		}

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

func CheckError(msg string, err error) error {

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
