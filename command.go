package main

import (
	"github.com/google/shlex"
	"github.com/spf13/pflag"
	"sync"
	"bufio"
	"strconv"
	"time"
)



const (
	SCOPE_MAIN int = 1 << iota
	SCOPE_INIT
	SCOPE_CONN
	SCOPE_FILE
	SCOPE_SHELL
)

type Command struct {
	flagset *pflag.FlagSet
	doInMain bool
	handle func(*pflag.FlagSet)string
	mu sync.Mutex
}

func buildCommandMap(scope int) map[string]*Command {
	var cmd *Command
	var fset *pflag.FlagSet
	var commands = make(map[string]*Command)

	// SET
	fset = pflag.NewFlagSet("set", pflag.ContinueOnError)
	cmd = new(Command)
	cmd.flagset = fset
	cmd.doInMain = false
	cmd.handle = handleCommand_set
	commands["set"] = cmd

	// Limit of init command list
	if scope == SCOPE_INIT {
		return commands
	}

	// GET
	fset = pflag.NewFlagSet("get", pflag.ContinueOnError)
	cmd = new(Command)
	cmd.flagset = fset
	cmd.doInMain = false
	cmd.handle = handleCommand_get
	commands["get"] = cmd

	// PING
	cmd = new(Command)
	cmd.doInMain = false
	cmd.handle = handleCommand_ping
	commands["ping"] = cmd

	// SAVE
	if scope & SCOPE_SHELL != 0 {
		fset = pflag.NewFlagSet("save", pflag.ContinueOnError)
		cmd = new(Command)
		cmd.flagset = fset
		cmd.doInMain = true
		cmd.handle = handleCommand_save
		commands["save"] = cmd
	}

	// TEXT
	fset = pflag.NewFlagSet("text", pflag.ContinueOnError)
	fset.IntSliceP("position"  , "p", []int{}, "a compact way to set x and y position")
	fset.IntP     ("x"         , "x", 0      , "x position of text box")
	fset.IntP     ("y"         , "y", 0      , "y position of text box")
	fset.BoolP    ("halign"    , "h", false  , "evaluate horizontal position from center of screen")
	fset.BoolP    ("valign"    , "v", false  , "evaluate vertical position from center of screen")
	fset.IntP     ("id"        , "i", 0      , "identifier (integer value) for the text box")
	fset.StringP  ("name"      , "n", ""     , "identifier (string) for the text box")
	fset.BoolP    ("kill"      , "k", false  , "try remove given box")
	fset.StringP  ("animation" , "a", ""     , "name of animation")
	fset.StringP  ("foreground", "f", ""     , "set foreground colour")
	fset.StringP  ("background", "b", ""     , "set background colour")
	cmd = new(Command)
	cmd.flagset = fset
	cmd.doInMain = true
	cmd.handle = handleCommand_text
	commands["text"] = cmd

	// BANNER
	fset = pflag.NewFlagSet("banner", pflag.ContinueOnError)
	cmd = new(Command)
	cmd.flagset = fset
	cmd.doInMain = true
	cmd.handle = handleCommand_banner
	commands["banner"] = cmd

	// STREAM
	if scope & SCOPE_SHELL == 0 {
		fset = pflag.NewFlagSet("stream", pflag.ContinueOnError)
		fset.IntSliceP("position"  , "p", []int{}, "a compact way to set x and y position")
		fset.IntP     ("x"         , "x", 0      , "x position of text box")
		fset.IntP     ("y"         , "y", 0      , "y position of text box")
		fset.IntP     ("width"     , "w", 10     , "box's width")
		fset.BoolP    ("halign"    , "h", false  , "evaluate horizontal position from center of screen")
		fset.BoolP    ("valign"    , "v", false  , "evaluate vertical position from center of screen")
		fset.IntP     ("id"        , "i", 0      , "identifier (integer value) for the text box")
		fset.StringP  ("name"      , "n", ""     , "identifier (string) for the text box")
		fset.BoolP    ("kill"      , "k", false  , "try remove given box")
		fset.StringP  ("animation" , "a", ""     , "name of animation")
		fset.StringP  ("foreground", "f", ""     , "set foreground colour")
		fset.StringP  ("background", "b", ""     , "set background colour")
		cmd = new(Command)
		cmd.flagset = fset
		cmd.doInMain = true
		cmd.handle = handleCommand_stream
		commands["stream"] = cmd
	}

	// SLEEP
	if scope & SCOPE_SHELL == 0 {
		fset = pflag.NewFlagSet("sleep", pflag.ContinueOnError)
		cmd = new(Command)
		cmd.flagset = fset
		cmd.doInMain = false
		cmd.handle = handleCommand_sleep
		commands["sleep"] = cmd
	}

	return commands
}





// CommandHandler allows reading line of commands
// from many diferent sources. once a line is readed,
// eval() transform a line in a request. if command
// is suposed to be ran in main loop it schedules a
// request by the request channel polled in main().
// if not, it is executed immediatelly.

type CommandHandler struct {
	commands map[string]*Command
	writer *bufio.Writer
}

func NewCommandHandler(scope int) *CommandHandler {
	self := &CommandHandler{}
	self.commands = buildCommandMap(scope)
	return self
}

func (self *CommandHandler) eval(line string) {
	args, err := shlex.Split(line)
	if err != nil || len(args) < 1 { return }
	cmd, exists := self.commands[args[0]]
	if !exists { return }
	cmd.mu.Lock()
	if len(args) > 1 && cmd.flagset != nil {
		cmd.flagset.Parse(args[1:])
	}
	if cmd.doInMain {
		ch_HandleRequests <- HandleRequest{self, cmd}
	} else {
		self.do(cmd)
	}
}

func (self *CommandHandler) do(cmd *Command) {
	result := cmd.handle(cmd.flagset)
	if cmd.flagset != nil {
		cmd.flagset.VisitAll(resetFlag)
	}
	cmd.mu.Unlock()
	if self.writer != nil && len(result) > 0 {
		self.writer.WriteString(result)
		self.writer.WriteString("\n")
		self.writer.Flush()
	}
}








func handleCommand_ping(fs *pflag.FlagSet) string {
	return "pong"
}


func handleCommand_state(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) == 0 {
		// maybe printout state info?
		return ""
	}
	switch args[0] {
		case "pause":
			rainStatus = !rainStatus
	}
	return ""
}

func handleCommand_sleep(fs *pflag.FlagSet) string {
	args := fs.Args()
	if len(args) == 0 {
		return ""
	}
	secs, err := strconv.ParseFloat(fs.Arg(0), 64)
	if err != nil {
		return err.Error()
	}
	d := time.Duration(secs * float64(time.Second))
	time.Sleep(d)
	return ""
}
