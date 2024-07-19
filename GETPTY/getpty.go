package getpty

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	check "test/CHECK"
	vars "test/VARS"
)

func GetConectedData() (vars.ConnectedData, error) {
	// Get the current process ID

	pid := os.Getpid()

	AppPTYName, err := getTerminalName(pid)
	check.Check("Error on Getting Last PTY", err)
	ppid, err := getParentProcessID(pid)

	if err != nil {
		fmt.Printf("Error getting parent process ID: %v\n", err)
		return vars.ConnectedData{}, err
	}

	grandppid, err := getParentProcessID(ppid)

	if err != nil {
		fmt.Printf("Error getting grandparent process ID: %v\n", err)
		return vars.ConnectedData{}, err
	}

	SSHPTYName, err := getTerminalName(grandppid)
	if err != nil {
		fmt.Printf("Error getting terminal name: %v\n", err)
		return vars.ConnectedData{}, err
	}

	originalUser := os.Getenv("SUDO_USER")
	if originalUser == "" {
		fmt.Println("The command was not run using sudo or the SUDO_USER environment variable is not set.")
	} else {
		fmt.Printf("Original User: %s\n", originalUser)
	}
	ip, err := getIP(SSHPTYName)
	check.Check("Error On getting IP from pty", err)
	AppTime, err := GetTimes(AppPTYName)
	check.Check("Error On Getting Time", err)
	SSHTime, err := GetTimes(SSHPTYName)
	check.Check("Error On Getting Time", err)
	sshpid, firstSpownpid, err := GetTerminalSSHFirsConnectionPID(pid)
	check.Check("Error On Getting SSH PID", err)
	ppts := vars.ConnectedData{IP: ip, User: originalUser, AppPTY: AppPTYName, SSHPTY: SSHPTYName, TimeLoginSSH: SSHTime, TimeProgrammStart: AppTime, SSHPID: sshpid, FirstSpownID: firstSpownpid}

	fmt.Printf("Terminal PID: %s Last Terminal PTY %s\n", SSHPTYName, AppPTYName)
	fmt.Println("Process Pid :", pid, "Parent Pid :", ppid, "Grand Parent Pid ", grandppid)
	return ppts, nil
	// Pouse()

}
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

func GetTerminalSSHFirsConnectionPID(pid int) (int, int, error) {
	idToCheck := 0
	for {
		command, err := getCommandName(pid)
		if err != nil {
			fmt.Printf("Error getting command name for PID %d: %v\n", pid, err)
			return 0, 0, err
		}

		fmt.Printf("PID: %d, Command: %s\n", pid, command)

		if strings.Contains(command, "sshd") {
			fmt.Printf("Found SSH process: PID %d\n", pid)
			idToCheck, err = getParentProcessID(pid)
			check.Check("Error on ID last", err)
			break
		}

		ppid, err := getParentProcessID(pid)
		if err != nil {
			fmt.Printf("Error getting parent PID for PID %d: %v\n", pid, err)
			return 0, 0, err
		}

		if ppid == 1 {
			fmt.Println("Reached the init process, stopping search.")
			break
		}

		pid = ppid
	}
	return pid, idToCheck, nil

}

// getParentProcessID returns the parent process ID for the given process ID
func getParentProcessID(pid int) (int, error) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "ppid=")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Convert the output to an integer
	ppidStr := strings.TrimSpace(string(output))
	ppid, err := strconv.Atoi(ppidStr)
	if err != nil {
		return 0, err
	}

	return ppid, nil
}

// getTerminalName returns the terminal name associated with the given process ID
func getTerminalName(pid int) (string, error) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "tty=")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	tty := strings.TrimSpace(string(output))
	return tty, nil
}

func getIP(cd string) (string, error) {
	output, err := exec.Command("who").Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		if strings.Contains(l, cd) {
			splitedField := strings.Fields(l)
			return splitedField[4][1 : len(splitedField[4])-1], nil
		}
	}
	return "", errors.New("NoIP")
}

func GetTimes(pty string) (string, error) {

	output, err := exec.Command("who").Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, l := range lines {
		if strings.Contains(l, pty) {
			splitedField := strings.Fields(l)
			return splitedField[2] + " " + splitedField[3], nil
		}
	}
	return "", errors.New("NoTime")

}

// getTerminalPID returns the PID of the terminal associated with the given terminal name
// func getTerminalPID(tty string) (int, error) {
// 	// Use the ps command to find the PID of the terminal
// 	cmd := exec.Command("ps", "ax", "-o", "pid,tty", "--no-headers")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Parse the output to find the PID associated with the terminal
// 	lines := strings.Split(string(output), "\n")
// 	for _, line := range lines {
// 		fields := strings.Fields(line)
// 		if len(fields) == 2 && fields[1] == tty {
// 			return strconv.Atoi(fields[0])
// 		}
// 	}

// 	return 0, fmt.Errorf("terminal PID not found")
// }

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
