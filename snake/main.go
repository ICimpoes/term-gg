package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type direction int

const (
	up = iota
	down
	left
	right
)

type point struct {
	x, y int
}

type snake struct {
	points    []point
	direction direction
	changeDir chan direction
}

func (s snake) drawSnake() {
	for _, p := range s.points {
		termbox.SetCell(p.x, p.y, '☒', termbox.ColorRed, termbox.ColorDefault)
	}
}

var (
	N    int
	star point
)

func (s *snake) move() {
	s.points = append(s.points[0:1], s.points[0:]...)
	select {
	case s.direction = <-s.changeDir:
	default:
	}
	switch s.direction {
	case up:
		s.points[0].y--
	case down:
		s.points[0].y++
	case left:
		s.points[0].x -= 2
	case right:
		s.points[0].x += 2
	}
	if s.points[0] == star {
		newStar()
		return
	}
	s.points = s.points[:len(s.points)-1]
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().Unix())

	var s = snake{points: []point{{4, 1}, {2, 1}, {1, 1}}, direction: right, changeDir: make(chan direction, 1)}

	_, N = termbox.Size()
	N /= 2
	newStar()
	exit := make(chan struct{})

	go sendDirection(exit, s.changeDir)

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			s.drawSnake()
			drawBoard()
			drawStar()
			showScore(len(s.points))
			termbox.Flush()
			s.move()
		case <-exit:
			fmt.Println("bye bye")
			return
		}
	}

}

func newStar() {
	rnd := func() int {
		return int(rand.Int31n(int32(N)-1) + 1)
	}
	x := rnd()
	if x%2 != 0 {
		x += x % (N - 1)
	}
	y := rnd()
	star = point{x: x, y: y}
}

func drawStar() {
	termbox.SetCell(star.x, star.y, '❉', termbox.ColorGreen, termbox.ColorDefault)
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
}

func showScore(score int) {
	for i, r := range fmt.Sprintf("Your score: %d", score) {
		termbox.SetCell((N/3)+i, N+2, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func sendDirection(done chan struct{}, ch chan direction) error {
	for {
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			continue
		}
		// clear channel buffer
		for len(ch) > 0 {
			<-ch
		}

		switch ev.Key {
		case termbox.KeyArrowUp:
			ch <- up
		case termbox.KeyArrowDown:
			ch <- down
		case termbox.KeyArrowLeft:
			ch <- left
		case termbox.KeyArrowRight:
			ch <- right
		case termbox.KeyEsc:
			done <- struct{}{}
		}
	}
}
