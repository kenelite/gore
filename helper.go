package gore

const (
	startOfText uint8 = 1
	endOfText   uint8 = 2
)

func getChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return endOfText
	}

	if pos < 0 {
		return startOfText
	}

	return input[pos]
}

func (s *state) check(input string, pos int) bool { // <1>
	ch := getChar(input, pos) // <2>

	if ch == endOfText && s.terminal { // <3>
		return true
	}

	if states := s.transitions[ch]; len(states) > 0 { // <4>
		nextState := states[0]
		if nextState.check(input, pos+1) { // <5>
			return true
		}
	}

	for _, state := range s.transitions[epsilonChar] { // <6>
		if state.check(input, pos) { // <7>
			return true
		}

		if ch == startOfText && state.check(input, pos+1) { // <8>
			return true
		}
	}

	return false // <9>
}
