package lastlog

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	vars "test/VARS"
	"time"
)

func ChangeLastLog(sIP *string, x *vars.DataToInfl, connectedUser *string) {

	file, err := os.Open(vars.LASTLOG_FILE)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", vars.LASTLOG_FILE, err)
		return
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}
	fileSize := fileInfo.Size()

	// Calculate number of entries
	numEntries := int(fileSize / vars.LINE_LENGTH_LASTLOG)

	// Read each entry
	for i := 0; i < numEntries; i++ {
		var ll vars.Lastlog

		// Seek to the position of the current entry
		offset := int64(i * vars.LINE_LENGTH_LASTLOG)
		_, err := file.Seek(offset, 0)
		if err != nil {
			fmt.Printf("Error seeking to offset %d: %v\n", offset, err)
			return
		}

		// Read the entry
		err = binary.Read(file, binary.LittleEndian, &ll)
		if err != nil {
			fmt.Printf("Error reading lastlog entry: %v\n", err)
			return
		}

		if ll.LastLoginTime != 0 {
			lastLogin := time.Unix(int64(ll.LastLoginTime), 0)

			u, err := user.LookupId(strconv.Itoa(i))
			if err != nil {
				fmt.Printf("Error getting username for entry %d: %v\n", i, err)
				continue
			}

			fmt.Printf("Username: %s, Last Login: %s\n", u.Name, lastLogin.String())
			if u.Name == *connectedUser {
				fileRepl, err := os.OpenFile(vars.LASTLOG_FILE, os.O_RDWR, 0644)
				if err != nil {
					fmt.Println(err)
				}
				defer fileRepl.Close()
				// x.Device
				changeTimestamp(i, fileRepl, &x.Time.Sec)
				changeDevice(i, fileRepl, &x.Device)
				changeIP(i, fileRepl, &x.AddrV6)
			}

		}
	}

}

// wine             pts/7    192.192.192.192  Παρ Μαρ  5 09:40:00 +0200 2021whoami
func changeIP(offset int, fileRepl *os.File, sIP *[16]byte) {
	var asciiBytes []byte
	strIP := ""
	for _, j := range sIP {
		if j != 0 {
			intValue := int(j)
			strValue := strconv.Itoa(intValue)
			strIP += strValue + "."
		}
	}

	strIP = strings.TrimRight(strIP, ".")
	if len(strIP) == 0 {
		asciiBytes = (*sIP)[:]
	} else {
		for i := 0; i < len(strIP); i++ {
			asciiBytes = append(asciiBytes, (strIP)[i])
		}
	}
	offsetStart := int64(offset) * int64(292)
	_, errfile := fileRepl.WriteAt(asciiBytes, offsetStart+36)
	if errfile != nil {
		fmt.Println(errfile)
	}
}

func changeDevice(i int, fileRepl *os.File, device *[32]byte) {
	data := []byte{}
	for _, d := range device {
		if d != 0 {
			data = append(data, d)
		}
	}

	if len(data) == 0 {
		for i := 0; i <= 31; i++ {
			data = append(data, 0)
		}
	}
	// data = []byte{112, 116, 115, 47, 57}
	offsetStart := int64(i) * int64(292)
	_, errfile := fileRepl.WriteAt(data, offsetStart+4)

	if errfile != nil {
		fmt.Println(errfile)
	}

}

func changeTimestamp(i int, fileRepl *os.File, timeToChange *int32) {
	offsetStart := int64(i) * int64(292)
	fmt.Println(*timeToChange)
	data := make([]byte, 4) // Assuming 32-bit integer (4 bytes)
	binary.LittleEndian.PutUint32(data, uint32((*timeToChange)))
	_, errfile := fileRepl.WriteAt(data, offsetStart)

	if errfile != nil {
		fmt.Println(errfile, "<<<<<<<<<<<<<")
	}
}
