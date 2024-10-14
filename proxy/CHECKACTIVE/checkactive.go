package checkactive

import (
	"fmt"
	"log"
	"strings"
	check "test/client/CHECK"
	"test/vars"
	"time"

	"golang.org/x/crypto/ssh"
)

var found = true
var prevCur int64 = 0

func Checkactive(conf vars.Config) {
	for { //BUG MAKE it witrh stdin write i supose
		if found {
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
			session, err := client.NewSession()
			if err != nil {
				log.Fatalf("Failed to create session: %s", err)
			}

			fmt.Println("Check", found)
			output, err := session.Output("stat -c %x /tmp/c")
			if err != nil {
				check.Check("Error Sending Command", err)
			}
			timestamp := strings.TrimSpace(string(output))

			layout := "2006-01-02 15:04:05.999999999 -0700"
			parsedTime, err := time.Parse(layout, timestamp)
			if err != nil {
				fmt.Println(err, "Error parsing time")
			}
			epochTime := parsedTime.Unix()

			if prevCur < epochTime {
				fmt.Println("NOT OK")
				prevCur = epochTime
			} else if epochTime == prevCur {
				fmt.Println("END RETRIEVE")
				found = false
				session.Close()
				client.Close()

				session, err := client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}

				output, err := session.Output("cat /tmp/c")
				if err != nil {
					fmt.Println(err, "Error retrieving file contents")
				}
				fmt.Println(string(output))
				session.Close()
				client.Close()

				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}

				_, err = session.Output("rm /tmp/c")
				if err != nil {
					fmt.Println(err, "Error removing file")
				}
				session.Close()
				client.Close()
				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}
				_, err = session.Output("sed -i '1d' filename")
				if err != nil {
					fmt.Println(err, "Error editing file")
				}
				session.Close()
				client.Close()
				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}
				_, err = session.Output("bash -c 'echo %s | sudo systemctl restart rsyslog ' \n ")
				if err != nil {
					fmt.Println(err, "Error restarting rsyslog")
				}
			}
			session.Close()
			client.Close()
			time.Sleep(10 * time.Second)
		} else {
			fmt.Println("END")
			break
		}
	}
}
