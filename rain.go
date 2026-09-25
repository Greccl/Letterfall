package main

import (
	"math/rand/v2"
	// "fmt"
)



type SyncGroup struct {
	speed float64
	count float64
	advance int
}

var rainStatus bool

var cols []Column
var density float64 = 0.75 // amount of new drops per chunk of 10 columns
var spawnLeft float64
var spawnRatio float64
var normalSyncGroups []SyncGroup
var normalHead, normalNeck, normalTail Color
var normalMinLen, normalMaxLen int
var normalMinSpeed, normalMaxSpeed float64
var normalSpeedStep float64
var normalGroupCount int
var normalCharset int
var maxNormalsPerColumn int = 1

var mutantChance float32
var mutantCharset int
var mutantHead, mutantNeck, mutantTail Color
var mutantMinLen, mutantMaxLen int
var mutantMinSpeed, mutantMaxSpeed float64

var backs []Column
var backHead, backNeck, backTail Color
var backCharset int
var backSpeed float64 = 2.5
var backCount float64
var maxBacksPerColumn int = 3
var backDensity float64 = 0.75 // amount of new drops per chunk of 10 columns
var backSpawnLeft float64
var backSpawnRatio float64
var backLength int
var backAlphas []int32
var backdropRune rune = '•'
// var backdropRune rune = '\u2588'
const BACK_SEGMENTS int = 3

func setBackLength(l int) {
	if l % 2 == 0 { l++ }
	backLength = l
	c := l * BACK_SEGMENTS
	half := c / 2
	backAlphas = make([]int32, c)
	for i:=0 ; i < half; i++ {
		a := int32(1000.0 / float64(half) * float64(i))
		if a < 0 { a = 0 }
		if a > 1000 { a = 1000 }
		backAlphas[i] = a
		backAlphas[c - 1 - i] = a
	}
	backAlphas[half] = 1000
}

func rain_init() {
	addTickCallback(0.0, rain_tick)
	normalizeNormalSpeed()
	setBackLength(7)
}

func rain_resize() {
	spawnRatio = float64(scrw) / 10.0  * density
	backSpawnRatio = float64(scrw) / 10.0  * backDensity
	if scrw == 0 { return }
	cols = SliceResize(cols, scrw)
	backs = SliceResize(backs, scrw)
	
	for i := range cols {
		cols[i].x = i
		cols[i].layer = LAYER_RAIN
		cols[i].resize()
	}
	for i := range backs {
		backs[i].x = i
		backs[i].layer = LAYER_BACK
		backs[i].resize()
	}
}

func rain_tick(dt float64) {
	// Normal drops
	generator_0(dt)
	for i := range normalSyncGroups {
		g := &normalSyncGroups[i]
		g.count += g.speed * dt
		if g.count >= 1.0 {
			g.advance = int(g.count)
			g.count -= float64(g.advance)
		} else {
			g.advance = 0
		}
	}
	for i := 0; i < len(cols); i++ {
		cols[i].tick(dt)
	}

	// Back drops
	back_generator_1(dt)
	for i := 0; i < len(backs); i++ {
		backs[i].tick(dt)
	}

	// Update screen
	screen_damage()
}

func normalizeNormalSpeed() {
	if normalMaxSpeed < normalMinSpeed {
		normalMaxSpeed = normalMinSpeed
	}
	f := (normalMaxSpeed - normalMinSpeed) / normalSpeedStep
	n := int(f)
	if n < 0 { n = 0 }
	n++
	normalGroupCount = n
	if len(normalSyncGroups) != n {
		normalSyncGroups = make([]SyncGroup, n)
	}
	for i := range normalSyncGroups {
		normalSyncGroups[i].speed = normalMinSpeed + (float64(i) * normalSpeedStep)
	}
	for i := range cols {
		col := &cols[i]
		for j := range col.drops {
			drop := &col.drops[j]
			l := drop.length
			drop.makeNormal()
			drop.length = l
		}
	}
}

func generator_0(dt float64) {
	spawnLeft += spawnRatio * dt

	var mutantSpawned bool // 1 mutant per generator iteration

	for ; spawnLeft >= 1.0; spawnLeft -= 1.0 {
		// inital column number to be spawned
		x := rand.IntN((scrw+oddOffset)/2) * 2

		// find columns with fewer drops
		allowDups := rand.Float64() < 0.001
		if !allowDups {
			initialx := x
			delta := 2
			MAX: for max:=0; max<maxNormalsPerColumn; max++ {
				x = initialx
				for {
					if cols[x].count == max { break MAX }
					x = x + delta
					if x >= scrw { x = 0 }
					if x == initialx { break }
				}
			}
		}

		// discard if reached max drop amount in target column
		if cols[x].count >= maxNormalsPerColumn {
			continue
		}

		// las drop in target column
		var zero *Drop
		if cols[x].count > 0 {
			zero = &cols[x].drops[cols[x].count-1]
			// if zero.pos < overlap { continue }
		}
		
		// select drop type
		mut := false
		if !mutantSpawned {
			mut = rand.Float32() < mutantChance
		}

		// setup new drop
		d := cols[x].newDrop()
		d.reset()
		if mut {
			d.makeMutant()
			mutantSpawned = true
		} else {
			d.makeNormal()
			if zero != nil {
				d.pos = zero.pos - zero.length
			}
			if d.pos > 0 { d.pos = 0 }
		}
	}
}

var lastx int

func back_generator_1(dt float64) {
	backSpawnLeft += backSpawnRatio * dt
	if backSpawnLeft < 1.0 { return }

	var spawn bool

	var zero *Drop
	if backs[lastx].count == 0 {
		spawn = true
	} else {
		zero = &backs[lastx].drops[backs[lastx].count-1]
		if zero.pos > backLength + 3 {
			spawn = true
		}
	}

	if !spawn { return }
	
	d := backs[lastx].newDrop()
	d.reset()
	d.makeBackdrop()
	lastx++
	if lastx >= scrw { lastx = 0 }
	backSpawnLeft -= 1.0
}

func back_generator_0(dt float64) {
	backSpawnLeft += backSpawnRatio * dt
	
	for ; backSpawnLeft >= 1.0; backSpawnLeft -= 1.0 {
		x := rand.IntN(scrw)
		initialx := x
		MAX: for max:=0; max<maxBacksPerColumn; max++ {
			x = initialx
			for {
				if backs[x].count == max { break MAX }
				x++
				if x >= scrw { x = 0 }
				if x == initialx { break }
			}
		}

		if backs[x].count >= maxBacksPerColumn {
			return
		}

		var zero *Drop
		if backs[x].count > 0 {
			zero = &backs[x].drops[backs[x].count-1]
			// if zero.pos < overlap { return }
		}
		
		d := backs[x].newDrop()
		d.reset()
		d.makeBackdrop()
		if zero != nil {
			d.pos = zero.pos - zero.length - 3
		}
		if d.pos > 0 { d.pos = 0 }
	}
}
