package main

import (
	"github.com/spf13/pflag"
)



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
