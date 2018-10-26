package main

import (
	"bufio"
	"io"
	"unicode"
)

type stateMachine uint8

const (
	/*
	   .
	   O -> o.ne
	   Z -> z.ero
	*/
	state0 stateMachine = iota

	/*
	   z.ero
	   E -> ze.ro
	   O -> o.ne
	   Z -> z.ero
	*/
	state1

	/*
	   ze.ro
	   O -> o.ne
	   R -> zer.o
	   Z -> z.ero
	*/
	state2

	/*
	   zer.o
	   O -> [zero] -> .
	   Z -> z.ero
	*/
	state3

	/*
	   o.ne
	   N -> on.e
	   O -> o.ne
	   Z -> z.ero
	*/
	state4

	/*
	   on.e
	   E -> [one] -> .
	   O -> o.ne
	   Z -> z.ero
	*/
	state5
)

func convertOneZero(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)

	state := state0

	bitAccum := [2]byte{0, 8}
	defer func() {
		if bitAccum[1] != 8 {
			_, _ = w.Write(bitAccum[:1])
		}
	}()

	addBit := func(bit uint8) error {
		bitAccum[1]--
		bitAccum[0] |= bit << bitAccum[1]
		if bitAccum[1] != 0 {
			return nil
		}
		_, err := w.Write(bitAccum[:1])
		bitAccum[0] = 0
		bitAccum[1] = 8
		return err
	}

	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}

		if unicode.IsSpace(rune(b)) {
			continue
		}

		nextState := state0
		switch state {
		case state0:
			switch b {
			case 'O':
				nextState = state4
			case 'Z':
				nextState = state1
			}
		case state1:
			switch b {
			case 'E':
				nextState = state2
			case 'O':
				nextState = state4
			case 'Z':
				nextState = state1
			}
		case state2:
			switch b {
			case 'O':
				nextState = state4
			case 'R':
				nextState = state3
			case 'Z':
				nextState = state1
			}
		case state3:
			switch b {
			case 'O':
				err = addBit(0)
			case 'Z':
				nextState = state1
			}
		case state4:
			switch b {
			case 'N':
				nextState = state5
			case 'O':
				nextState = state4
			case 'Z':
				nextState = state1
			}
		case state5:
			switch b {
			case 'E':
				err = addBit(1)
			case 'O':
				nextState = state4
			case 'Z':
				nextState = state1
			}
		}

		if err != nil {
			return err
		}

		state = nextState
	}
}
