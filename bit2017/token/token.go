package token

// Token represents a word or pair of words in the BIT language.
type Token uint8

// Possible values for Token.
const (
	_ Token = iota

	Address
	At
	Beyond
	Close
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
	Address:      "ADDRESS",
	At:           "AT",
	Beyond:       "BEYOND",
	Close:        "CLOSE",
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
