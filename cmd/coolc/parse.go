package main

import "os"

func (ast *AST) ParseFile(name string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()

	panic("unimplemented")
}
