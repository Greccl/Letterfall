package main

import (
	"github.com/Greccl/tcell/v2"
)



type Column struct {
	drops []Drop
	count int
	x int
	layer int8
}

func (self *Column) resize() {
	for i := range self.drops {
		self.drops[i].runes = SliceResize(self.drops[i].runes, scrh)
	}
	// self.backs.runes = SliceResize(self.backs.runes, scrh)
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
	if self.count == 0 { return }

	var d *Drop

	// Advance
	for i:=self.count-1; i>=0; i-- {
		d = &self.drops[i]

		switch d.class {
			case CLASS_NORMAL:
				g := &normalSyncGroups[d.group]
				if g.advance > 0 {
					d.pos += g.advance
					d.dirty = true
				}
			case CLASS_MUTANT:
				d.count += d.speed * dt
				if d.count < 1.0 { continue }
				d.dirty = true
				adv := int(d.count)
				d.pos += adv
				d.count -= float64(adv)
			case CLASS_BACK:
				d.count += d.speed * dt * float64(BACK_SEGMENTS)
				if d.count < 1.0 { continue }
				d.dirty = true
				if d.group < BACK_SEGMENTS-1 {
					d.group++
				} else {
					d.group = 0
					d.pos++
				}
				d.count = 0.0
		}

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
		screen_releaseCell(self.layer, self.x, scrh-1)
	}

	// Drawing
	force := false
	for i:=0; i<self.count; i++ {
		force = force || self.drops[i].dirty
		if force && d.pos >= 0 {
			if d.class == CLASS_BACK {
				self.drawBackDrop(i)
			} else {
				self.drawNormalDrop(i)
			}
		}
	}
}

func (self *Column) drawNormalDrop(i int) {
	d := &self.drops[i]
	s := tcell.StyleDefault
	var y int

	head := normalHead
	neck := normalNeck
	tail := normalTail
	if d.class == CLASS_MUTANT {
		head = mutantHead
		neck = mutantNeck
		tail = mutantTail		
	}
	
	l := d.pos - d.end + 1
	for p := 0; p <= l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= scrh { continue }
		if p == 0 {
			s.SetForegroundRGB(head.r, head.g, head.b)
		} else if p == l {
			screen_releaseCell(self.layer, self.x, y)
			continue
		} else {
			alfa := int32((d.length - p)*1000/d.length)
			var c Color
			c = blend(neck, tail, 1000-alfa)
			s.SetForegroundRGB(c.r, c.g, c.b)
		}
		screen_drawCell(self.layer, self.x, y, d.runes[y], s)
	}
	
	d.dirty = false
}

func (self *Column) drawBackDrop(i int) {
	d := &self.drops[i]
	s := tcell.StyleDefault
	var y int

	// head := Color{48, 50, 45}
	// head := Color{80, 83, 75}
	head := Color{64, 99, 60}
	tail := Color{32, 33, 30}
	
	l := d.pos - d.end + 1
	for p := 0; p < l; p++ {
		y = d.pos-p
		if y < 0 { return }
		if y >= scrh { continue }
		alpha := backAlphas[p*BACK_SEGMENTS+d.group]
		var c Color
		c = blend(head, tail, 1000-alpha)
		s.SetForegroundRGB(c.r, c.g, c.b)
		screen_drawCell(self.layer, self.x, y, d.runes[y], s)
	}
	screen_releaseCell(self.layer, self.x, y - 1)

	d.dirty = false
}
