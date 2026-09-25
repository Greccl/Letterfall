package main

import (
	"math/rand/v2"
)

const (
	CHARSET_CUSTOM_A int = iota
	CHARSET_CUSTOM_B
	CHARSET_DEFAULT
	CHARSET_ONEZERO
	CHARSET_NUMBERS
	CHARSET_UPPERCASE
	CHARSET_ASCII
	CHARSET_ALNUM
	CHARSET_HIRAGANA
	CHARSET_KATAKANA
	CHARSET_HALFKANA
	CHARSET_GREEK
	CHARSET_HEXAGRAM
	CHARSET_BRAILE
	CHARSET_CYRILLIC
	CHARSET_COUNT
	CHARSET_BACKDROP
)
var charsetNames = [CHARSET_COUNT]string {
	"customA", "customB", "default", "onezero", "numbers", "uppercase", "ascii", "alnum",
	"hiragana", "katakana", "halfkana", "greek", "hexagram", "braile", "cyrillic",
}

const (
	CLASS_NORMAL int8 = iota
	CLASS_MUTANT
	CLASS_BACK
)

var customCharsetA []rune
var customCharsetB []rune

type Drop struct {
	runes []rune
	class int8

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
	self.class = CLASS_NORMAL
	self.group = rand.IntN(normalGroupCount)
	self.length = rand.IntN(normalMaxLen-normalMinLen) + normalMinLen
	self.resetRunes(normalCharset)
}

func (self *Drop) makeMutant() {
	self.class = CLASS_MUTANT
	self.speed = mutantMinSpeed + rand.Float64() * (mutantMaxSpeed - mutantMinSpeed)
	self.length = rand.IntN(mutantMaxLen-mutantMinLen) + mutantMinLen
	self.resetRunes(mutantCharset)
}

func (self *Drop) makeBackdrop() {
	self.class = CLASS_BACK
	self.length = backLength
	self.speed = backSpeed
	self.count = 0.0
	self.group = 0
	self.resetRunes(CHARSET_BACKDROP)
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
			case CHARSET_NUMBERS:
				self.runes[i] = 0x0030 + rand.Int32N(10)
			case CHARSET_UPPERCASE:
				self.runes[i] = 0x0041 + rand.Int32N(26)
			case CHARSET_ASCII:
				self.runes[i] = 0x0021 + rand.Int32N(94)
			case CHARSET_ALNUM:
				r := rand.Int32N(62)
				if r >= 36 {
					r += 0x0061 - 36
				} else if r >= 10 {
					r += 0x0041 - 10
				} else {
					r += 0x0030
				}
				self.runes[i] = r
		case CHARSET_HIRAGANA:
				self.runes[i] = 0x3040 + rand.Int32N(96)
			case CHARSET_KATAKANA:
				self.runes[i] = 0x30A0 + rand.Int32N(96)
			case CHARSET_HALFKANA:
				self.runes[i] = 0xFF66 + rand.Int32N(56)
			case CHARSET_GREEK:
				r := rand.Int32N(49)
				if r >= 24 {
					r += 0x03B1 - 24
				} else if r >= 17 {
					r += 0x03A3 - 17
				} else {
					r += 0x0391
				}
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