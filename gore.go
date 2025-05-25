package gore

func Match(regex string, input string) bool {
	start := Compile(regex)
	return start.check(input, 0)
}

func Compile(regex string) *State {
	ctx := parse(regex)
	return toNfa(ctx)
}

func FindString(regex, input string) (string, bool) {
	nfa := Compile(regex)

	for start := 0; start <= len(input); start++ {
		if nfa.check(input[start:], 0) {
			// Walk forward to find the end of the match
			for end := len(input); end >= start; end-- {
				if nfa.check(input[start:end], 0) {
					return input[start:end], true
				}
			}
		}
	}
	return "", false
}

func FindAllString(regex, input string) []string {
	nfa := Compile(regex)
	var matches []string

	i := 0
	for i <= len(input) {
		if nfa.check(input[i:], 0) {
			for j := len(input); j >= i; j-- {
				if nfa.check(input[i:j], 0) {
					matches = append(matches, input[i:j])
					i = j // skip to the end of this match (non-overlapping)
					break
				}
			}
		} else {
			i++
		}
	}
	return matches
}

func FindStringIndex(regex, input string) (int, int) {
	nfa := Compile(regex)

	for start := 0; start <= len(input); start++ {
		if nfa.check(input[start:], 0) {
			for end := len(input); end >= start; end-- {
				if nfa.check(input[start:end], 0) {
					return start, end
				}
			}
		}
	}
	return -1, -1
}
