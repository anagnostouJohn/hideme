package main

import (
	"bufio"
	"fmt"
	"os"
)

var AUTH_LOG string = "/var/log/auth.log"

func main() {

	lineNumber := 3 // Example: delete the 3rd line

	err := deleteLine(AUTH_LOG, lineNumber)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Line deleted successfully")
	}
}

func deleteLine(filePath string, lineNumber int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	// defer file.Close()

	tempFile, err := os.Create("/tmp/aaaattt.txt")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	currentLine := 0
	for scanner.Scan() {
		currentLine++
		if currentLine == lineNumber {
			continue // Skip the line to be deleted
		}
		x, err := tempFile.WriteString(scanner.Text() + "\n")
		fmt.Println(x)
		if err != nil {
			return err
		}
	}

	input, err := os.ReadFile(tempFile.Name())
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(input)
	err = os.WriteFile("/home/ubuntu/go/src/ttt.txt", input, 0644)
	if err != nil {
		fmt.Println("Error creating")
		fmt.Println(err)

	}
	// os.Remove(tempFile.Name())

	// if err := scanner.Err(); err != nil {
	// 	return err
	// }

	// tempFileName := tempFile.Name()
	// fmt.Println(tempFileName)
	// file.Close()
	// tempFile.Close()

	// err = os.Rename(tempFileName, filePath)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// package main

// import (
// 	"fmt"
// 	"strings"
// )

// func main() {

// 	x := []byte{112, 116, 115, 47, 50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	// z :=
// 	// fmt.Println(z, "<<<<<<<<<<<<<<<<")
// 	if strings.Contains(strings.TrimRight(string(x), "\x00"), "pts") {
// 		fmt.Println("OK")
// 	}

// }

// func test(ip *string) {
// 	var asciiBytes []byte

// 	for i := 0; i < len(*ip); i++ {
// 		// fmt.Println(ip[i])
// 		asciiBytes = append(asciiBytes, (*ip)[i])
// 	}

// 	// Print the ASCII byte values
// 	// for _, b := range asciiBytes {
// 	// 	fmt.Printf("%d ", b)
// 	// }
// 	fmt.Println(asciiBytes)
// }

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	ip := "192.168.23.23"
// 	test(&ip)
// }

// func test(ip *string) {
// 	var asciiBytes []byte

// 	for i := 0; i < len(*ip); i++ {
// 		// fmt.Println(ip[i])
// 		asciiBytes = append(asciiBytes, (*ip)[i])
// 	}

// 	// Print the ASCII byte values
// 	// for _, b := range asciiBytes {
// 	// 	fmt.Printf("%d ", b)
// 	// }
// 	fmt.Println(asciiBytes)
// }

// package main

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// )

// func main() {
// 	var number int = 192

// 	// To handle both 32-bit and 64-bit architectures, you can use int32 or int64
// 	var number8 int8 = int8(number) // Convert to int64 for this example

// 	// Create a buffer
// 	buffer := new(bytes.Buffer)

// 	// Write the int64 number to the buffer in BigEndian order
// 	err := binary.Write(buffer, binary.BigEndian, number8)
// 	if err != nil {
// 		fmt.Println("binary.Write failed:", err)
// 	}

// 	// Get the final byte array
// 	finalByteArray := buffer.Bytes()
// 	fmt.Printf("Final byte array: %v\n", finalByteArray)
// }

// package main

// import (
// 	"fmt"
// 	"reflect"
// )

// func main() {
// 	slice1 := []byte{192, 168, 23, 23, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
// 	slice2 := []byte{192, 168, 23, 23, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// 	if reflect.DeepEqual(slice1, slice2) {
// 		fmt.Println("The slices are equal.")
// 	} else {
// 		fmt.Println("The slices are not equal.")
// 	}
// }

// package main

// import "fmt"

// func main() {
// 	// Define the initial part of the byte array
// 	initialBytes := []byte{192, 168, 23, 23}

// 	// Define the number of additional zero bytes needed
// 	zeroBytes := make([]byte, 12)

// 	// Append the zero bytes to the initial byte array
// 	byteArray := append(initialBytes, zeroBytes...)

// 	// Print the byte array
// 	fmt.Println(byteArray)
// }

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	filePath := "example.txt"
// 	start := int64(5)          // Start deleting from the 6th byte (zero-based index)
// 	bytesToDelete := int64(10) // Number of bytes to delete

// 	err := deleteBytesFromFile(filePath, start, bytesToDelete)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	fmt.Println("Bytes deleted successfully.")
// }

// var network bytes.Buffer
// enc := gob.NewEncoder(&network)
// enc.Encode(P{3, 4, 5, "Pythagoras"})

// fmt.Println(network.Bytes())
// file, err := os.OpenFile(LASTLOG_FILEa, os.O_RDWR, 0644)
// if err != nil {
// 	fmt.Println(err)
// }

// defer file.Close()

// for i := int64(1); i <= 20; i++ {
// 	_, err := file.WriteAt([]byte("\x00"), i)
// 	if err != nil {
// 		log.Fatalf("failed writing to file: %s", err)
// 	}
// }
// // fmt.Printf("\nLength: %d bytes", len)
// fmt.Printf("\nFile Name: %s", file.Name())

// filePath := "hi.txt" // Path to the file whose modification time you want to change

// // Get the current time
// currentTime := time.Now()

// // Define the desired modification time (replace with your desired time)
// desiredModTime := time.Date(2022, time.January, 1, 12, 34, 56, 0, time.UTC)

// // Change the modification time of the file
// err := os.Chtimes(filePath, currentTime, desiredModTime)
// if err != nil {
// 	fmt.Printf("Error changing modification time: %v\n", err)
// 	return
// }

// fmt.Println("Modification time changed successfully")
// }
