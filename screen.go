package main

import (
	"github.com/Greccl/tcell/v2"

)



const (
	LAYER_BASE int8 = iota
	LAYER_BACK
	LAYER_RAIN
	LAYER_TEXT
	LAYER_OVERLAY
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



func screen_resize() {
	scrw, scrh = scr.Size()
	scrh -= reservedHeight
	oddOffset = scrw % 2
	state = SliceResize(state, scrw)
	for i := range state {
		state[i] = SliceResize(state[i], scrh)
	}
	rain_resize()
	for y:=0; y<scrh; y++ {
		for x:=0; x<scrw; x++ {
			owner := state[x][y].owner
			cell := &state[x][y].data[owner]
			screen_printCell(x, y, cell.r, cell.s)
		}
	}
}

func screen_reserveHeight(h int) {
	reservedHeight += h
	if reservedHeight > scrh { reservedHeight = scrh }
	if reservedHeight < 0 { reservedHeight = 0 }
	screen_resize()
}

func screen_damage() {
	select {
		case ch_Draw <- true:
		default:
	}
}

func screen_printText(x, y int, s string) {
	for i, r := range s {
		scr.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
	screen_damage()
}

func screen_printCell(x, y int, r rune, s tcell.Style) {
	scr.SetContent(x, y, r, nil, s)
	screen_damage()
}

func screen_drawCell(layer int8, x, y int, r rune, s tcell.Style) {
	if x < 0 || x >= scrw { return }
	if y < 0 || y >= scrh { return }
	cell := &state[x][y]
	cell.data[layer].r = r
	cell.data[layer].s = s
	if layer > cell.owner { cell.owner = layer }
	if layer == cell.owner {
		scr.SetContent(x, y, r, nil, s)
		screen_damage()
	}
}

func screen_releaseCell(level int8, x, y int) {
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
			screen_damage()
			return
		}
	}
	scr.SetContent(x, y, ' ', nil, tcell.StyleDefault)
	screen_damage()
}


func screen_message() {
	
}
