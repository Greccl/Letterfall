package main

import (
	"github.com/spf13/pflag"
	"strings"
	"fmt"
	"strconv"
)

//
// INternal color representation
//

type Color struct {
	r, g, b int32
}

func (c *Color) toString() string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.r, c.g, c.b)
}

func blend(a, b Color, alfa int32) Color {
	var c Color
	beta := 1000 - alfa
	c.r = ((b.r * alfa) + (a.r * beta)) / 1000
	c.g = ((b.g * alfa) + (a.g * beta)) / 1000
	c.b = ((b.b * alfa) + (a.b * beta)) / 1000
	return c
}

func parseColor(s string) (c Color, err error) {
	s = strings.TrimSpace(s)

	// 0xRRGGBB
	if strings.HasPrefix(s, "0x") {
		if len(s) != 8 { return }
		var i64 int64
		i64, err = strconv.ParseInt(s[2:4], 16, 0)
		if err != nil { return }
		c.r = int32(i64)
		i64, err = strconv.ParseInt(s[4:6], 16, 0)
		if err != nil { return }
		c.g = int32(i64)
		i64, err = strconv.ParseInt(s[6:8], 16, 0)
		if err != nil { return }
		c.b = int32(i64)
		return c, nil
	}

	// Formato: rgb(R,G,B)
	if strings.HasPrefix(strings.ToLower(s), "rgb(") && strings.HasSuffix(s, ")") {
		inner := s[4:len(s)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 3 { return c, fmt.Errorf("not enough values") }
		var val int
		val, err = strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil { return }
		if val < 0 || val > 255 { return c, fmt.Errorf("red value out of range") }
		c.r = int32(val)
		val, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil { return }
		if val < 0 || val > 255 { return c, fmt.Errorf("green value out of range") }
		c.g = int32(val)
		val, err = strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil { return }
		if val < 0 || val > 255 { return c, fmt.Errorf("blue value out of range") }
		c.b = int32(val)
		return c, nil
	}

	return c, fmt.Errorf("unknown color format")
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
