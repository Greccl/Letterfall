package main

import (
	// "math"
	"math/rand/v2"
)





var cols []Column
var backs []Column
var syncCount int
var syncAdvance bool
var rainStatus bool


func rain_init() {
	addTickCallback(0.0, rain_tick)
	normalizeNormalSpeed()
}

func rain_resize() {
	spawnRatio = float64(scrw) / 10.0  * density
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
	screen_damage()
}





var density float64 = 1.0 // amount of new drops per chunk of 10 columns
var spawnLeft float64
var spawnRatio float64

func generator_0(dt float64) {
	spawnLeft += spawnRatio * dt

	// guard. 1 mutant per generator iteration
	var mutantSpawned bool

	// generation loop
	for ; spawnLeft >= 1.0; spawnLeft -= 1.0 {
		// inital column number to be spawned
		x := rand.IntN((scrw+oddOffset)/2) * 2

		// find columns with fewer drops
		allowDups := rand.Float64() < 0.001
		if !allowDups {
			initialx := x
			delta := 2
			MAX: for max:=0; max<maxDropsPerColumn; max++ {
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
		if cols[x].count >= maxDropsPerColumn {
			continue
		}

		// las drop in target column
		var zero *Drop
		if cols[x].count > 0 {
			zero = &cols[x].drops[cols[x].count-1]
			if zero.pos < overlap { continue }
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
			if syncSpeed < 0 {
				// d.count = syncCount % d.speed
			}
		}
	}
}


















var backSpeed int = 50
var backCount int

func tick_back(dt float64) {
	back_generator_0()
	backCount++
	if backCount < backSpeed { return }
	backCount = 0
	for i := 0; i < len(backs); i++ {
		backs[i].tick(dt)
	}
}

var backSpawnDelay int = 1200
var backSpawnAmount int = 1
var backSpawnCounter int
var maxBacksPerColumn int = 2

func back_generator_0() {
	backSpawnCounter += frameDuration
	if backSpawnCounter < backSpawnDelay { return }
	backSpawnCounter = 0

	for i:=0; i<backSpawnAmount; i++ {
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
			if zero.pos < overlap { return }
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
