package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	board     [16]int
	emptyIndx int
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().Unix())
	prepareBoard()

	for !isSolved() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		showBoard()
		makeMove()
	}
	showBoard()
	fmt.Println("\ncongrats!")
}

func showBoard() {
	clearScreen()
	for i, v := range board {
		if i%4 == 0 {
			fmt.Println("\n --- --- --- ---")
		}
		s := "    "
		if v != 0 {
			s = fmt.Sprintf("%4d", v)
		}
		fmt.Print(s)
	}
}

func prepareBoard() {
	copy(board[:], rand.Perm(16))
	for i, v := range board {
		if v == 0 {
			emptyIndx = i
			return
		}
	}
}

func isSolved() bool {
	for i, v := range board[:15] {
		if v-i != 1 {
			return false
		}
	}
	return true
}

func makeMove() {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return
	}
	old := emptyIndx
	switch ev.Key {
	case termbox.KeyArrowUp:
		if emptyIndx-4 >= 0 {
			emptyIndx -= 4
		}
	case termbox.KeyArrowDown:
		if emptyIndx+4 <= 15 {
			emptyIndx += 4
		}
	case termbox.KeyArrowLeft:
		if emptyIndx%4 != 0 {
			emptyIndx -= 1
		}
	case termbox.KeyArrowRight:
		if emptyIndx%4 != 3 {
			emptyIndx += 1
		}
	case termbox.KeyEsc:
		fmt.Println("Bye!")
		os.Exit(0)
	}
	board[old] = board[emptyIndx]
	board[emptyIndx] = 0
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//
//type Runner interface {
//	Run() error
//}
//
//type Go struct {}
//
//func (Go) Run() error {
//	currentPos += 1
//	return checkIfOk()
//}
//
//type For struct {
//	i int
//	r Runner
//}
//
//func (f For) Run() error {
//	for i := 0; i<f.i; i++ {
//		if err := f.r.Run(); err != nil {
//			return nil
//		}
//	}
//	return nil
//}
//
//func parseCMDs(cmd string) {
//}
//
//func run(rs ...Runner) error {
//
//}
//
//func checkIfOk() error {
//	return nil
//}
