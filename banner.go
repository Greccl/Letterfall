package main

import (
	"github.com/spf13/pflag"
	"github.com/Greccl/tcell/v2"
)



type Banner struct {
	x, y int
	pix [4][3]rune
}

type Letter struct {
	
}

type Font struct {
	lets []Letter
}


func handleCommand_banner(fs *pflag.FlagSet) string {
	var b Banner
	var s tcell.Style

	b.x = 15
	b.y = 3
	// :▖:▗:▘:▙:▚:▛:▜:▝:▞:▟:▌:▐:▄:▀:
	b.pix[0] = [3]rune{' ','▄',' '}
	b.pix[1] = [3]rune{'▞',' ','▚'}
	b.pix[2] = [3]rune{'▙','▄','▟'}
	b.pix[3] = [3]rune{'▌',' ','▐'}

b.pix[0] = [3]rune{'⠀','⠤','⠀'}   // '▄' → '⠤' (línea media)
b.pix[1] = [3]rune{'⠜','⠀','⠣'}   // '▞' y '▚' → diagonales equivalentes
b.pix[2] = [3]rune{'⠟','⠤','⠻'}   // '▙','▄','▟' → parte inferior rellena
b.pix[3] = [3]rune{'⠸','⠀','⠸'}   // '▌' y '▐' → bordes laterales


	if true {
		for y := range b.pix {
			for x := range b.pix[y] {
				screen_drawCell(LAYER_TEXT, b.x+x, b.y+y, b.pix[y][x], s)
			}
		}
	}
	return ""
}
