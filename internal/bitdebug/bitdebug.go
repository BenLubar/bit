package bitdebug

import "github.com/BenLubar/bit/bitio"

var RunTrace func(program interface{}, in bitio.BitReader, out bitio.BitWriter, trace func(line, ctx interface{})) (interface{}, error)
var LineNum func(line interface{}) uint64
var LineStmt func(line interface{}) interface{}
var LineOpt func(line interface{}) interface{}
