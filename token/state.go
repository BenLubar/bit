package token

type state struct {
	rest   string
	choice *[26]state
	token  Token
}

var parser = state{
	choice: &[26]state{
		'A' - 'A': {
			choice: &[26]state{
				'D' - 'A': {
					rest:  "DRESS",
					token: Address,
				},
				'T' - 'A': {
					rest:  "",
					token: At,
				},
			},
		},
		'B' - 'A': {
			rest:  "EYOND",
			token: Beyond,
		},
		'C' - 'A': {
			choice: &[26]state{
				'L' - 'A': {
					rest:  "OSE",
					token: Close,
				},
				'O' - 'A': {
					rest:  "DE",
					token: Code,
				},
			},
		},
		'E' - 'A': {
			rest:  "QUALS",
			token: Equals,
		},
		'G' - 'A': {
			rest:  "OTO",
			token: Goto,
		},
		'I' - 'A': {
			choice: &[26]state{
				'F' - 'A': {
					rest:  "",
					token: If,
				},
				'S' - 'A': {
					rest:  "",
					token: Is,
				},
			},
		},
		'J' - 'A': {
			rest:  "UMP REGISTER",
			token: JumpRegister,
		},
		'L' - 'A': {
			rest:  "INE NUMBER",
			token: LineNumber,
		},
		'N' - 'A': {
			rest:  "AND",
			token: Nand,
		},
		'O' - 'A': {
			choice: &[26]state{
				'N' - 'A': {
					rest:  "E",
					token: One,
				},
				'P' - 'A': {
					rest:  "EN",
					token: Open,
				},
			},
		},
		'P' - 'A': {
			choice: &[26]state{
				'A' - 'A': {
					rest:  "RENTHESIS",
					token: Parenthesis,
				},
				'R' - 'A': {
					rest:  "INT",
					token: Print,
				},
			},
		},
		'R' - 'A': {
			rest:  "EAD",
			token: Read,
		},
		'T' - 'A': {
			rest:  "HE",
			token: The,
		},
		'V' - 'A': {
			rest: "A",
			choice: &[26]state{
				'L' - 'A': {
					rest:  "UE",
					token: Value,
				},
				'R' - 'A': {
					rest:  "IABLE",
					token: Variable,
				},
			},
		},
		'Z' - 'A': {
			rest:  "ERO",
			token: Zero,
		},
	},
}
