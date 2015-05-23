package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
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

const (
	// variable zero - pointer, data pointer
	// variable one - pointer, scratch
	// variable one zero...one zero zero zero one zero - bits, data index
	// variable one zero zero zero one zero...one zero zero zero zero one zero - bits, scratch
	// variable one zero zero zero zero one zero - bit, first bit of data
	// variable one zero zero one zero one zero is the first bit of the second cell, and so on.

	varDataPtr        = "VARIABLE ZERO"
	varScratchPtr     = "VARIABLE ONE"
	varDataIndexStart = "VARIABLE ONE ZERO"
	varScratch64Start = "VARIABLE ONE ZERO ZERO ZERO ONE ZERO"
	varUserStart      = "VARIABLE ONE ZERO ZERO ZERO ZERO ONE ZERO"

	addr  = "THE ADDRESS OF "
	deref = "THE VALUE AT "
)

type Writer struct {
	w     *bufio.Writer
	n     uint64
	taken map[uint64]bool
}

func (w *Writer) Close() error {
	return w.w.Flush()
}

func (w *Writer) reserve() uint64 {
	n := w.n
	w.n++
	return n
}

func (w *Writer) number(n uint64) {
	if n == 0 {
		w.w.WriteString(" ZERO")
	} else {
		w.numberBits(n)
	}
}

func (w *Writer) numberBits(n uint64) {
	if n == 0 {
		return
	}

	w.numberBits(n >> 1)
	if n&1 == 0 {
		w.w.WriteString(" ZERO")
	} else {
		w.w.WriteString(" ONE")
	}
}

func (w *Writer) line(n uint64, line string, goto0, goto1 uint64) {
	if w.taken[n] {
		panic("INTERNAL COMPILER ERROR: DUPLICATE LINE NUMBER")
	}
	w.taken[n] = true
	w.w.WriteString("LINE NUMBER")
	w.number(n)
	w.w.WriteString(" CODE ")
	w.w.WriteString(line)

	// XXX: we don't allow jumping to line zero.
	if goto0 == goto1 {
		if goto0 != 0 {
			w.w.WriteString(" GOTO")
			w.number(goto0)
		}
	} else {
		if goto0 != 0 {
			w.w.WriteString(" GOTO")
			w.number(goto0)
			w.w.WriteString(" IF THE JUMP REGISTER IS ZERO")
		}
		if goto1 != 0 {
			w.w.WriteString(" GOTO")
			w.number(goto1)
			w.w.WriteString(" IF THE JUMP REGISTER IS ONE")
		}
	}
	w.w.WriteString("\n")
}

func (w *Writer) command(n uint64, cmd Command, done uint64) {
	switch cmd.Token {
	case Right:
		w.right(n, done)
	case Left:
		w.left(n, done)
	case Increment:
		next := w.reserve()
		// scratch = ptr;
		w.assign(n, varScratchPtr, varDataPtr, next)
		// *scratch++;
		w.increment(next, 8, done, done)
	case Decrement:
		next := w.reserve()
		// scratch = ptr;
		w.assign(n, varScratchPtr, varDataPtr, next)
		// *scratch--;
		w.decrement(next, 8, done, done)
	case Output:
		w.output(n, done)
	case Input:
		w.input(n, done)
	case Begin:
		w.loop(n, cmd.Loop, done)
	}
}

func (w *Writer) commands(list []Command, done uint64) uint64 {
	for i := len(list) - 1; i >= 0; i-- {
		n := w.reserve()

		w.command(n, list[i], done)

		done = n
	}

	return done
}

func (w *Writer) Program(list []Command) {
	var first uint64
	if len(list) != 0 {
		first = w.reserve()
	}

	// char *ptr = &array[0];
	w.assign(0, varDataPtr, addr+varUserStart, first)

	// // the rest of the program
	if len(list) != 0 {
		second := w.commands(list[1:], 0)
		w.command(first, list[0], second)
	}
}

