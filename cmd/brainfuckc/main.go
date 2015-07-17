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
	varScratch64Start = "VARIABLE ONE ZERO ZERO ZERO ZERO ONE ZERO"
	varUserStart      = "VARIABLE ONE ZERO ZERO ZERO ZERO ZERO ONE ZERO"

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
	n1 := w.reserve()
	n2 := w.reserve()

	// ptr++;
	w.assign(start, varDataPtr, w.offset(varDataPtr, 8), n1)
	// scratch = &index;
	w.assign(n1, varScratchPtr, addr+varDataIndexStart, n2)
	// *scratch++;
	w.increment(n2, 64, end, 0)
}

func (w *Writer) left(start, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()
	var n3 [64]uint64
	for i := range n3 {
		n3[i] = w.reserve()
	}
	n4 := w.reserve()
	n5 := w.reserve()
	n6 := w.reserve()

	// scratch = &index;
	w.assign(start, varScratchPtr, addr+varDataIndexStart, n1)

	// *scratch--;
	w.decrement(n1, 64, n2, 0)

	// scratch64 = index;
	for i := range n3 {
		w.assign(n2, deref+w.offset(addr+varScratch64Start, i), deref+w.offset(addr+varDataIndexStart, i), n3[i])
		n2 = n3[i]
	}

	// ptr = &array[0];
	w.assign(n2, varDataPtr, addr+varUserStart, n4)
	// scratch = &scratch64;
	w.assign(n4, varScratchPtr, addr+varScratch64Start, n5)
	// while (*scratch--)
	w.decrement(n5, 64, n6, end)
	// ptr++;
	w.assign(n6, varDataPtr, w.offset(varDataPtr, 8), n5)
}

func (w *Writer) output(start, end uint64) {
	var n1 [8]uint64
	for i := range n1[1:] {
		n1[i+1] = w.reserve()
	}
	n1[0] = end

	for i := 8 - 1; i >= 0; i-- {
		// write(*(ptr + i));
		w.outputBit(start, w.offset(varDataPtr, i), n1[i])
		start = n1[i]
	}
}

func (w *Writer) outputBit(start uint64, addr string, end uint64) {
	n1 := w.reserve()
	n2 := w.reserve()

	// jump = *addr;
	w.jump(start, deref+addr, n1, n2)
	// if (jump == 0) write(0);
	w.line(n1, "PRINT ZERO", end, end)
	// else write(1);
	w.line(n2, "PRINT ONE", end, end)
}

func (w *Writer) input(start, end uint64) {
	var n1 [8]uint64
	for i := range n1[1:] {
		n1[i+1] = w.reserve()
	}
	n1[0] = end

	for i := 8 - 1; i >= 0; i-- {
		// *(ptr + i) = read();
		w.inputBit(start, w.offset(varDataPtr, i), n1[i])
		start = n1[i]
	}
}

func (w *Writer) inputBit(start uint64, addr string, end uint64) {
	n1 := w.reserve()

	// jump = read();
	w.line(start, "READ", n1, n1)
	// *addr = jump;
	w.assign(n1, deref+addr, "THE JUMP REGISTER", end)
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

	for i := 0; i < bits; i++ {
		current := deref + w.offset(varScratchPtr, i)

		n1 := w.reserve()
		n2 := w.reserve()

		// if (*(scratch + i) == 0) goto n1; else goto n2;
		w.jump(next, current, n1, n2)
		// n1: *(scratch + i) = 1; goto end;
		w.assign(n1, current, "ONE", end)

		if i == bits-1 {
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

	for i := 0; i < bits; i++ {
		current := deref + w.offset(varScratchPtr, i)

		n1 := w.reserve()
		n2 := w.reserve()

		// if (*(scratch + i) == 0) goto n1; else goto n2;
		w.jump(next, current, n1, n2)

		if i == bits-1 {
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
