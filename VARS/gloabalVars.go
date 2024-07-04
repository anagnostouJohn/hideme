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
}

type Lastlog struct {
	LastLoginTime int32
	Unused        [256]byte
}

type ExitStatus struct {
	Termination int16
	Exit        int16
}

type DataLogin struct {
	Username string
	Datetime string
	Ip       string
	PTY      string
}

// type Utmp struct {
// 	Type int16
// 	// alignment
// 	_       [2]byte
// 	Pid     int32
// 	Device  [LineSize]byte
// 	Id      [4]byte
// 	User    [NameSize]byte
// 	Host    [HostSize]byte
// 	Exit    ExitStatus
// 	Session int32
// 	Time    vars.TimeVal
// 	AddrV6  [16]byte
// 	// Reserved member
// 	Reserved [20]byte
// }

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

var (
	BrFile        string
	BrFileHomeDir string
	Host          string
	Port          string
	User          string
	Pass          string
	ConnectedUser string
	Threads       int
	HideMe        bool
	Combo         bool
	Destr         bool
)
