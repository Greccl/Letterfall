package main

import (
	"math/rand/v2"
	"github.com/Greccl/tcell/v2"
	"github.com/spf13/pflag"
)





var boxes []*Text

func findTextById(id int) *Text {
	for i := range boxes {
		if boxes[i].id == id {
			return boxes[i]
		}
	}
	return nil
}

func findTextByName(name string) *Text {
	for i := range boxes {
		if boxes[i].name == name {
			return boxes[i]
		}
	}
	return nil
}





type TextAnimator interface {
	init(t *Text)
	tick(t *Text)
	draw(t *Text)
}



type Text struct {
	id int
	name string
	runes []rune
	x, y int
	halign, valign bool
	setx, sety int
	fg, bg Color
	anim TextAnimator
}

func NewText() *Text {
	self := new(Text)
	// self.x = -1
	self.y = -1
	self.setAnimation("none")
	self.fg = Color{255, 255, 255}
	self.bg = Color{  0,   0,   0}
	return self
}

func (self *Text) movex() {
	var x int
	if self.halign { x = (scrw-len(self.runes)) / 2 } else
	if self.setx < 0 { x = scrw - len(self.runes) + 1 }
	x += self.setx
	if x != self.x {
		self.releaseAll()
		self.x = x
	}
}

func (self *Text) movey() {
	var y int
	if self.valign { y = scrh / 2 } else
	if self.sety < 0 { y = scrh }
	y += self.sety
	if y != self.y {
		self.releaseAll()
		self.y = y
	}
}

func (t *Text) setText(str string) {
	str = " " + str + " "
	oldLen := len(t.runes) - 1
	t.runes = []rune(str)
	for oldLen >= len(t.runes) {
		releaseCell(LAYER_TEXT, t.x+oldLen, t.y)
		oldLen--
	}
	t.anim.init(t)
}

func (t *Text) setAnimation(name string) {
	switch name {
		case "none":
			t.anim = new(TextAnimator0)
		case "fadein":
			t.anim = new(TextAnimator1)
		case "shine":
			t.anim = new(TextAnimator2)
		case "progresive":
			t.anim = new(TextAnimator3)
		default:
			return
	}
	t.anim.init(t)
}

func (t *Text) releaseAll() {
	for x:=0; x<len(t.runes); x++ {
		releaseCell(LAYER_TEXT, t.x+x, t.y)
	}
}

func (t *Text) autoremove() {
	for i := range boxes {
		if boxes[i] == t {
			boxes = SliceRemove(boxes, i)
			break
		}
	}
	t.releaseAll()
}





//
//   Animator 0
//
type TextAnimator0 struct {
	drawn bool
}

func (self *TextAnimator0) init(t *Text) {

}

func (self *TextAnimator0) tick(t *Text) {
	if self.drawn { return }
	self.draw(t)
	self.drawn = true
}

func (self *TextAnimator0) draw(t *Text) {
	s := tcell.StyleDefault
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	s.SetForegroundRGB(t.fg.r, t.fg.g, t.fg.b)
	for i, r := range t.runes {
		drawCell(LAYER_TEXT, t.x+i, t.y, r, s)
	}
}





//
//   Animator 1
//
type TextAnimator1 struct {
	blends []int32
	fg Color
	i int
}

func (self *TextAnimator1) init(t *Text) {
	self.blends = SliceResize(self.blends, len(t.runes))
	for i := range self.blends {
		self.blends[i] = (rand.Int32N(5)*60) - 200
	}
}

func (self *TextAnimator1) tick(t *Text) {
	self.tick2(t, true)
}

func (self *TextAnimator1) tick2(t *Text, tick bool) {
	for i := range self.blends {
		blender := t.fg
		alfa := int32(0)
		if self.blends[i] < 1000 {
			if tick { self.blends[i] += 20 }
			alfa = self.blends[i]
			if alfa < 0 { alfa = 0 }
			blender = Color{}
		} else if tick {
			continue
		}
		self.fg = blend(t.fg, blender, 1000-alfa)
		self.i = i
		self.draw(t)
	}
}

