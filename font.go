package main

import (
	"fmt"
	// "os"
	// "time"
	"github.com/Greccl/tcell/v2"
)






type Grid struct {
	pix [4][4]rune
	posx, posy int
}

var output, editor Grid
var selx, sely int

func onKey_editor(key tcell.Key) {
	switch key {
		case tcell.KeyUp:
			sely--
			if sely < 0 { sely = 0 }
		case tcell.KeyDown:
			sely++
			if sely > 3 { sely = 3 }
		case tcell.KeyLeft:
			selx--
			if selx < 0 { selx = 0 }
		case tcell.KeyRight:
			selx++
			if selx > 3 { selx = 3 }
	}
}

func onRune_editor(r rune) {
	switch r {
		case 'z':
			setHandler(CHARS)
	}	
}





var chars = []rune("#+(!@:)")

var charPos int
var charField string

func onKey_chars(key tcell.Key) {
	switch key {
		case tcell.KeyLeft:
			charPos--
			if charPos < 0 { charPos = 0 }
			charField = fmt.Sprintf("%d", chars[charPos])
		case tcell.KeyRight:
			charPos++
			if charPos >= len(chars) { charPos = len(chars)-1 }
			charField = fmt.Sprintf("%d", chars[charPos])
	}
}

func onRune_chars(r rune) {
	switch r {
		case 'z':
			setHandler(EDITOR)
	}
}















const (
	CHARS int = iota
	EDITOR
	TARGET
)

func setHandler(h int) {
	switch h {
		case CHARS:
			onKey = onKey_chars
			onRune = onRune_chars
		case EDITOR:
			onKey = onKey_editor
			onRune = onRune_editor
		case TARGET:
	}
}

func resizeEditor() {
	if onKey == nil { onKey = onKey_editor }
	if onRune == nil { onRune = onRune_editor }
	output.posx = 10
	editor.posx = 0
	drawEditor()
}

func drawEditor() {
	s := tcell.StyleDefault
	s.SetBackgroundRGB(50, 50, 50)

	basex, basey := output.posx, output.posy
	printText(basex, basey, "output")
	basey++
	for y := range output.pix {
		for x := range output.pix[y] {
			printCell(basex+x, basey+y, output.pix[y][x], s)
		}
	}

	basex, basey = editor.posx, editor.posy
	printText(basex, basey, "editor")
	basey++
	for y := range editor.pix {
		for x := range editor.pix[y] {
			printCell(basex+x, basey+y, editor.pix[y][x], s)
		}
	}
	s.SetBackgroundRGB(0, 100, 0)
	printCell(basex+selx, basey+sely, editor.pix[sely][selx], s)
	
	basey = 7
	s.SetBackgroundRGB(0,0,0)
	for x := range chars {
		printCell(x, basey, chars[x], s)
		printCell(x, basey+1, ' ', s)
	}
	basey++
	printCell(charPos, basey, '^', s)
	basey++
	printText(0, basey, charField)
	
}
