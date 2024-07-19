package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Function to get the parent process ID of a given PID
func getParentProcessID(pid int) (int, error) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "ppid=")
	fmt.Println("ppid ->", cmd)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	ppidStr := strings.TrimSpace(string(output))
	ppid, err := strconv.Atoi(ppidStr)
	if err != nil {
		return 0, err
	}

	return ppid, nil
}

// Function to get the command name of a given PID
func getCommandName(pid int) (string, error) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	fmt.Println(cmd)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	command := strings.TrimSpace(string(output))
	return command, nil
}

func main() {
	pid := os.Getpid() // Replace with the actual PID of the `test` process

	for {
		command, err := getCommandName(pid)
		if err != nil {
			fmt.Printf("Error getting command name for PID %d: %v\n", pid, err)
			return
		}

		fmt.Printf("PID: %d, Command: %s\n", pid, command)

		if strings.Contains(command, "sshd") {
			fmt.Printf("Found SSH process: PID %d\n", pid)
			break
		}

		ppid, err := getParentProcessID(pid)
		if err != nil {
			fmt.Printf("Error getting parent PID for PID %d: %v\n", pid, err)
			return
		}

		if ppid == 1 {
			fmt.Println("Reached the init process, stopping search.")
			break
		}

		pid = ppid
	}
	time.Sleep(20 * time.Second)
}
