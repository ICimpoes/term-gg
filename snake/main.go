package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
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
	sync.Mutex
}

func (s *snake) drawSnake() {
	s.Lock()
	defer s.Unlock()
	for _, p := range s.points {
		termbox.SetCell(p.x, p.y, '☒', termbox.ColorRed, termbox.ColorDefault)
	}
}

var (
	N  int
	sx int
	sy int
)

func (s *snake) move() {
	s.Lock()
	defer s.Unlock()
	first := s.points[0]
	switch s.direction {
	case up:
		first.y--
	case down:
		first.y++
	case left:
		first.x -= 2
	case right:
		first.x += 2
	}
	if first.x == sx && first.y == sy {
		s.points = append([]point{first}, s.points...)
		sx = rand.Int()%N + 2
		if sx%2 != 0 {
			sx--
		}
		sy = rand.Int()%N + 1
		return
	}
	s.shift()
	s.points[0] = first
}
func (s *snake) shift() {
	l := len(s.points)
	for i := range s.points[:l-1] {
		s.points[l-i-1] = s.points[l-i-2]
	}
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().Unix())

	var s = snake{points: []point{{8, 1}, {6, 1}, {4, 1}, {2, 1}, {1, 1}}, direction: right}
	ch := time.Tick(100 * time.Millisecond)

	_, N = termbox.Size()
	N = N / 2
	sx = rand.Int()%N + 2
	if sx%2 != 0 {
		sx--
	}
	sy = rand.Int()%N + 1
	go func() {
		for {
			<-ch
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			s.drawSnake()
			drawBoard()
			drawStar()
			termbox.Flush()
			s.move()
		}
	}()

	for err := s.makeMove(); err == nil; {
		err = s.makeMove()
	}

	fmt.Println("\ncongrats!")
}

func drawStar() {
	termbox.SetCell(sx, sy, '❉', termbox.ColorRed, termbox.ColorDefault)
	//termbox.SetCursor(10, 10)
	//fmt.Println(sx, sy)
}

func drawBoard() {
	termbox.SetCell(0, 0, '┏', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(N*2, 0, '┓', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(0, N, '┗', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(N*2, N, '┛', termbox.ColorWhite, termbox.ColorBlack)
	for i := 1; i < N; i++ {
		termbox.SetCell(i*2, 0, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2+1, 0, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2, N, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2+1, N, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(0, i, '┃', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(N*2, i, '┃', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func (s *snake) makeMove() error {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return nil
	}
	switch ev.Key {
	case termbox.KeyArrowUp:
		s.direction = up
	case termbox.KeyArrowDown:
		s.direction = down
	case termbox.KeyArrowLeft:
		s.direction = left
	case termbox.KeyArrowRight:
		s.direction = right
	case termbox.KeyEsc:
		return errors.New("bye")
	}
	return nil
}
