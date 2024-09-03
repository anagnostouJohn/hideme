// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// // Function to get the parent process ID
// func getParentProcessID(pid int) (int, error) {
// 	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "ppid=")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return 0, err
// 	}

// 	ppidStr := strings.TrimSpace(string(output))
// 	ppid, err := strconv.Atoi(ppidStr)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ppid, nil
// }

// // Function to get the terminal name associated with a process ID
// func getTerminalName(pid int) (string, error) {
// 	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "tty=")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", err
// 	}
// 	tty := strings.TrimSpace(string(output))
// 	return tty, nil
// }

// func main() {
// 	// Get the current process ID
// 	pid := os.Getpid()

// 	// Get the terminal name for the current process
// 	AppPTYName, err := getTerminalName(pid)
// 	if err != nil {
// 		fmt.Printf("Error getting terminal name for PID %d: %v\n", pid, err)
// 		return
// 	}

// 	// Get the parent process ID
// 	ppid, err := getParentProcessID(pid)
// 	if err != nil {
// 		fmt.Printf("Error getting parent process ID for PID %d: %v\n", pid, err)
// 		return
// 	}

// 	// Get the grandparent process ID
// 	grandppid, err := getParentProcessID(ppid)
// 	if err != nil {
// 		fmt.Printf("Error getting grandparent process ID for PID %d: %v\n", ppid, err)
// 		return
// 	}

// 	// Get the terminal name for the grandparent process
// 	SSHPTYName, err := getTerminalName(grandppid)
// 	if err != nil {
// 		fmt.Printf("Error getting terminal name for grandparent PID %d: %v\n", grandppid, err)
// 		return
// 	}

// 	// Get the original user if running with sudo
// 	originalUser := os.Getenv("SUDO_USER")
// 	if originalUser == "" {
// 		fmt.Println("The command was not run using sudo or the SUDO_USER environment variable is not set.")
// 	} else {
// 		fmt.Printf("Original User: %s\n", originalUser)
// 	}

//		fmt.Printf("Terminal PID: %s Last Terminal PTY %s\n", SSHPTYName, AppPTYName)
//		fmt.Println("Process Pid :", pid, "Parent Pid :", ppid, "Grand Parent Pid ", grandppid)
//		time.Sleep(10 * time.Second)
//	}
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

var b bytes.Buffer

func main() {

	session, err := CreateConn()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Create a pipe to send commands to the session
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %v", err)
	}

	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	// Start the shell
	if err := session.Shell(); err != nil {
		log.Fatalf("Failed to start shell: %v", err)
	}

	// Feed commands to the session
	commands := []string{
		"echo test >> test.txt",
		"echo my test >> test.txt",
		"echo testme >> test.txt",
		"cat test.txt", // This will allow us to see the contents of test.txt in the buffer
	}

	for _, cmd := range commands {
		fmt.Fprintln(stdin, cmd)
	}

	// Optionally, close the stdin pipe when done
	stdin.Close()

	// Wait for the session to finish
	if err := session.Wait(); err != nil && err != io.EOF {
		log.Fatalf("Failed to wait for session: %v", err)
	}

	// Print the output captured in the buffer
	fmt.Println(b.String())

	// err = s.Run("echo test >> test.txt")
	// fmt.Println(err, "1")
	// err = s.Run("echo my test >> test.txt")
	// fmt.Println(err, "2")
	// f := fmt.Sprintf("echo %s >> test.txt", "testme")
	// err = s.Run(f)
	// fmt.Println(err, "3")
	// s.Run("echo ")

	// fmt.Println(b.String())

}

func CreateConn() (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: "wine",
		Auth: []ssh.AuthMethod{
			ssh.Password("1234"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	hostPort := fmt.Sprintf("%s:%s", "192.168.23.89", "22")
	client, err := ssh.Dial("tcp", hostPort, config)

	if err != nil {

		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
