package main

import (
	"os"
	"bufio"
	"net"
	"fmt"
	"strings"
	"strconv"
	"path/filepath"
	"github.com/spf13/pflag"
)



var profile string = "default"
// var frameDuration int
// var overlap int

func handleCommand_get(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) < 1 { return "" }
	switch args[0] {
		case "normalColors":
			return fmt.Sprintf("head:%s / neck:%s / tail:%s", normalHead.toString(), normalNeck.toString(), normalTail.toString())
		case "normalSpeed":
			return fmt.Sprintf("min:%f / max:%f / step:%f", normalMinSpeed, normalMaxSpeed, normalSpeedStep)
		case "normalSpeedDebug":
			var s string
			for i := range normalSyncGroups {
				s += fmt.Sprintf(" %.3f", normalSyncGroups[i].speed)
			}
			return fmt.Sprintf("groupCount:%d / groups(%d):%s", normalGroupCount, len(normalSyncGroups), s)
	}
	return ""
}

func handleCommand_set(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) < 2 {
		return "few arguments for set command"
	}
	switch args[0] {
		case "normalHead":
			color, err := parseColor(args[1])
			if err == nil { normalHead = color }
		case "normalNeck":
			color, err := parseColor(args[1])
			if err == nil { normalNeck = color }
		case "normalTail":
			color, err := parseColor(args[1])
			if err == nil { normalTail = color }
		case "normalMinLen":
			value, err := strconv.Atoi(args[1])
			if err == nil { normalMinLen = value }
		case "normalMaxLen":
			value, err := strconv.Atoi(args[1])
			if err == nil { normalMaxLen = value }
		case "normalMinSpeed":
			value, err := strconv.ParseFloat(args[1], 64)
			if err != nil { break }
			normalMinSpeed = value
			normalizeNormalSpeed()
		case "normalMaxSpeed":
			value, err := strconv.ParseFloat(args[1], 64)
			if err != nil { break }
			normalMaxSpeed = value
			normalizeNormalSpeed()
		case "normalSpeedStep":
			value, err := strconv.ParseFloat(args[1], 64)
			if err != nil { break }
			normalSpeedStep = value
			normalizeNormalSpeed()
		case "normalSpeed":
			if len(args) < 4 { break }
			v1, err := strconv.ParseFloat(args[1], 64)
			if err != nil { return err.Error() }
			v2, err := strconv.ParseFloat(args[2], 64)
			if err != nil { return err.Error() }
			v3, err := strconv.ParseFloat(args[3], 64)
			if err != nil { return err.Error() }
			normalMinSpeed = v1
			normalMaxSpeed = v2
			normalSpeedStep = v3
			normalizeNormalSpeed()
		case "normalCharset":
			switch args[1] {
				case "next":
					normalCharset++
					if normalCharset >= CHARSET_COUNT {
						if customCharsetA != nil {
							normalCharset = CHARSET_CUSTOM_A
						} else if customCharsetB != nil {
							normalCharset = CHARSET_CUSTOM_B
						} else {
							normalCharset = CHARSET_DEFAULT
						}
					}
					return "normal chatset -> " + charsetNames[normalCharset]
				case "customA":
					if len(customCharsetA) > 0 {
						normalCharset = CHARSET_CUSTOM_A
						return ""
					}
					return "no custom charset defined"
				case "customB":
					if len(customCharsetB) > 0 {
						normalCharset = CHARSET_CUSTOM_B
						return ""
					}
					return "no custom charset defined"
				default:
					i := CHARSET_DEFAULT
					for ; i<CHARSET_COUNT; i++ {
						if args[1] == charsetNames[i] {
							normalCharset = i
							break
						}
					}
					if i != normalCharset {
						return "invalid charset name: " + args[1]
					}
			}
		case "customCharsetA":
			if len(args[1]) < 1 {
				return "provided charset is empty"
			}
			customCharsetA = []rune(args[1])
			if len(args) >= 3 && args[2] == "now" {
				normalCharset = CHARSET_CUSTOM_A
			}
		case "customCharsetB":
			if len(args[1]) < 1 {
				return "provided charset is empty"
			}
			customCharsetB = []rune(args[1])
			if len(args) >= 3 && args[2] == "now" {
				normalCharset = CHARSET_CUSTOM_B
			}

		case "mutantHead":
			color, err := parseColor(args[1])
			if err == nil { mutantHead = color }
		case "mutantNeck":
			color, err := parseColor(args[1])
			if err == nil { mutantNeck = color }
		case "mutantTail":
			color, err := parseColor(args[1])
			if err == nil { mutantTail = color }
		default:
			return "invalid property: " + args[0]
	}
	return ""
}

