package main

import (
	// "math/rand/v2"
	"time"
	"github.com/Greccl/tcell/v2"
	"github.com/spf13/pflag"
)





var streams []*Stream

func findStreamById(id int) *Stream {
	for i := range streams {
		if streams[i].id == id {
			return streams[i]
		}
	}
	return nil
}

func findStreamByName(name string) *Stream {
	for i := range streams {
		if streams[i].name == name {
			return streams[i]
		}
	}
	return nil
}





type StreamAnimator interface {
	init(*Stream)
	tick(*Stream,time.Duration)
	draw(*Stream)
	next(*Stream)
}





type Stream struct {
	Widget
	id int
	name string
	items []StreamItem
	fg, bg Color
	anim StreamAnimator
}

func NewStream() *Stream {
	self := new(Stream)
	self.layer = LAYER_TEXT
	self.setw = 1
	self.seth = 1
	return self
}

func (self *Stream) appendString(s string) {
	// self.queue = append(self.queue, s)
	// self.anim.next(self)
}

func (self *Stream) addItem(item StreamItem) {
	self.items = append(self.items, item)
}

func (self *Stream) autoremove() {
	for i := range streams {
		if streams[i] == self {
			streams = SliceRemove(streams, i)
			break
		}
	}
	self.releaseAll()
}

func (self *Stream) setAnimation(int) {

}





const (
	STREAM_IMMEDIATE int8 = iota
	STREAM_VERTICAL
	STREAM_SLIDE
)

const (
	STREAM_ONCE int8 = iota
	STREAM_FOREVER
	SYREAM_FILL
)



type StreamItem struct {
	runes []rune
	mode int8
	loop int8
	adj bool
}

type StreamAnimator0 struct {
	item *StreamItem
	index int
	// elapsed int64

	i int
	delay int
	l, r int
}

func (self *StreamAnimator0) init(t *Stream) {
	self.index = -1
	self.delay = 0
}

func (self *StreamAnimator0) tick(t *Stream, d time.Duration) {
	if self.item == nil { return }
	self.delay++
	if self.delay < 17 { return }
	self.delay = 0

	switch self.item.mode {
		case STREAM_SLIDE:
			x := 0
			y := 0
			for ; x < self.l && x < t.w; x++ {
				t.drawCell(x, y, ' ')
			}
			if self.i < len(self.item.runes) {
				i := self.i
				for x < t.w && i < len(self.item.runes) {
					r := self.item.runes[i]
					t.drawCell(x, y, r)
					x++
					i++
				}
				self.i++
			} else {
				self.next(t)
			}
			for ; x < t.w; x++ {
				t.drawCell(x, y, ' ')
			}
	}
	self.draw(t)
}

func (self *StreamAnimator0) draw(t *Stream) {
	
}

func (self *StreamAnimator0) next(t *Stream) {
	if len(t.items) == 0 { return }
	self.index++
	if self.index >= len(t.items) {
		self.index = 0
		self.item = nil
		return
	}
	self.item = &t.items[self.index]
	self.l = t.w
	self.r = 0
	self.i = 0
}





