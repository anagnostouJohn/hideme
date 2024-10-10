package checkactive

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
	check "test/client/CHECK"
	"test/vars"
	"time"

	"golang.org/x/crypto/ssh"
)

var wg sync.WaitGroup
var found = true
var cur int64 = 0
var prevCur int64 = 0

func Checkactive(conf vars.Config) {
	config := &ssh.ClientConfig{
		User: conf.Client.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Client.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote host
	hostPort := fmt.Sprintf("%s:%s", conf.Client.Host, conf.Client.Port)
	client, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()
	//TODO execute that command stat -c %c c
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	fmt.Println("sadasdasdasd")
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to create stdin pipe: %s", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %s", err)
	}

	wg.Add(1)
	go CheckStds(stdout)
	defer stdin.Close()

	// Start the session shell
	// err = session.Shell()
	if err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}
	for {
		if found {
			fmt.Println("Check", found)
			err = session.Run("stat -c %x /tmp/c")
			if err != nil {
				check.Check("Error Sending Command", err)
			}
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("END")
			break
		}

	}
	wg.Wait()

}

func CheckStds(stdout io.Reader) {
	defer wg.Done()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		Result := scanner.Text()
		fmt.Println(Result)
		// timeStr := "2024-10-07 10:14:26.210630517 +0300"

		// Define the layout to match the input string format
		layout := "2006-01-02 15:04:05.999999999 -0700"

		// Parse the time string into a time.Time object
		parsedTime, err := time.Parse(layout, Result)
		if err != nil {
			log.Fatalf("Error parsing time: %s", err)
		}

		// Convert to Unix epoch time (seconds since Jan 1, 1970)
		epochTime := parsedTime.Unix()

		// Print the Unix epoch time
		cur = epochTime
		fmt.Println(cur, prevCur)
		if cur < prevCur {
			fmt.Println("NOT OK")
			prevCur = cur
		} else if cur == prevCur {
			fmt.Println("END RETREAVE")
			found = false

		}
	}
}
