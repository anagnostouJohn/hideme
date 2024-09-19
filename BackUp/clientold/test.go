package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	vars "test/VARS"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh"
)

var confa vars.Config
var isSudo bool
var commands []string

func init() {
	if _, err := toml.DecodeFile("config.toml", &confa); err != nil {
		log.Fatal(err)
	}
}

func main() {

	filePath := "test2" // replace with your file path
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	hexString := hex.EncodeToString(fileContent)
	z := splitString(hexString, 100000)

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: confa.Server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(confa.Server.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote host
	hostPort := fmt.Sprintf("%s:%s", confa.Server.Host, confa.Server.Port)
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
	session.Stdout = io.Discard // This discards the welcome message (MOTD)
	session.Stderr = io.Discard
	defer session.Close()

	// Prepare buffers to capture stdout and stderr
	// var stdoutBuf, stderrBuf bytes.Buffer
	// session.Stdout = &stdoutBuf
	// session.Stderr = &stderrBuf

	// Create a pipe for session stdin
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
			fmt.Println(Result)
			patern := fmt.Sprintf(`\b%s\b.*\bsudo\b|\bsudo\b.*\b%s\b`, confa.Server.User, confa.Server.User)
			re := regexp.MustCompile(patern)
			if re.MatchString(Result) {
				isSudo = true
			}
		}
	}()

	_, err = stdin.Write([]byte("groups $USER\n"))

	time.Sleep(5 * time.Second)
	if err != nil {
		log.Fatalf("Failed to send command: %s", err)
	}

	for _, chunk := range z {
		commands = append(commands, fmt.Sprintf("echo \"%s\" >> skata.txt \n", chunk))
		// fmt.Println(i)
	}
	isSudo = false
	commands = append(commands, "xxd -r -p skata.txt > output_skata \n")
	commands = append(commands, "chmod +x output_skata \n")
	if isSudo {
		commands = append(commands, "echo 1234 | sudo -S nohup sleep 20 > /dev/null 2>&1 &\n")
	} else {
		commands = append(commands, "echo test >> skata2.txt \n")
	}

	for i, command := range commands {
		fmt.Println("Sending Command :", i, "from", len(commands))
		_, err = stdin.Write([]byte(command))
		if err != nil {
			log.Fatalf("Failed to send command: %s", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	stdin.Close()

	// err = session.Wait()
	err = session.Close()
	if err != nil {
		log.Fatalf("Failed to wait for session: %s", err)
	}

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
