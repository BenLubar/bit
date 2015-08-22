package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	"github.com/BenLubar/bit/bitgen"
)

// BF represents a single Brainfuck token.
type BF byte

// Definitions from Wikipedia.
const (
	// increment the data pointer (to point to the next cell to the right).
	Right BF = '>'
	// decrement the data pointer (to point to the next cell to the left).
	Left BF = '<'
	// increment (increase by one) the byte at the data pointer.
	Increment BF = '+'
	// decrement (decrease by one) the byte at the data pointer.
	Decrement BF = '-'
	// output the byte at the data pointer.
	Output BF = '.'
	// accept one byte of input, storing its value in the byte at the data
	// pointer.
	Input BF = ','
	// if the byte at the data pointer is zero, then instead of moving the
	// instruction pointer forward to the next command, jump it forward to
	// the command after the matching ] command.
	Begin BF = '['
	// if the byte at the data pointer is nonzero, then instead of moving
	// the instruction pointer forward to the next command, jump it back to
	// the command after the matching [ command.
	End BF = ']'
)

var (
	ErrUnmatchedBegin = errors.New("unmatched '['")
	ErrUnmatchedEnd   = errors.New("unmatched ']'")
)

type Command struct {
	// Token is never End
	Token BF

	// Loop is nil unless Token == Begin
	Loop []Command
}

// Tokenize returns two channels. The first channel will receive the tokens
// from the reader, skipping any non-BF characters. The first channel will be
// closed before the second channel receives. The second channel receives
// the error encountered or nil if the error was io.EOF.
func Tokenize(r io.Reader) (<-chan BF, <-chan error) {
	ch := make(chan BF)
	errch := make(chan error)

	go func() {
		br := bufio.NewReader(r)

		// If a [ is encountered as the first token, we skip until the
		// matching ].
		comment := 0

		for {
			c, err := br.ReadByte()
			if err != nil {
				if err == io.EOF {
					if comment <= 0 {
						err = nil
					} else {
						err = ErrUnmatchedBegin
					}
				}
				close(ch)
				errch <- err
				close(errch)
				return
			}

			bf := BF(c)

			switch bf {
			case Right, Left, Increment, Decrement, Input, Output:
				if comment == 0 {
					comment = -1
				}
				if comment == -1 {
					ch <- bf
				}

			case Begin:
				if comment == -1 {
					ch <- bf
				} else {
					comment++
				}

			case End:
				if comment == 0 {
					comment = -1
				}
				if comment == -1 {
					ch <- bf
				} else {
					comment--
				}
			}
		}
	}()

	return ch, errch
}

func Parse(tokens <-chan BF, err <-chan error) ([]Command, error) {
	return parse(tokens, err, 0)
}

func parse(tokens <-chan BF, err <-chan error, depth int) ([]Command, error) {
	var list []Command

	for tok := range tokens {
		switch tok {
		case Begin:
			l, e := parse(tokens, err, depth+1)
			list = append(list, Command{Token: Begin, Loop: l})
			if e != nil {
				return list, e
			}

		case End:
			if depth == 0 {
				for range tokens {
				}
				e := <-err
				if e == nil {
					e = ErrUnmatchedEnd
				}
				return list, e
			}
			return list, nil

		default:
			list = append(list, Command{Token: tok})
		}
	}

	e := <-err
	if e == nil && depth != 0 {
		e = ErrUnmatchedBegin
	}
	return list, e
}

type Writer struct {
	*bitgen.Writer
	DataPtr    bitgen.Variable
	ScratchPtr bitgen.Variable
	DataIndex  bitgen.Integer
	Scratch64  bitgen.Integer
	UserStart  bitgen.Variable
}

func (w *Writer) Command(n bitgen.Line, cmd Command, done bitgen.Line) {
	cell := bitgen.Integer{bitgen.ValueAt{w.DataPtr}, 8}

	switch cmd.Token {
	case Right:
		w.Right(n, done)
	case Left:
		w.Left(n, done)
	case Increment:
		w.Increment(n, cell, done, done)
	case Decrement:
		w.Decrement(n, cell, done, done)
	case Output:
		w.Output(n, cell, done)
	case Input:
		w.Input(n, cell, done)
	case Begin:
		w.Loop(n, cell, cmd.Loop, done)
	}
}

func (w *Writer) Commands(list []Command, done bitgen.Line) bitgen.Line {
	for i := len(list) - 1; i >= 0; i-- {
		n := w.ReserveLine()

		w.Command(n, list[i], done)

		done = n
	}

	return done
}

func (w *Writer) Program(list []Command) {
	var first bitgen.Line
	if len(list) != 0 {
		first = w.ReserveLine()
	}

	w.Assign(0, w.DataPtr, bitgen.AddressOf{w.UserStart}, first)

	if len(list) != 0 {
		second := w.Commands(list[1:], 0)
		w.Command(first, list[0], second)
	}
}

func (w *Writer) Right(start, end bitgen.Line) {
	n1 := w.ReserveLine()

	w.Assign(start, w.DataPtr, bitgen.Offset{w.DataPtr, 8}, n1)
	w.Increment(n1, w.DataIndex, end, 0)
}

func (w *Writer) Left(start, end bitgen.Line) {
	n1 := w.ReserveLine()
	n2 := w.ReserveLine()
	n3 := w.ReserveLine()
	n4 := w.ReserveLine()

	w.Decrement(start, w.DataIndex, n1, 0)
	w.Copy(n1, w.Scratch64, w.DataIndex, n2)
	w.Assign(n2, w.DataPtr, bitgen.AddressOf{w.UserStart}, n3)
	w.Decrement(n3, w.Scratch64, n4, end)
	w.Assign(n4, w.DataPtr, bitgen.Offset{w.DataPtr, 8}, n3)
}

func (w *Writer) Loop(start bitgen.Line, value bitgen.Integer, list []Command, end bitgen.Line) {
	inner := w.Commands(list, start)
	w.Cmp(start, value, 0, end, inner)
}

func NewWriter(writer io.Writer) *Writer {
	w := &Writer{Writer: bitgen.NewWriter(writer)}
	w.DataPtr = w.ReserveVariable()
	w.ScratchPtr = w.ReserveVariable()
	w.DataIndex = w.ReserveInteger(64)
	w.Scratch64 = w.ReserveInteger(64)
	w.UserStart = w.ReserveHeap()
	return w
}

func main() {
	list, err := Parse(Tokenize(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}

	w := NewWriter(os.Stdout)

	w.Program(list)

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
