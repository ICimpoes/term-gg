package main

import (
	"fmt"
	"math/rand"
	"time"

	"sync"

	"errors"

	"github.com/nsf/termbox-go"
)

type direction int

const (
	up = iota
	down
	left
	right
)

type action int

const (
	move action = iota
	changeDirection
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

func (s *snake) move() {
	s.Lock()
	defer s.Unlock()
	s.shift()
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
}
func (s *snake) shift() {
	l := len(s.points)
	for i := range s.points[:l-1] {
		s.points[l-i-1] = s.points[l-i-2]
	}
}

func main() {

	var s = snake{points: []point{{8, 0}, {6, 0}, {4, 0}, {2, 0}, {0, 0}}, direction: right}
	ch := time.Tick(100 * time.Millisecond)

	go func() {
		for {
			<-ch
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			s.drawSnake()
			drawBoard()
			termbox.Flush()
			s.move()
		}
	}()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for err := s.makeMove(); err == nil; {
		err = s.makeMove()
	}

	rand.Seed(time.Now().Unix())

	fmt.Println("\ncongrats!")
}
func drawBoard() {
	_, n := termbox.Size()
	termbox.SetCell(0, 0, '┏', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(n*2, 0, '┓', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(0, n, '┗', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(n*2, n, '┛', termbox.ColorWhite, termbox.ColorBlack)
	for i := 1; i < n; i++ {
		termbox.SetCell(i*2, 0, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2+1, 0, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2, n, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i*2+1, n, '━', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(0, i, '┃', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(n*2, i, '┃', termbox.ColorWhite, termbox.ColorBlack)
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
