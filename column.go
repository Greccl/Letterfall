package main

import (
	// "fmt"
	// "math/rand/v2"
	"github.com/Greccl/tcell/v2"
)



type Column struct {
	drops []Drop
	count int
	visibleCount int
	x int
	layer int8

	backs Drop
}






func (self *Column) resize() {
	for i := range self.drops {
		self.drops[i].runes = SliceResize(self.drops[i].runes, scrh)
	}
	self.backs.runes = SliceResize(self.backs.runes, scrh)
}

func (self *Column) newDrop() *Drop {
	if self.count == len(self.drops) {
		self.drops = append(self.drops, Drop{})
		self.drops[self.count].runes = make([]rune, scrh)
	}
	self.count++
	return &self.drops[self.count-1]
}

func (self *Column) remove(i int) {
	if i < self.count-1 {
		old := self.drops[i].runes
		copy(self.drops[i:], self.drops[i+1:])
		self.drops[self.count-1].runes = old
	}
	self.count--
}

func (self *Column) tick(dt float64) {
	// self.backTick()
	if self.count == 0 { return }
	self.visibleCount = 0

	var d *Drop

	// Advance
	for i:=self.count-1; i>=0; i-- {
		d = &self.drops[i]
		var g *SyncGroup
		if d.mutant {
			continue
		} else {
			g = &normalSyncGroups[d.group]
		}
		if g.advance > 0 {
			d.pos += g.advance
			d.dirty = true
		} else {
			continue
		}
/*
		d.count += d.speed * dt
		if d.count < 1.0 { continue }
		d.dirty = true
		adv := int(d.count)
		d.pos += adv
		d.count -= float64(adv)
*/

/*		trueAdvance := false
		if d.mutant || syncSpeed <= 0 {
			d.count++
			if d.count >= d.speed { trueAdvance = true }
		} else

		if syncAdvance {
			trueAdvance = true
		}

		if trueAdvance {
			d.pos++
			if d.pos >= 0 { d.dirty = true }
			d.count = 0			
		} else {
			continue
		}
*/		
		// Check overlaps
		if i > 0 {
			if d.pos >= self.drops[i-1].pos {
				self.remove(i-1)
				i--
			}
		}
	}

	// Endings
	d = &self.drops[self.count-1]
	d.end = d.pos - d.length + 1
	if d.end < 0 { d.end = 0 }
	if d.end >= d.pos { d.end = d.pos }
	for i:=self.count-2; i>=0; i-- {
		d = &self.drops[i]
		d.end = d.pos - d.length + 1
		prev := self.drops[i+1].pos + 1
		if d.end < prev { d.end = prev }
		if d.end >= d.pos { d.end = d.pos }
	}

	// Remove completed
	if self.drops[0].end >= scrh {
		self.remove(0)
		releaseCell(self.layer, self.x, scrh-1)
	}

	// Drawing
	force := false
	for i:=0; i<self.count; i++ {
		force = force || self.drops[i].dirty
		if force && d.pos >= 0 {
			self.visibleCount++
			self.draw(i)
		}
	}

	/*
	s := fmt.Sprintf("%2d", self.count)
	for i := range s {
		scr.SetContent(self.x+i, scrh, rune(s[i]), nil, tcell.StyleDefault)
	}	
	*/
}

func (self *Column) draw(i int) {
	d := &self.drops[i]
	s := tcell.StyleDefault
	var y int

	head := normalHead
	neck := normalNeck
	tail := normalTail
	if d.mutant {
		head = mutantHead
		neck = mutantNeck
		tail = mutantTail		
	} else
	if d.back {
		head = Color{90, 90, 100}
		neck = Color{90, 90, 90}
		tail = Color{10, 10, 20}
	} else
	if d.lucent {
		
	}
	
	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= scrh { continue }
		if p == 0 {
			s.SetForegroundRGB(head.r, head.g, head.b)
		} else if p == l {
			releaseCell(self.layer, self.x, y)
			continue
		} else {
			alfa := int32((d.length - p)*1000/d.length)
			var c Color
			c = blend(neck, tail, 1000-alfa)
			s.SetForegroundRGB(c.r, c.g, c.b)
		}
		drawCell(self.layer, self.x, y, d.runes[y], s)
	}
	
	d.dirty = false
}












var backrunes = []rune{9679, 9670, 9643, 9642, 9702, 9711}
