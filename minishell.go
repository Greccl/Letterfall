package main

import (
	"github.com/Greccl/tcell/v2"
	"slices"
)

var mshActive bool
var mshBuffer []rune = make([]rune, 0, 256)
var mshSavedBuffer []rune
var mshEnabled bool
var mshStyle = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorBlack)
var mshStyleRev = tcell.StyleDefault.Background(tcell.ColorBrown).Foreground(tcell.ColorGray)
var mshHeight int = -1
var mshCursorPos int
var mshCursorBlink bool
var mshTickId int = -1
var mshResizeId int = -1
var mshCommandHandler = NewCommandHandler(SCOPE_SHELL)
var mshHistory = make([]string, 0)
var mshHistoryPos int = -1

func minishell_init() {
	mshTickId = addTickCallback(0.6, minishell_tickCallback)
	mshResizeId = addResizeListener(minishell_resizeCallback)
}

func minishell_tickCallback(dt float64) {
	mshCursorBlink = !mshCursorBlink
	minishell_draw()
}

func minishell_resizeCallback() {
	minishell_resize()
	minishell_draw()
}

func minishell_setString(str string) {
	minishell_clear()
	for _, r := range str {
		minishell_insert(r)
	}
}

func minishell_clear() {
	mshBuffer = mshBuffer[0:0]
	mshCursorPos = 0
}

func minishell_saveBuffer() {
	mshSavedBuffer = slices.Clone(mshBuffer)
}

func minishell_restoreBuffer() {
	mshBuffer = mshSavedBuffer
	mshSavedBuffer = nil
}

func minishell_insert(r rune) {
	mshBuffer = slices.Insert(mshBuffer, mshCursorPos, r)
	mshCursorPos++
}

func minishell_print(x, y int, i int, r rune) int {
	if i == mshCursorPos && mshCursorBlink {
		printCell(x, y, r, mshStyleRev)
	} else {
		printCell(x, y, r, mshStyle)
	}
	return i+1
}

func minishell_draw() {
	if scrw < 1 { return }
	if scrh < 1 { return }
	z := scrh
	h := mshHeight
	x := 0
	y := 1
	i := 0
	for x = 0; x<scrw; x++ {
		minishell_print(x, z, -1, '\u2580')
	}
	x = 0
	for ; y < h; y++ {
		for x = 0; x<scrw && i<len(mshBuffer); x++ {
			i = minishell_print(x, z+y, i, mshBuffer[i])
		}
	}
	y--	
	for ; x<scrw; x++ {
		i = minishell_print(x, z+y, i, ' ')
	}
}

func minishell_resize() {
	if scrw < 1 { return }
	if scrh < 1 { return }
	if !mshActive {
		mshHeight = 0
		reservedHeight = 0
		resize()
		return
	}
	h := len(mshBuffer) / scrw + 2 // one more line for top delimiter
	if h != mshHeight {
		mshHeight = h
		reservedHeight = h
		resize()
	}
	minishell_draw()
}

func minishell_toggle() {
	mshActive = !mshActive
	if mshActive {
		setKeyboardHandlers_minishell()
		mshBuffer = mshBuffer[:]
		mshCursorPos = 0
		mshHeight = 0
		minishell_draw()
	} else {
		setKeyboardHandlers_default()
	}
	minishell_resize()
}

func setKeyboardHandlers_minishell() {
	onRune = onRune_minishell
	onKey = onKey_minishell
}

func onRune_minishell(r rune) {
	minishell_insert(r)
	mshCursorBlink = true
	resetTick(mshTickId)
	minishell_resize()
	minishell_draw()
}

func onKey_minishell(k tcell.Key) {
	switch k {
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if mshCursorPos == 0 { break }
			mshBuffer = slices.Delete(mshBuffer, mshCursorPos - 1, mshCursorPos)
			mshCursorPos--
			minishell_resize()
		case tcell.KeyDelete:
			if mshCursorPos == len(mshBuffer) { break }
			mshBuffer = slices.Delete(mshBuffer, mshCursorPos, mshCursorPos + 1)
			minishell_resize()
		case tcell.KeyLeft:
			if mshCursorPos > 0 {
				mshCursorPos--
			}
		case tcell.KeyRight:
			if mshCursorPos < len(mshBuffer) {
				mshCursorPos++
			}
		case tcell.KeyHome:
			mshCursorPos = 0
		case tcell.KeyEnd:
			mshCursorPos = len(mshBuffer)
		case tcell.KeyEnter:
			line := string(mshBuffer)
			mshHistory = append(mshHistory, line)
			mshHistoryPos = -1
			mshSavedBuffer = nil
			minishell_clear()
			mshCommandHandler.eval(line)
		case tcell.KeyUp:
			if len(mshHistory) == 0 {
				return
			}
			if mshHistoryPos < len(mshHistory) - 1 {
				if mshHistoryPos == -1 {
					minishell_saveBuffer()
				}
				mshHistoryPos++
				minishell_setString(mshHistory[mshHistoryPos])
				mshCursorPos = len(mshBuffer)
			}
		case tcell.KeyDown:
			if mshHistoryPos >= 0 {
				mshHistoryPos--
				if mshHistoryPos == -1 {
					minishell_restoreBuffer()
				} else {
					minishell_setString(mshHistory[mshHistoryPos])
				}
				mshCursorPos = len(mshBuffer)
			}
		default:
			return
	}
	mshCursorBlink = true
	resetTick(mshTickId)
	minishell_draw()
}