package main

import (
	"github.com/Greccl/tcell/v2"
    "fmt"
    // "image/color"
    "strconv"
    "strings"
)





type Color struct {
	r, g, b int32
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

	// #RRGGBB
	if strings.HasPrefix(s, "#") {
		if len(s) != 7 { return }
		var i64 int64
		i64, err = strconv.ParseInt(s[1:3], 16, 0)
		if err != nil { return }
		c.r = int32(i64)
		i64, err = strconv.ParseInt(s[3:5], 16, 0)
		if err != nil { return }
		c.g = int32(i64)
		i64, err = strconv.ParseInt(s[5:7], 16, 0)
		if err != nil { return }
		c.b = int32(i64)
		// if true {panic(err)}
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





const (
	LAYER_BASE int8 = iota
	LAYER_BACK
	LAYER_RAIN
	LAYER_TEXT
	LAYER_COUNT
)

type Cell struct {
	r rune
	s tcell.Style
}

type CellState struct {
	data [LAYER_COUNT]Cell
	owner int8
}





var scr tcell.Screen
var scrw, scrh int
var ch_Draw = make(chan bool, 1)
var state [][]CellState
var oddOffset int
var reservedHeight int





func resize() {
	scrw, scrh = scr.Size()
	scrh -= reservedHeight
	oddOffset = scrw % 2
	state = SliceResize(state, scrw)
	for i := range state {
		state[i] = SliceResize(state[i], scrh)
	}
	rain_resize()
	scr.Clear()
	// resizeEditor()
}

func damage() {
	select {
		case ch_Draw <- true:
		default:
	}
}

func printText(x, y int, s string) {
	for i, r := range s {
		scr.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
	damage()
}

func printCell(x, y int, r rune, s tcell.Style) {
	scr.SetContent(x, y, r, nil, s)
	damage()
}

func drawCell(layer int8, x, y int, r rune, s tcell.Style) {
	if x < 0 || x >= scrw { return }
	if y < 0 || y >= scrh { return }
	cell := &state[x][y]
	cell.data[layer].r = r
	cell.data[layer].s = s
	if layer > cell.owner { cell.owner = layer }
	if layer == cell.owner {
		scr.SetContent(x, y, r, nil, s)
		damage()
	}
// drawOwner(x, y)
}

func drawOwner(x, y int) {
	cell := &state[x][y]
	switch cell.owner {
		case 0:
			scr.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		case 1:
			scr.SetContent(x, y, '1', nil, tcell.StyleDefault)
		case 2:
			scr.SetContent(x, y, '2', nil, tcell.StyleDefault)
	}
	damage()	
}

func releaseCell(level int8, x, y int) {
	if x < 0 || x >= scrw { return }
	if y < 0 || y >= scrh { return }
	cell := &state[x][y]
	cell.data[level].r = 0
	if level < cell.owner { return }
	cell.owner = 0
	for l:=int8(2); l>=0; l-- {
		if cell.data[l].r != 0 {
			cell.owner = l
			scr.SetContent(x, y, cell.data[l].r, nil, cell.data[l].s)
			damage()
// drawOwner(x, y)
			return
		}
	}
	scr.SetContent(x, y, ' ', nil, tcell.StyleDefault)
	damage()
// drawOwner(x, y)
}