func handleCommand_save(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) > 1 {
		return ""
	}
	if len(args) == 0 {
		saveProfile(profile)
	} else {
		saveProfile(args[0])
	}
	return ""
}

func saveProfile(name string) {
	configDir, err := getConfigDir()
	if err != nil { return }

	tempPath := filepath.Join(configDir, name + ".temp.letterfall")
	f, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil { return }
	defer f.Close()

	fmt.Fprintln(f, "set normalHead", normalHead.toString())
	fmt.Fprintln(f, "set normalNeck", normalNeck.toString())
	fmt.Fprintln(f, "set normalTail", normalTail.toString())

	basePath := filepath.Join(configDir, name + ".letterfall")
	base, err := os.Open(basePath)
	if err == nil {
		defer base.Close()
		scanner := bufio.NewScanner(base)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "set ") {
				fmt.Fprintln(f, line)
			}
		}
	}

	f.Close()
	base.Close()
	os.Rename(tempPath, basePath)
}

func loadDefaults() {
	normalHead = Color{255, 153,   0}
	normalNeck = Color{224,  51,   0}
	normalTail = Color{22,    5,   5}
	// normalHead = Color{136, 204,   0}
	// normalNeck = Color{ 51, 153,  51}
	// normalTail = Color{  0,  25,   0}
	normalMinSpeed = 5
	normalMaxSpeed = 5
	normalSpeedStep = 1
	normalMinLen = 8
	normalMaxLen = 16
	normalCharset = CHARSET_DEFAULT

	backCharset = CHARSET_BACKDROP

	mutantHead = Color{204, 153, 255}
	mutantNeck = Color{250, 255, 250}
	mutantTail = Color{  0, 204, 122}
	mutantMinSpeed = 18.0
	mutantMaxSpeed = 24.0
	// mutantSpeedStep = 1
	mutantMinLen = 12
	mutantMaxLen = 24
	mutantCharset = CHARSET_UPPERCASE
	mutantChance = 0.1

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
	hnd := NewCommandHandler(SCOPE_INIT)

	socket := pflag.StringP("socket", "s", "0", "create a server socket to read commands")
	sockpath := pflag.Bool("socket-path", false, "print socket path and exit")
	pflag.Parse()

	configDir, err := getConfigDir()

	if err == nil {
		if pflag.NArg() > 0 {
			path := filepath.Join(configDir, pflag.Args()[0] + ".letterfall")
			f, err := os.Open(path)
			if err == nil {
				profile = pflag.Args()[0]
				scanFile(f, hnd)
			} else {
				fmt.Printf("error loading profile %s: %s", pflag.Args()[0], err.Error())
				os.Exit(1)
			}
		} else {
			path := filepath.Join(configDir, "default.letterfall")
			profile = "default"
			f, err := os.Open(path)
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

	// Check if our program was redirected from a pipe
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		// its a pipe, read stdin as a command source
		go func() {
			hnd := NewCommandHandler(SCOPE_FILE)
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
	hnd := NewCommandHandler(SCOPE_CONN)
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
	hnd := NewCommandHandler(SCOPE_FILE)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		hnd.eval(line)
	}
}

