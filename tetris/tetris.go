package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

var (
	N int
)

type point struct {
	x, y int
}

type letter []point

type action int
type angle int

const (
	left action = iota
	right
	rotate

	a angle = iota
	b
	c
	d
)

func O(_ angle) letter {
	return letter{{0, 0}, {0, 1}, {2, 0}, {2, 1}}
}

func I(ang angle) letter {
	if ang == a || ang == c {
		return letter{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
	}
	return letter{{0, 0}, {2, 0}, {4, 0}, {6, 0}}
}

func Z(ang angle) letter {
	if ang == a || ang == c {
		return letter{{0, 0}, {2, 0}, {2, 1}, {4, 1}}
	}
	return letter{{0, 1}, {0, 2}, {2, 0}, {2, 1}}
}

func S(ang angle) letter {
	if ang == a || ang == c {
		return letter{{0, 1}, {2, 0}, {2, 1}, {4, 0}}
	}
	return letter{{0, 0}, {0, 1}, {2, 1}, {2, 2}}
}

func L(ang angle) letter {
	switch ang {
	case a:
		return letter{{0, 0}, {0, 1}, {0, 2}, {2, 2}}
	case b:
		return letter{{0, 2}, {2, 2}, {4, 1}, {4, 2}}
	case c:
		return letter{{0, 0}, {2, 0}, {2, 1}, {2, 2}}
	default:
		return letter{{0, 1}, {0, 2}, {2, 1}, {4, 1}}
	}
}

func J(ang angle) letter {
	switch ang {
	case a:
		return letter{{0, 2}, {2, 0}, {2, 1}, {2, 2}}
	case b:
		return letter{{0, 1}, {2, 1}, {4, 1}, {4, 2}}
	case c:
		return letter{{0, 0}, {0, 1}, {0, 2}, {2, 0}}
	default:
		return letter{{0, 1}, {0, 2}, {2, 2}, {4, 2}}
	}
}

func (l letter) drawAt(p point) {
	drawLetter(p, l)
}
func (l letter) at(at point) letter {
	lAt := make(letter, len(l))
	for i, p := range l {
		lAt[i] = point{p.x + at.x, p.y + at.y}
	}
	return lAt
}

type block struct {
	ang angle
	let func(angle) letter
}

func (b block) show() letter {
	return b.let(b.ang)
}

func (b *block) rotate() {
	b.ang++
	if b.ang > d {
		b.ang = a
	}
}

var board = make(map[point]bool)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	_, N = termbox.Size()
	N /= 2
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i, ang := range []angle{a, b, c, d} {
		for j, l := range []letter{O(ang), I(ang), Z(ang), L(ang), J(ang)} {
			drawLetter(point{j * 10, i * 10}, l)
		}
	}

	termbox.Flush()

	dir := make(chan action, 1)
	exit := make(chan struct{})
	go listenForAction(exit, dir)

	b := block{
		let: Z,
		ang: a,
	}
	i := 0
	ls := []func(angle) letter{O, I, Z, L, J, S}

	tick := time.Tick(500 * time.Millisecond)
	sP := point{N / 2, 0}
	for {
		select {
		case <-tick:
			newSp := sP
			newSp.y++
			if stop(b.show().at(newSp)) {
				for _, p := range b.show().at(sP) {
					board[p] = true
				}
				sP.y = 0
				b.let = ls[i]
				i++
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawBoard()
			b.show().drawAt(sP)
			termbox.Flush()
			sP.y++
		case act := <-dir:
			switch act {
			case left:
				sP.x -= 2
			case right:
				sP.x += 2
			case rotate:
				b.rotate()
			}
		case <-exit:
			return
		}
	}
}

func stop(l letter) bool {
	for _, p := range l {
		if board[p] ||
			p.x == 0 ||
			p.x == N*2 ||
			p.y == N {
			return true
		}
	}
	return false
}

func drawLetter(at point, l letter) {
	for _, p := range l {
		termbox.SetCell(p.x+at.x, p.y+at.y, '☒', termbox.ColorRed, termbox.ColorDefault)
	}
}

func drawBoard() {
	for i := 0; i < N; i++ {
		for _, j := range []int{0, N} {
			termbox.SetCell(i*2, j, '━', termbox.ColorWhite, termbox.ColorDefault)
			termbox.SetCell(i*2+1, j, '━', termbox.ColorWhite, termbox.ColorDefault)
			termbox.SetCell(2*j, i, '┃', termbox.ColorWhite, termbox.ColorDefault)
		}
	}
	termbox.SetCell(0, 0, '┏', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(N*2, 0, '┓', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, N, '┗', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(N*2, N, '┛', termbox.ColorWhite, termbox.ColorDefault)
	for p := range board {
		termbox.SetCell(p.x, p.y, '☒', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func listenForAction(done chan struct{}, ch chan action) {
	for {
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			continue
		}
		// clear channel buffer
		select {
		case <-ch:
		default:
		}

		switch ev.Key {
		case termbox.KeySpace:
			ch <- rotate
		case termbox.KeyArrowRight:
			ch <- right
		case termbox.KeyArrowLeft:
			ch <- left
		case termbox.KeyEsc:
			done <- struct{}{}
			break
		}
	}
}
