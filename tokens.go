package gore

type tokenType uint8 // <1>

const ( // <2>
	group           tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	repeat          tokenType = iota
	literal         tokenType = iota
	groupUncaptured tokenType = iota
)

type token struct { // <3>
	tokenType tokenType
	// the payload required for each token will be different
	// so we need to be flexible with the type
	value interface{}
}

type parseContext struct { // <4>
	// the index of the character we're processing
	// in the regex string
	pos    int
	tokens []token
}
