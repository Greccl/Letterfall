package main

import (
	"github.com/google/shlex"
	"github.com/spf13/pflag"
	"sync"
	"bufio"
)








type Command struct {
	flagset *pflag.FlagSet
	doInMain bool
	handle func(*pflag.FlagSet)string
	mu sync.Mutex
}

type HandleRequest struct {
	hnd *CommandHandler
	cmd *Command
}





const (
	HANDLER_TYPE_MAIN int = iota
	HANDLER_TYPE_INIT
	HANDLER_TYPE_CONN
	HANDLER_TYPE_FILE
)

func NewCommandMap(kind int) map[string]*Command {
	commands := make(map[string]*Command)

	var cmd *Command
	var fset *pflag.FlagSet

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
	// cmd := new(Command)
	// commands["text"] = cmd
	// cmd.fs = fset
	// cmd.fn = handleCommand_text

	// fset = pflag.NewFlagSet("state", pflag.ContinueOnError)
	// commands["state"] = fset
	// handlers["state"] = handleCommand_state

/*
	fset = pflag.NewFlagSet("delay", pflag.ContinueOnError)
	commands["automate"] = fset
	handlers["state"] = handleCommand_state
*/
	return commands
}








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
	if self.writer != nil {
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
