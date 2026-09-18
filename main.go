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



const (
	EVENT_EXIT int = iota
	EVENT_PLAY
	EVENT_STEP
)

type InternalEvent struct {
	tcell.EventTime
	event int
}

func postEvent(e int) {
	var ie InternalEvent
	ie.event = e
	scr.PostEvent(&ie)
}

var onRune func(rune)
var onKey func(tcell.Key)

func setKeyboardHandlers_default() {
	onRune = defaultRuneHandler
	onKey = defaultKeyHandler
}

func defaultRuneHandler(r rune) {
	
}

func defaultKeyHandler(k tcell.Key) {
	
}

func handleKey(ev *tcell.EventKey) {
	switch ev.Key() {
		case tcell.KeyRune:
			if onRune != nil {
				onRune(ev.Rune())
			}
		case tcell.KeyCtrlS:
			minishell_toggle()
		case tcell.KeyEscape, tcell.KeyCtrlQ, tcell.KeyCtrlC:
			postEvent(EVENT_EXIT)
		default:
			if onKey != nil {
				onKey(ev.Key())
			}
	}
}




type TickTask struct {
	interval float64
	elapsed float64
	runing bool
	callback func(float64)
}

var tickRegistry = make([]TickTask, 0)

func addTickCallback(interval float64, callback func(float64)) int {
	var task TickTask
	task.interval = interval
	task.callback = callback
	task.runing = true
	tickRegistry = append(tickRegistry, task)
	return len(tickRegistry) - 1
}

func setTickStatus(i int, status bool) {
	if i >= len(tickRegistry) { return }
	tickRegistry[i].runing = status
	tickRegistry[i].elapsed = 0.0
}

func resetTick(i int) {
	if i >= len(tickRegistry) { return }
	tickRegistry[i].elapsed = 0
}

func tickCycle(dt float64) {
	for i := range tickRegistry {
		t := &tickRegistry[i]
		if !t.runing { continue }
		t.elapsed += dt
		if t.elapsed >= t.interval {
			t.callback(t.elapsed)
			t.elapsed = 0.0
		}
	}
}


type ResizeListener func()

var resizeListeners []ResizeListener 

func addResizeListener(f ResizeListener) int {
	resizeListeners = append(resizeListeners, f)
	return len(resizeListeners)
}

func callResizeListeners() {
	for _, f := range resizeListeners {
		f()
	}
}

type HandleRequest struct {
	hnd *CommandHandler
	cmd *Command
}

var ch_HandleRequests = make(chan HandleRequest, 32)



func main() {
	loadDefaults()
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

	// Other inits
	minishell_init()
	rain_init()

	minishell_toggle()

	// Global ticker
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	last := time.Now()

	// Setup step, wait for resize event or abort
	// if a timeout is reached (is 1 second enough?)
	INIT: for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape {
							fmt.Println("initialization aborted")
							return
						}
					case *tcell.EventResize:
						resize()
						callResizeListeners()
						break INIT
				}
			case <- ticker.C:
		   	now := time.Now()
				if now.Sub(last).Seconds() >= 1.0 {
					fmt.Println("screen initialization tiemout reached")
					return
				}
				last = now
		}
	}

	// Main loop
	LOOP:	for {
		select {
			case ev := <- ch_ScreenEvents:
				switch ev := ev.(type) {
					case *tcell.EventKey:
						handleKey(ev)
					case *tcell.EventResize:
						resize()
						callResizeListeners()
					case *InternalEvent:
						switch ev.event {
							case EVENT_EXIT:
								break LOOP
							case EVENT_PLAY:
								rainStatus = !rainStatus
							case EVENT_STEP:
								if !rainStatus {
									tickCycle(0.25)
								}
						}
				}
			case req := <- ch_HandleRequests:
				req.hnd.do(req.cmd)
			case <- ch_Draw:
				scr.Show()
			case <- ticker.C:
		   	now := time.Now()
		   	dt := now.Sub(last)
		   	last = now
				if rainStatus {
					tickCycle(dt.Seconds())
				}
		}
	}
}


