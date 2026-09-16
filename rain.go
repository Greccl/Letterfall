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



func rain_resize() {
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

func tick_rain() {
	generator_0()
	if syncSpeed > 0 {
		syncCount++
		if syncCount >= syncSpeed {
			syncCount = 0
			syncAdvance = true
		} else {
			syncAdvance = false
		}
	} else if syncSpeed < 0 {
		syncCount++
	}
	for i := 0; i < len(cols); i++ {
		cols[i].tick()
	}
	damage()
}




/*
var gen1time int

func fgen(a, b float64) int {
	for i:=0; i<3; i++ {
		
	}
}

func f0(x, max int) int {
	sin := math.Sin
	xx := float64(x)
	// y := int(1*sin(1*xx) + 2*sin(0.5*xx) + 0.25*sin(4*xx))
	y := int(2*sin(xx/4))
	if y > max { y = max }
	return y
}

func f1(x, max int) int {
	sin := math.Sin
	xx := float64(x)
	// y := int(1*sin(1*xx) + 2*sin(0.5*xx) + 0.25*sin(4*xx))
	y := int(3*sin(xx/4))+2
	if y > max { y = max }
	return y
}

func generator_1() {
	if gen1time > 0 {
		gen1time--
		return
	}
	
	gen1time = 200
	boxh := 10
	speed := 10
	// var lasty int
	var change bool
	var f = f0
	var lap int
	for x := range cols {
		// y := 0 - rand.IntN(boxh)
		var y int
		if x > scrw/2 { change = true }
		if change {
			switch lap {
				case 0: f = f1
			}
			y = -5 + f(x, boxh)
			change = false
		} else {
			y = -5 + f(x, boxh)
		}
		// lasty = y
		d := cols[x].newDrop()
		d.reset()
		d.makeNormal()
		d.speed = speed
		d.pos = y
	}
}
*/

var genCounter int
var genFrameDuration int = 350
var density float32 = 1.0 // amount of new drops per chunk
var spawnLeft float32

func generator_0() {
	genCounter += frameDuration
	if genCounter < genFrameDuration { return }
	genCounter -= genFrameDuration

	chunks := float32(scrw) / 10.0
   // chunkSize := scrw / chunks
	newCount := chunks + spawnLeft
	newCount *= density

	// guard. 1 mutant per generator iteration
	var mutantSpawned bool

	// generation loop
	for ; newCount >= 1.0; newCount -= 1.0 {
		// inital column number to be spawned
		x := rand.IntN((scrw+oddOffset)/2) * 2

		// find columns with fewer drops
		allowDups := rand.Float32() < 0.001
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
				d.count = syncCount % d.speed
			}
		}
	}

	spawnLeft = newCount
	/*
	for i := 0; i < newCount; i++ {
		if rand.IntN(1000) < 25 { // 15
			x := rand.IntN(scrw/2)
			x *= 2
			x++
			if x >= scrw { continue }
			if cols[x].backs.speed > 0 {
				continue
			}
			cols[x].addBackDrop()
		}
	}
	*/
}


















var backSpeed int = 50
var backCount int

func tick_back() {
	back_generator_0()
	backCount++
	if backCount < backSpeed { return }
	backCount = 0
	for i := 0; i < len(backs); i++ {
		backs[i].tick()
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