func (self *TextAnimator1) draw(t *Text) {
	if self.i < 0 {
		self.tick2(t, false)
		self.i = -1
		return
	}
	s := tcell.StyleDefault
	s.SetForegroundRGB(self.fg.r, self.fg.g, self.fg.b)
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	drawCell(LAYER_TEXT, t.x+self.i, t.y, t.runes[self.i], s)
	self.i = -1
}





//
//   Animator 2
//
type TextAnimator2 struct {
	blends []int32
	fg Color
	i int
}

func (self *TextAnimator2) init(t *Text) {
	self.blends = SliceResize(self.blends, len(t.runes))
	for i := range self.blends {
		self.blends[i] = 1
	}
}

func (self *TextAnimator2) tick(t *Text) {
	self.tick2(t, true)
}

func (self *TextAnimator2) tick2(t *Text, tick bool) {
	for i := range self.blends {
		blender := t.fg
		if tick && rand.Float32() < 0.0100 {
			self.blends[i] = 1000
		}
		if self.blends[i] > 100 {
			if tick { self.blends[i] -= 50 }
			blender = Color{255, 255, 255}
		} else if self.blends[i] > 0 {
			self.blends[i] = 0
		} else if tick {
			continue
		}
		self.fg = blend(t.fg, blender, self.blends[i])
		self.i = i
		self.draw(t)
	}
}

func (self *TextAnimator2) draw(t *Text) {
	if self.i < 0 {
		self.tick2(t, false)
		self.i = -1
		return
	}
	s := tcell.StyleDefault
	s.SetForegroundRGB(self.fg.r, self.fg.g, self.fg.b)
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	drawCell(LAYER_TEXT, t.x+self.i, t.y, t.runes[self.i], s)
	self.i = -1
}





//
//   Animator 3
//
type TextAnimator3 struct {
	last int
}

func (self *TextAnimator3) init(t *Text) {
	self.last = -1
}

func (self *TextAnimator3) tick(t *Text) {
	if self.last >= len(t.runes) { return }
	self.last++
	if self.last < len(t.runes) {
		s := tcell.StyleDefault
		s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
		s.SetForegroundRGB(t.fg.r, t.fg.g, t.fg.b)
		drawCell(LAYER_TEXT, t.x+self.last, t.y, t.runes[self.last], s)		
	}
}

func (self *TextAnimator3) draw(t *Text) {
	s := tcell.StyleDefault
	s.SetBackgroundRGB(t.bg.r, t.bg.g, t.bg.b)
	s.SetForegroundRGB(t.fg.r, t.fg.g, t.fg.b)
	for i:=0; i<self.last; i++ {
		drawCell(LAYER_TEXT, t.x+i, t.y, t.runes[i], s)
	}
}






//
// ***text*** command handler
//
func handleCommand_text(fs *pflag.FlagSet) string {
	var t *Text
	var draw, movex, movey bool

	// Find target text
	id, _ := fs.GetInt("id")
	if id > 0 {
		t = findTextById(id)
	}
	name, _ := fs.GetString("name")
	if t == nil && name != "" {
		t = findTextByName(name)
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
		t = NewText()
		t.id = id
		t.name = name
		boxes = append(boxes, t)
		movey = true
	}

	// Text
	s := fs.Arg(0)
	if len(s) > 0 {
		t.setText(s)
		movex = true
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
		movex = true
	}
	if changed, value := getBool(fs, "valign"); changed {
		t.valign = value
		movey = true
	}

	// Position
	if changed, value := getIntSlice(fs, "position"); changed {
		if len(value) == 2 {
			t.setx = value[0]
			t.sety = value[1]
			movex = true
			movey = true
		}
	} else {
		if changed, value := getInt(fs, "x"); changed {
			t.setx = value
			movex = true
		}
		if changed, value := getInt(fs, "y"); changed {
			t.sety = value
			movey = true
		}
	}

	// Animation type
	if changed, value := getString(fs, "animation"); changed {
		t.setAnimation(value)
		draw = true
	}

	// recalculate position
	if movex {
		t.movex()
		draw = true
	}

	if movey {
		t.movey()
		draw = true
	}

	// redraw needed
	if draw {
		t.anim.draw(t)
	}

	return ""
}





//
// Text animation loop
//

func tick_text() {
	for i := range boxes {
		t := boxes[i]
		t.anim.tick(t)
	}
}
