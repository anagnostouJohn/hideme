package vars

import "sync"

const (
	Empty        = 0x0
	RunLevel     = 0x1
	BootTime     = 0x2
	NewTime      = 0x3
	OldTime      = 0x4
	InitProcess  = 0x5
	LoginProcess = 0x6
	UserProcess  = 0x7
	DeadProcess  = 0x8
	Accounting   = 0x9
)

type ConnectedData struct {
	IP                string
	User              string
	TimeLoginSSH      string
	TimeLoginSSHEpoch int32
	TimeProgrammStart string
	SSHPTY            string
	AppPTY            string
	SessionNumber     string
	SSHPID            int
	FirstSpownID      int
}

type Dessisions struct {
	Dessision []int
}

type TimeVal struct {
	Sec  int32
	Usec int32
}

const (
	LineSize = 32
	NameSize = 32
	HostSize = 256
)

type DataToInfl struct {
	User    string
	Pid     int32
	Time    TimeVal
	TimeEnd TimeVal
	AddrV6  [16]byte
	Device  [LineSize]byte
	ConData ConnectedData
}

type Lastlog struct {
	LastLoginTime int32
	Unused        [256]byte
}

type ExitStatus struct {
	Termination int16
	Exit        int16
}

type Utmp struct {
	Type   int16
	Pad1   int16
	Pid    int32
	Device [LineSize]byte
	ID     [4]byte
	User   [NameSize]byte
	Host   [HostSize]byte
	Exit   struct {
		Termination int16
		Exit        int16
	}
	Session int32
	Time    struct {
		Sec  int32
		Usec int32
	}
	AddrV6   [16]byte
	Reserved [20]byte
}

type Connection struct {
	Host     string
	Port     string
	Username string
	Password string
	IsUsed   bool
	Place    string
}

type AllConnections struct {
	Mu   sync.Mutex
	Conn []Connection
}

type DelaConnection struct {
	Single bool
	Conn   Connection
}

// //////////UTMP///////////
var UTMP_FILE = "/var/run/utmp"
var UTMP_SIZE = 384

// //////////LASTLOG///////////
const LINE_LENGTH_LASTLOG = 292
const LASTLOG_FILE = "/var/log/lastlog"

// //////////AUTH_LOG ///////////
var AUTH_LOG string = "/var/log/auth.log"

var WTMP string = "/var/log/wtmp"

// var AUTH_LOG string = "/var/log/auth.log"
var SYSLOG string = "/var/log/syslog"

///////////////////////////////////////

var Sessions = "/run/systemd/sessions"

var BrFileHomeDir string

type Config struct {
	Client ClientConfig `toml:"server"`
	Flags  AppFlags     `toml:"flags"`
}

type ClientConfig struct {
	User string
	Port string
	Host string
	Pass string
}

type AppFlags struct {
	PreFile       string
	MainFile      string
	Destr         bool
	Combo         bool
	BrFile        string
	Threads       int
	Key           string
	RundomTimeSec int
	Pids          string
	SessionId     string
}

// ANSI escape code for background color
// Format: "\033[<code>m"
// Background color codes:
// 40: Black
// 41: Red
// 42: Green
// 43: Yellow
// 44: Blue
// 45: Magenta
// 46: Cyan
// 47: White

// Set background to blue (44)
var Blue = "\033[44m"
var Yellow = "\033[43m"
var Red = "\033[41m"

var Green = "\033[42m"

// // Print some text with the blue background
// fmt.Println("This text has a blue background!")

// // Reset to default terminal colors
var Reset = "\033[0m"

// ANSI escape codes for text (foreground) colors
// Format: "\033[<code>m"
// Text color codes:
// 30: Black
// 31: Red
// 32: Green
// 33: Yellow
// 34: Blue
// 35: Magenta
// 36: Cyan
// 37: White

var GreenString = "\033[32m"

// Set text color to red (31)
// fmt.Print("\033[31m")

// // Print text with red color
// fmt.Println("This text is red!")

// // Set text color to green (32)
// fmt.Print("\033[32m")
// fmt.Println("This text is green!")

// // Reset to default terminal colors
// fmt.Print("\033[0m")
