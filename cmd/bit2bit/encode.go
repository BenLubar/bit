package main

import (
	"bufio"
	"io"
)

var nibbleBits = [16]string{
	"ZERO\nZERO\nZERO\nZERO\n",
	"ZERO\nZERO\nZERO\nONE\n",
	"ZERO\nZERO\nONE\nZERO\n",
	"ZERO\nZERO\nONE\nONE\n",
	"ZERO\nONE\nZERO\nZERO\n",
	"ZERO\nONE\nZERO\nONE\n",
	"ZERO\nONE\nONE\nZERO\n",
	"ZERO\nONE\nONE\nONE\n",
	"ONE\nZERO\nZERO\nZERO\n",
	"ONE\nZERO\nZERO\nONE\n",
	"ONE\nZERO\nONE\nZERO\n",
	"ONE\nZERO\nONE\nONE\n",
	"ONE\nONE\nZERO\nZERO\n",
	"ONE\nONE\nZERO\nONE\n",
	"ONE\nONE\nONE\nZERO\n",
	"ONE\nONE\nONE\nONE\n",
}

func convertBit(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)

	for {
		b, err := br.ReadByte()
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, nibbleBits[b>>4])
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, nibbleBits[b&0xf])
		if err != nil {
			return err
		}
	}
}
