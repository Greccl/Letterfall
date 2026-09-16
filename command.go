package main

import (
	"github.com/google/shlex"
	"github.com/spf13/pflag"
	"sync"
	"bufio"
	"strconv"
	"time"
)



type Command struct {
	flagset *pflag.FlagSet
	doInMain bool
	handle func(*pflag.FlagSet)string
	mu sync.Mutex
}

func NewCommandMap(kind int) map[string]*Command {
	var cmd *Command
	var fset *pflag.FlagSet
	var commands = make(map[string]*Command)

	// SET
	fset = pflag.NewFlagSet("set", pflag.ContinueOnError)
	cmd = new(Command)
	cmd.flagset = fset
	cmd.doInMain = true
	cmd.handle = handleCommand_set
	commands["set"] = cmd

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
	if kind != HANDLER_TYPE_SHELL {
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
	if kind != HANDLER_TYPE_INIT && kind != HANDLER_TYPE_SHELL {
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

const (
	HANDLER_TYPE_MAIN int = iota
	HANDLER_TYPE_INIT
	HANDLER_TYPE_CONN
	HANDLER_TYPE_FILE
	HANDLER_TYPE_SHELL
)

type CommandHandler struct {
	commands map[string]*Command
	writer *bufio.Writer
}

func NewCommandHandler(kind int) *CommandHandler {
	self := new(CommandHandler)
	self.commands = NewCommandMap(kind)
	return self
}

func (self *CommandHandler) eval(line string) {
	args, err := shlex.Split(line)
	if err != nil || len(args) < 1 { return }
	cmd, exists := self.commands[args[0]]
	if !exists { return }
	cmd.mu.Lock()
	if len(args) > 1 {
		cmd.flagset.Parse(args[1:])
	}
	if cmd.doInMain {
		ch_HandleRequests <- HandleRequest{self, cmd}
		return
	}
	self.do(cmd)
}

func (self *CommandHandler) do(cmd *Command) {
	result := cmd.handle(cmd.flagset)
	cmd.flagset.VisitAll(resetFlag)
	cmd.mu.Unlock()
	if self.writer != nil && len(result) > 0 {
		self.writer.WriteString(result)
		self.writer.WriteString("\n")
		self.writer.Flush()
	}
}



//
// Helpers for FlagSet manipulation
//

func resetFlag(f *pflag.Flag) {
	f.Value.Set(f.DefValue)
	f.Changed = false
}

func getBool(fs *pflag.FlagSet, name string) (bool,bool) {
	f := fs.Lookup(name)
	if !f.Changed { return false, false }
	value, _ := fs.GetBool(name)
	return true, value
}

func getInt(fs *pflag.FlagSet, name string) (bool,int) {
	f := fs.Lookup(name)
	if !f.Changed { return false, 0 }
	value, _ := fs.GetInt(name)
	return true, value
}

func getIntSlice(fs *pflag.FlagSet, name string) (bool,[]int) {
	f := fs.Lookup(name)
	if !f.Changed { return false, []int{} }
	value, _ := fs.GetIntSlice(name)
	return true, value
}

func getString(fs *pflag.FlagSet, name string) (bool,string) {
	f := fs.Lookup(name)
	if !f.Changed { return false, "" }
	value, _ := fs.GetString(name)
	return true, value
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