func (w *Writer) jump(start uint64, value string, goto0, goto1 uint64) {
	// if (!value) goto goto0; else goto goto1;
	w.line(start, "THE JUMP REGISTER EQUALS "+value, goto0, goto1)
}

func (w *Writer) assign(start uint64, left, right string, end uint64) {
	// left = right;
	w.line(start, left+" EQUALS "+right, end, end)
}

func (w *Writer) right(start, end uint64) {
	n2 := w.reserve()
	n3 := w.reserve()

	// ptr++;
	w.assign(start, varDataPtr, w.offset(varDataPtr, 8), n2)
	// scratch = &index;
	w.assign(n2, varScratchPtr, addr+varDataIndexStart, n3)
	// *scratch++;
	w.increment(n3, 64, end, 0)
}

func (w *Writer) left(start, end uint64) {
	n2 := w.reserve()
	n3 := w.reserve()
	n4 := w.reserve()
	n5 := w.reserve()
	var n6 [64]uint64
	for i := range n6 {
		n6[i] = w.reserve()
	}
	n7 := w.reserve()
	n8 := w.reserve()
	n9 := w.reserve()

	// scratch = &index;
	w.assign(start, varScratchPtr, addr+varDataIndexStart, n2)
	// *scratch--;
	w.decrement(n2, 64, n3, 0)

	// // copy the current index to scratch64.
	// ptr = &scratch64;
	w.assign(n3, varDataPtr, addr+varScratch64Start, n4)
	// scratch = &index;
	w.assign(n4, varScratchPtr, addr+varDataIndexStart, n5)
	for i := 0; i < 64; i++ {
		// *(ptr + i) = *(scratch + i);
		w.assign(n5, deref+w.offset(varDataPtr, i), deref+w.offset(varScratchPtr, i), n6[i])
		n5 = n6[i]
	}

	// ptr = &array[0];
	w.assign(n5, varDataPtr, addr+varUserStart, n7)
	// scratch = &scratch64;
	w.assign(n7, varScratchPtr, addr+varScratch64Start, n8)
	// while (*scratch--)
	w.decrement(n8, 64, n9, end)
	// ptr++;
	w.assign(n9, varDataPtr, w.offset(varDataPtr, 8), n7)
}

func (w *Writer) output(start, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()
	n3 := w.reserve()
	n4 := w.reserve()
	n5 := w.reserve()
	n6 := w.reserve()
	n7 := w.reserve()
	n8 := w.reserve()

	// scratch = ptr;
	w.assign(start, varScratchPtr, varDataPtr, n1)
	// write(*scratch++);
	w.outputBit(n1, n2)
	// write(*scratch++);
	w.outputBit(n2, n3)
	// write(*scratch++);
	w.outputBit(n3, n4)
	// write(*scratch++);
	w.outputBit(n4, n5)
	// write(*scratch++);
	w.outputBit(n5, n6)
	// write(*scratch++);
	w.outputBit(n6, n7)
	// write(*scratch++);
	w.outputBit(n7, n8)
	// write(*scratch++);
	w.outputBit(n8, end)
}

func (w *Writer) outputBit(start, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()
	n3 := w.reserve()

	// jump = *scratch;
	w.jump(start, deref+varScratchPtr, n1, n2)
	// if (jump == 0) write(0);
	w.line(n1, "PRINT ZERO", n3, n3)
	// else write(1);
	w.line(n2, "PRINT ONE", n3, n3)
	// scratch++;
	w.assign(n3, varScratchPtr, w.offset(varScratchPtr, 1), end)
}

func (w *Writer) input(start, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()
	n3 := w.reserve()
	n4 := w.reserve()
	n5 := w.reserve()
	n6 := w.reserve()
	n7 := w.reserve()
	n8 := w.reserve()

	// scratch = ptr;
	w.assign(start, varScratchPtr, varDataPtr, n1)
	// *scratch++ = read();
	w.inputBit(n1, n2)
	// *scratch++ = read();
	w.inputBit(n2, n3)
	// *scratch++ = read();
	w.inputBit(n3, n4)
	// *scratch++ = read();
	w.inputBit(n4, n5)
	// *scratch++ = read();
	w.inputBit(n5, n6)
	// *scratch++ = read();
	w.inputBit(n6, n7)
	// *scratch++ = read();
	w.inputBit(n7, n8)
	// *scratch++ = read();
	w.inputBit(n8, end)
}

