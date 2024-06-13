package main

import (
	"fmt"
	"regexp"
)

func main() {
	strings := []string{
		"string1",
		"string2",
		"string3",
		"string4",
		"string5",
		"string6",
		"string7",
		"string8",
		"string9",
		"string10",
	}

	for i, j := range strings[4:] {
		fmt.Println(i+4, "---", j)
	}

	systemdLogInd := "7"
	numberToCheck := "897"

	// Create the regex pattern with the systemdLogInd and numberToCheck values
	patternSystemLogInd := fmt.Sprintf(`^.*systemd-logind\[%s\].*(Session %s logged out|Removed session %s|New session %s)`, numberToCheck, systemdLogInd, systemdLogInd, systemdLogInd)

	// Compile the regex
	re := regexp.MustCompile(patternSystemLogInd)

	lines := []string{
		"Jun 12 09:27:50 ubuntu systemd-logind[897]: Session 7 logged out. Waiting for processes to exit. KKKKKKKKKKK 897",
		"Jun 12 09:27:50 ubuntu systemd-logind[897]: Removed session 7. KKKKKKKKKKK 897",
		"Some other log line",
		"Jun 12 09:27:50 ubuntu systemd-logind[1234]: New session 7",
		"Jun 12 09:27:50 ubuntu systemd-logind[897]: Session 7 logged out. Waiting for processes to exit.",
	}

	for _, line := range lines {
		if re.MatchString(line) {
			fmt.Println("Match found in line:", line)
		} else {
			fmt.Println("No match in line:", line)
		}
	}

	// Open the file for writing, create it if it doesn't exist, truncate it if it does
	// file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()

	// for _, j := range strings {
	// 	_, err := file.WriteString(j + "\n")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }

	fmt.Println("All strings written to file successfully")
}
