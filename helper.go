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

func (s *State) check(input string, pos int) bool {
	ch := getChar(input, pos)

	// If we're at the end of input and in a terminal state, return true
	if ch == endOfText && s.terminal {
		return true
	}

	// Explore all regular transitions (consume one character)
	if states, ok := s.transitions[ch]; ok {
		for _, next := range states {
			if next.check(input, pos+1) {
				return true
			}
		}
	}

	// Explore all epsilon transitions (consume no characters)
	if epsilonStates, ok := s.transitions[epsilonChar]; ok {
		for _, next := range epsilonStates {
			if next.check(input, pos) {
				return true
			}
		}
	}

	return false
}
