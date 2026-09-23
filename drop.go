package main

import (
	"math/rand/v2"
)

const (
	CHARSET_CUSTOM_A int = iota
	CHARSET_CUSTOM_B
	CHARSET_DEFAULT
	CHARSET_ONEZERO
	CHARSET_UPPERCASE
	CHARSET_ASCII
	CHARSET_HIRAGANA
	CHARSET_KATAKANA
	CHARSET_GREEK
	CHARSET_HEXAGRAM
	CHARSET_BRAILE
	CHARSET_CYRILLIC
	CHARSET_COUNT
	CHARSET_BACKDROP
)
var charsetNames = [CHARSET_COUNT]string {
	"customA", "customB", "default", "onezero", "uppercase", "ascii",
	"hiragana", "katakana", "greek", "hexagram", "braile", "cyrillic",
}

var backdropRune rune = '•'
var customCharsetA []rune
var customCharsetB []rune

type Drop struct {
	runes []rune
	mutant bool
	back bool

	pos int
	length int
	head int
	end int
	lastEnd int
	dirty bool

	group int
	speed float64
	count float64
}

func (self *Drop) makeNormal() {
	self.mutant = false
	self.group = rand.IntN(normalGroupCount)
	self.length = rand.IntN(normalMaxLen-normalMinLen) + normalMinLen
	self.resetRunes(normalCharset)
}

func (self *Drop) makeMutant() {
	self.mutant = true
	self.speed = mutantMinSpeed + rand.Float64() * (mutantMaxSpeed - mutantMinSpeed)
	self.length = rand.IntN(mutantMaxLen-mutantMinLen) + mutantMinLen
	self.resetRunes(mutantCharset)
}

func (self *Drop) makeBackdrop() {
	self.back = true
	self.mutant = false
	self.speed = 1
	self.length = 5
	self.resetRunes(-1)
}

func (self *Drop) reset() {
	self.pos = 0
	self.count = 0.0
	self.resetRunes(normalCharset)
}

func (self *Drop) resetRunes(n int) {
	if n == CHARSET_BACKDROP {
		for i := range self.runes {
			self.runes[i] = backdropRune
		}
		return
	} else if n == CHARSET_CUSTOM_A {
		if len(customCharsetA) > 0 {
			for i := range self.runes {
				j := rand.IntN(len(customCharsetA))
				self.runes[i] = customCharsetA[j]
			}
			return
		} else {
			n = CHARSET_DEFAULT
		}
	} else if n == CHARSET_CUSTOM_B {
		if len(customCharsetB) > 0 {
			for i := range self.runes {
				j := rand.IntN(len(customCharsetB))
				self.runes[i] = customCharsetB[j]
			}
			return
		} else {
			n = CHARSET_DEFAULT
		}
	}

	for i := range self.runes {
		switch n {
			case CHARSET_DEFAULT:
				self.runes[i] = '\u250B'
			case CHARSET_ONEZERO:
				self.runes[i] = 0x0030 + rand.Int32N(2)
			case CHARSET_UPPERCASE:
				self.runes[i] = 0x0041 + rand.Int32N(26)
			case CHARSET_ASCII:
				self.runes[i] = 0x0021 + rand.Int32N(94)
			case CHARSET_HIRAGANA:
				self.runes[i] = 0x3040 + rand.Int32N(96)
			case CHARSET_KATAKANA:
				self.runes[i] = 0x30A0 + rand.Int32N(96)
			case CHARSET_GREEK:
				r := 0x0391 + rand.Int32N(48)
				if r > 0x03A1 { r += 1 }
				if r > 0x03A8 { r += 7 }
				self.runes[i] = r
			case CHARSET_HEXAGRAM:
				self.runes[i] = 0x4DC0 + rand.Int32N(62)
			case CHARSET_BRAILE:
				self.runes[i] = 0x2801 + rand.Int32N(255)
			case CHARSET_CYRILLIC:
				self.runes[i] = 0x0410 + rand.Int32N(64)
		}
	}
}