func (w *Writer) inputBit(start, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()

	// jump = read();
	w.line(start, "READ", n1, n1)
	// *scratch = jump;
	w.assign(n1, deref+varScratchPtr, "THE JUMP REGISTER", n2)
	// scratch++;
	w.assign(n2, varScratchPtr, w.offset(varScratchPtr, 1), end)
}

func (w *Writer) loop(start uint64, list []Command, end uint64) {
	inner := w.commands(list, start)
	n1 := w.reserve()
	n2 := w.reserve()
	n3 := w.reserve()
	n4 := w.reserve()
	n5 := w.reserve()
	n6 := w.reserve()
	n7 := w.reserve()

	// if (*ptr != 0) goto inner;
	w.jump(start, deref+varDataPtr, n1, inner)
	// if (*(ptr + 1) != 0) goto inner;
	w.jump(n1, deref+w.offset(varDataPtr, 1), n2, inner)
	// if (*(ptr + 2) != 0) goto inner;
	w.jump(n2, deref+w.offset(varDataPtr, 2), n3, inner)
	// if (*(ptr + 3) != 0) goto inner;
	w.jump(n3, deref+w.offset(varDataPtr, 3), n4, inner)
	// if (*(ptr + 4) != 0) goto inner;
	w.jump(n4, deref+w.offset(varDataPtr, 4), n5, inner)
	// if (*(ptr + 5) != 0) goto inner;
	w.jump(n5, deref+w.offset(varDataPtr, 5), n6, inner)
	// if (*(ptr + 6) != 0) goto inner;
	w.jump(n6, deref+w.offset(varDataPtr, 6), n7, inner)
	// if (*(ptr + 7) != 0) goto inner;
	w.jump(n7, deref+w.offset(varDataPtr, 7), end, inner)
}

func (w *Writer) offset(ptr string, n int) string {
	return strings.Repeat("OPEN PARENTHESIS THE ADDRESS OF OPEN PARENTHESIS THE VALUE BEYOND ", n) + "OPEN PARENTHESIS " + ptr + strings.Repeat(" CLOSE PARENTHESIS", 2*n+1)
}

func (w *Writer) increment(start uint64, bits int, end, overflow uint64) {
	next := start

	for i := bits - 1; i >= 0; i-- {
		current := deref + w.offset(varScratchPtr, i)

		n1 := w.reserve()
		n2 := w.reserve()

		// if (*(scratch + i) == 0) goto n1; else goto n2;
		w.jump(next, current, n1, n2)
		// n1: *(scratch + i) = 1; goto end;
		w.assign(n1, current, "ONE", end)

		if i == 0 {
			next = overflow
		} else {
			next = w.reserve()
		}
		// n2: *(scratch + i) = 0; goto next;
		w.assign(n2, current, "ZERO", next)
	}
}

func (w *Writer) decrement(start uint64, bits int, end, underflow uint64) {
	next := start

	for i := bits - 1; i >= 0; i-- {
		current := deref + w.offset(varScratchPtr, i)

		n1 := w.reserve()
		n2 := w.reserve()

		// if (*(scratch + i) == 0) goto n1; else goto n2;
		w.jump(next, current, n1, n2)

		if i == 0 {
			next = underflow
		} else {
			next = w.reserve()
		}
		// n1: *(scratch + i) = 1; goto next;
		w.assign(n1, current, "ONE", next)
		// n2: *(scratch + i) = 0; goto end;
		w.assign(n2, current, "ZERO", end)
	}
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     bufio.NewWriter(w),
		n:     1,
		taken: make(map[uint64]bool),
	}
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
