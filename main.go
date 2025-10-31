package main

import (
	"fmt"
	"os"
	"time"
	"github.com/Greccl/tcell/v2"
)





func SliceResize[T any](s []T, n int) []T {
	if n == len(s) { return s }
	if n <= cap(s) { return s[:n] }
	newSlice := make([]T, n)
	copy(newSlice, s)
	return newSlice
}

func SliceRemove[T any](s []T, i int) []T {
	last := len(s) - 1
	if i < last {
		copy(s[i:], s[i+1:])
	}
	return s[:last]
}




var ch_HandleRequests = make(chan HandleRequest, 32)





func main() {
	// Read command line arguments
	defaults()
	readCommandLine()

	// Init tcell screen
	var e error
	scr, e = tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = scr.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defer scr.Fini()

	ch_ScreenEvents := make(chan tcell.Event)

	go func() {
		for {
			ev := scr.PollEvent()
			ch_ScreenEvents <- ev
		}
	}()

	// A timer to update animations
	ch_Tick := time.Tick(time.Duration(frameDuration)*time.Millisecond)

	// Setup step, wait for resize event or abort
	// if a timeout is reached (is 1 second enough?)
	initTimeout := 0

	INIT:
	for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape {
							initTimeout = -2
							break INIT
						}
					case *tcell.EventResize:
						resize()
						initTimeout = 0
						break INIT
				}
			case <- ch_Tick:
				initTimeout += frameDuration
				if initTimeout >= 1000 {
					break INIT
				}
		}
	}
	
	if initTimeout != 0 {
		// Timeout reached or aborted by Esc key
		return
	}

	// Main loop
	LOOP:
	for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape { break LOOP }
						if ev.Key() == tcell.KeyCtrlQ { break LOOP }
						if ev.Key() == tcell.KeyRune {
							switch ev.Rune() {
								case 'p':
									rainStatus = !rainStatus
								case 's':
									if !rainStatus {
										text_tick()
										rain_tick()
										back_tick()
									}
							}
						}
					case *tcell.EventResize:
						resize()
				}
			case req := <- ch_HandleRequests:
				req.hnd.do(req.cmd)
			case <- ch_Draw:
				scr.Show()
			case <- ch_Tick:
				if rainStatus {
					text_tick()
					rain_tick()
					back_tick()
				}
		}
	}
}