//
// ***text*** command handler
//
func handleCommand_stream(fs *pflag.FlagSet) string {
	var t *Stream
	var draw, layout bool

	// Find target text
	id, _ := fs.GetInt("id")
	if id > 0 {
		t = findStreamById(id)
	}
	name, _ := fs.GetString("name")
	if t == nil && name != "" {
		t = findStreamByName(name)
	}

	// Issue a kill command
	if b, _ := fs.GetBool("kill"); b {
		if t != nil {
			t.autoremove()
		}
		return ""
	}

	// Create if doesnt exists
	if t == nil {
		t = NewStream()
		t.id = id
		t.name = name
		streams = append(streams, t)
		t.anim = new(StreamAnimator0)
		t.anim.init(t)
		layout = true
	}

	// Text
	s := fs.Arg(0)
	// var msg bool
	if len(s) > 0 {
		var item StreamItem
		item.runes = []rune(s)
		item.adj, _ = fs.GetBool("adjust")
		t.addItem(item)
	}

	// Colours
	s, _ = fs.GetString("foreground")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil {
			t.fg = c
			draw = true
		}
	}

	s, _ = fs.GetString("background")
	if len(s) > 0 {
		c, err := parseColor(s)
		if err == nil {
			t.bg = c
			draw = true
		}
	}

	// Align
	if changed, value := getBool(fs, "halign"); changed {
		t.halign = value
		layout = true
	}
	if changed, value := getBool(fs, "valign"); changed {
		t.valign = value
		layout = true
	}

	// Position
	if changed, value := getIntSlice(fs, "position"); changed {
		if len(value) == 2 {
			t.setx = value[0]
			t.sety = value[1]
			layout = true
		}
	} else {
		if changed, value := getInt(fs, "x"); changed {
			t.setx = value
			layout = true
		}
		if changed, value := getInt(fs, "y"); changed {
			t.sety = value
			layout = true
		}
	}

	// Size
	if changed, value := getInt(fs, "width"); changed {
		if value == -1 {
			t.setw = 10
		} else {
			t.setw = value
		}
		layout = true
	}

	// Animation type
	/*
	if changed, value := getString(fs, "animation"); changed {
		t.setAnimation(value)
		draw = true
	}
	*/

	// recalculate position
	if layout {
		t.layout()
		draw = true
	}

	// redraw needed
	if draw {
		t.anim.draw(t)
	}

	return ""
}




type Widget struct {
	x, y int // Position
	w, h int // Size
	halign, valign bool // Align
	setx, sety int // <- this are configured values
	setw, seth int // <- needed gor recalculations
	layer int8
	style tcell.Style
	style2 tcell.Style
	cx, cy int // cursor position
}

func (self *Widget) setPosition(x, y int) {
	self.setx = x
	self.sety = y
	self.layout()
}

func (self *Widget) setSize(w, h int) {
	self.setw = w
	self.seth = h
	self.layout()
}

func (self *Widget) setFg(r, g, b int32) {
	self.style.SetForegroundRGB(r, g, b)
}

func (self *Widget) setBg(r, g, b int32) {
	self.style.SetBackgroundRGB(r, g, b)
}

func (self *Widget) save(b bool) {
	if b {
		self.style2 = self.style
	} else {
		self.style = self.style2
	}
}

func (self *Widget) drawCell(x, y int, r rune) {
	if x < 0 || x >= self.w { return }
	if y < 0 || y >= self.h { return }
	drawCell(self.layer, self.x+x, self.y+y, r, self.style)
}

func (self *Widget) putc(r rune) {
	drawCell(self.layer, self.x+self.cx, self.y+self.cy, r, self.style)
	self.cursorInc()
}

func (self *Widget) puts(s string) {
	for _, r := range s {
		self.putc(r)
	}
}

func (self *Widget) cursorInc() {
	self.cx++
	if self.cx >= self.w {
		self.cy++
		if self.cy >= self.h {
			self.cy--
		} else {
			self.cx = 0
		}
	}
}

func (self *Widget) cursorMove(x, y int) {
	if x < 0 { x = 0 }
	if x >= self.w { x = self.w - 1}
	if y < 0 { y = 0 }
	if y >= self.h { y = self.h - 1}
	self.cx = x
	self.cy = y
}

func (self *Widget) layout() {
	var x int
	if self.halign { x = (scrw-self.setw) / 2 } else
	if self.setx < 0 { x = scrw - self.setw + 1 }
	x += self.setx
	if x != self.x {
		self.releaseAll()
		self.x = x
	}
	var y int
	if self.valign   { y = (scrh - self.seth) / 2 } else
	if self.sety < 0 { y = scrh - self.seth }
	y += self.sety
	if y != self.y {
		self.releaseAll()
		self.y = y
	}
}

func (self *Widget) releaseAll() {
	for y:=0; y<self.seth; y++ {
		for x:=0; x<self.setw; x++ {
			releaseCell(self.layer, self.x+x, self.y+y)
		}
	}
}





//
// Stream animation loop
//
var streamLastTime time.Time

func tick_stream() {
	now := time.Now()
	d := now.Sub(streamLastTime)
	streamLastTime = now
	for i := range streams {
		t := streams[i]
		t.anim.tick(t, d)
	}
}
