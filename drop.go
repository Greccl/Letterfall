package main

import (
	"math/rand/v2"
)


type Drop struct {
	runes []rune
	mutant bool
	lucent bool
	back bool

	pos int
	length int
	head int
	end int
	dirty bool

	speed int
	count int
}



func (self *Drop) makeNormal() {
	self.mutant = false
	// self.lucent = false
	self.speed = rand.IntN(normalMaxSpeed-normalMinSpeed+1) + normalMinSpeed
	self.speed = self.speed * normalSpeedStep
	self.length = rand.IntN(normalMaxLen-normalMinLen) + normalMinLen
	self.resetRunes(normalCharset)
}

func (self *Drop) makeMutant() {
	self.mutant = true
	// self.lucent = false
	self.speed = rand.IntN(mutantMaxSpeed-mutantMinSpeed+1) + mutantMinSpeed
	self.speed *= mutantSpeedStep
	self.length = rand.IntN(mutantMaxLen-mutantMinLen) + mutantMinLen
	self.resetRunes(mutantCharset)
}

func (self *Drop) makeBackdrop() {
	self.back = true
	self.mutant = false
	self.speed = 1
	self.length = 5
	self.resetRunes(backCharset)
}

func (self *Drop) makeLucent() {
	self.makeNormal()
	self.lucent = true
}







func (self *Drop) reset() {
	self.pos = 0
	self.count = 0
	self.resetRunes(normalCharset)
}

func (self *Drop) resetRunes(n int) {
	for r := range self.runes {
		switch n {
			case -1:
				// self.runes[r] = '\u25AE'
				self.runes[r] = '|'
			case -2:
				self.runes[r] = '°'
				// self.runes[r] = '\u2591'
			case 0:
				self.runes[r] = rand.Int32N(27) + 65
			case 1:
				self.runes[r] = rand.Int32N(2) + 48
			case 2:
				self.runes[r] = rand.Int32N(93) + 33
			case 3:
				self.runes[r] = rand.Int32N(96) + 0x30A0
			case 4:
				self.runes[r] = rand.Int32N(96) + 0x3040
			case 5:
				self.runes[r] = rand.Int32N(144) + 0x0370
		}
	}
}


