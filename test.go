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
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	exeName := filepath.Base(exePath)
	fmt.Println("Executable Name:", exeName)
	time.Sleep(10 * time.Second)
	fmt.Println("Ending")
}
