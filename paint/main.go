package main

import (
	"errors"

	"github.com/nsf/termbox-go"
)

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	for err = makeMove(); err == nil; {
		err = makeMove()
		termbox.Flush()
	}
}

func makeMove() error {
	switch ev := termbox.PollEvent(); ev.Key {
	case termbox.MouseLeft:
		termbox.SetCell(ev.MouseX, ev.MouseY, '▢', termbox.ColorRed, termbox.ColorDefault)
	case termbox.KeyCtrlC:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		return errors.New("stop")
	}
	return nil
}
