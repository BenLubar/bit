// Package token implements tokenization of BIT code.
package token

// Token represents a word or pair of words in the BIT language.
type Token uint8

// Possible values for Token.
const (
	_ Token = iota

	AddressOf
	At
	Beyond
	Close
	Code
	Equals
	Goto
	If
	Is
	JumpRegister
	LineNumber
	Nand
	One
	Open
	Parenthesis
	Print
	Read
	The
	Value
	Variable
	Zero
)

// String returns the token's source code representation, or "INVALID" if the
// token is not a valid Token constant from this package.
func (t Token) String() string {
	if t == 0 || int(t) > len(names) {
		return "INVALID"
	}

	return names[t]
}

var names = [...]string{
	AddressOf:    "ADDRESS OF",
	At:           "AT",
	Beyond:       "BEYOND",
	Close:        "CLOSE",
	Code:         "CODE",
	Equals:       "EQUALS",
	Goto:         "GOTO",
	If:           "IF",
	Is:           "IS",
	JumpRegister: "JUMP REGISTER",
	LineNumber:   "LINE NUMBER",
	Nand:         "NAND",
	One:          "ONE",
	Open:         "OPEN",
	Parenthesis:  "PARENTHESIS",
	Print:        "PRINT",
	Read:         "READ",
	The:          "THE",
	Value:        "VALUE",
	Variable:     "VARIBLE",
	Zero:         "ZERO",
}
