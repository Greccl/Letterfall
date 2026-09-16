package main

import (
	"os"
	"bufio"
	"net"
	"fmt"
	// "errors"
	"path/filepath"
	"github.com/spf13/pflag"
)

var lucentHead, lucentBody Color
var backHead, backNeck, backTail Color
var backCharset int


var normalHead, normalNeck, normalTail Color
var normalMinLen, normalMaxLen int
var normalMinSpeed, normalMaxSpeed int
var normalSpeedStep int
var normalCharset int

var mutantHead, mutantNeck, mutantTail Color
var mutantMinLen, mutantMaxLen int
var mutantMinSpeed, mutantMaxSpeed int
var mutantSpeedStep int
var mutantCharset int
var mutantChance float32

var frameDuration int
var overlap int
var maxDropsPerColumn int
var syncSpeed int

func handleCommand_set(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) < 2 {
		return ""
	}
	switch args[0] {
		case "normalHead":
			color, err := parseColor(args[1])
			if err == nil {
				normalHead = color
			}
		case "normalNeck":
			color, err := parseColor(args[1])
			if err == nil {
				normalNeck = color
			}
		case "normalTail":
			color, err := parseColor(args[1])
			if err == nil {
				normalTail = color
			}
		case "mutantHead":
			color, err := parseColor(args[1])
			if err == nil {
				mutantHead = color
			}
	}
	return ""
}



func loadDefaults() {
	normalHead = Color{255, 153,   0}
	normalNeck = Color{224,  51,   0}
	normalTail = Color{22,    5,   5}
	// normalHead = Color{136, 204,   0}
	// normalNeck = Color{ 51, 153,  51}
	// normalTail = Color{  0,  25,   0}
	normalMinSpeed = 1
	normalMaxSpeed = 2
	normalSpeedStep = 5
	normalMinLen = 8
	normalMaxLen = 16
	normalCharset = 1

	backCharset = -2

	mutantHead = Color{204, 153, 255}
	mutantNeck = Color{250, 255, 250}
	mutantTail = Color{  0, 204, 122}
	mutantMinSpeed = 1
	mutantMaxSpeed = 1
	mutantSpeedStep = 1
	mutantMinLen = 12
	mutantMaxLen = 24
	mutantCharset = 0
	mutantChance = 0.4

	frameDuration = 20
	overlap = 5
	maxDropsPerColumn = 1
	reservedHeight = 0
	syncSpeed = -1
	rainStatus = true
}

func getConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return dir, err
	}
	dir = filepath.Join(dir, "letterfall")
	return dir, nil
}

func scanFile(f *os.File, hnd *CommandHandler) {
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		hnd.eval(line)
	}
}

func readCommandLine() {
	hnd := NewCommandHandler(HANDLER_TYPE_INIT)

	// files      := pflag.StringSliceP("file"    , "f", []string{}, "[path] file for reading commands")
	// configPath := pflag.StringP     ("config"  , "c", ""        , "[path] configuration file")
	socket     := pflag.StringP     ("socket"  , "s", "0"       , "create a server socket to read commands")
	sockpath   := pflag.Bool("socket-path", false, "print socket path and exit")
	pflag.Parse()

	configDir, err := getConfigDir()

	if err == nil {
		if pflag.NArg() > 0 {
			profile := filepath.Join(configDir, pflag.Args()[0], ".letterfall")
			f, err := os.Open(profile)
			if err == nil {
				scanFile(f, hnd)
			} else {
				fmt.Printf("error loading profile %s: %s", pflag.Args()[0], err.Error())
				os.Exit(1)
			}
		} else {
			profile := filepath.Join(configDir, "config", ".letterfall")
			f, err := os.Open(profile)
			if err == nil {
				scanFile(f, hnd)
			}
		}
	}

	// Start the server socket
	if flag := pflag.Lookup("socket"); flag.Changed {
		socketPath := os.TempDir() + "/letterfall." + *socket + ".sock"
		os.Remove(socketPath)
		if *sockpath {
			fmt.Printf("%s", socketPath)
			os.Exit(0)
		}
		go initServer(socketPath)
	}

	// proccess non-flag arguments as comands
	// for _, cmd := range pflag.Args() {
		// hnd.eval(cmd)
	// }

	// Check if our program was redirected from a pipe
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		// its a pipe, read stdin as a command source
		go func() {
			hnd := NewCommandHandler(HANDLER_TYPE_FILE)
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				line := sc.Text()
				hnd.eval(line)
			}
		}()
	}
}

func initServer(path string) {
	l, err := net.Listen("unix", path)
	if err != nil { panic(err) }

	defer l.Close()
	defer os.Remove(path)
	for {
		conn, err := l.Accept()
		if err != nil { continue }
		go handleConn(conn)
	}	
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewScanner(conn)
	hnd := NewCommandHandler(HANDLER_TYPE_CONN)
	hnd.writer = bufio.NewWriter(conn)

	for reader.Scan() {
		line := reader.Text()
		hnd.eval(line)
	}
}

func readCommandFile(path string) {
	f, err := os.Open(path)
	if err != nil { return }
	defer f.Close()
	hnd := NewCommandHandler(HANDLER_TYPE_FILE)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		hnd.eval(line)
	}
}

