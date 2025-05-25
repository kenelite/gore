package gore

const (
	epsilonChar uint8 = 0
)

type State struct {
	start       bool
	terminal    bool
	transitions map[uint8][]*State
}

// empty character

func toNfa(ctx *parseContext) *State {
	if len(ctx.tokens) == 0 {
		// Special case: empty pattern
		start := &State{
			transitions: map[uint8][]*State{},
			start:       true,
		}
		end := &State{
			transitions: map[uint8][]*State{},
			terminal:    true,
		}
		start.transitions[epsilonChar] = []*State{end}
		return start
	}

	startState, endState := tokenToNfa(&ctx.tokens[0])

	for i := 1; i < len(ctx.tokens); i++ {
		nextStart, nextEnd := tokenToNfa(&ctx.tokens[i])

		// Add epsilon transition from current end to next start
		if endState.transitions == nil {
			endState.transitions = make(map[uint8][]*State)
		}
		endState.transitions[epsilonChar] = append(endState.transitions[epsilonChar], nextStart)

		endState = nextEnd
	}

	// Wrap with outer start and terminal end state
	outerStart := &State{
		transitions: map[uint8][]*State{
			epsilonChar: {startState},
		},
		start: true,
	}

	outerEnd := &State{
		transitions: map[uint8][]*State{},
		terminal:    true,
	}

	// Add epsilon transition from last real end to outerEnd
	if endState.transitions == nil {
		endState.transitions = make(map[uint8][]*State)
	}
	endState.transitions[epsilonChar] = append(endState.transitions[epsilonChar], outerEnd)

	return outerStart
}

func tokenToNfa(t *token) (*State, *State) {
	start := &State{
		transitions: map[uint8][]*State{},
	}
	end := &State{
		transitions: map[uint8][]*State{},
	}

	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*State{end}

	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		s1, e1 := tokenToNfa(&left)
		s2, e2 := tokenToNfa(&right)

		start.transitions[epsilonChar] = []*State{s1, s2}
		e1.transitions[epsilonChar] = []*State{end}
		e2.transitions[epsilonChar] = []*State{end}

	case bracket:
		literals := t.value.(map[uint8]bool)
		for l := range literals {
			start.transitions[l] = []*State{end}
		}

	case group, groupUncaptured:
		tokens := t.value.([]token)
		start, end = tokensToNfa(tokens)

	case repeat:
		p := t.value.(repeatPayload)

		if p.min == 0 {
			start.transitions[epsilonChar] = []*State{end}
		}

		var copyCount int
		if p.max == repeatInfinity {
			if p.min == 0 {
				copyCount = 1
			} else {
				copyCount = p.min
			}
		} else {
			copyCount = p.max
		}

		makeNfa := func(tok token) (*State, *State) {
			if tok.tokenType == group || tok.tokenType == groupUncaptured {
				toks := tok.value.([]token)
				return tokensToNfa(toks)
			}
			return tokenToNfa(&tok)
		}

		from, to := makeNfa(p.token)
		start.transitions[epsilonChar] = append(start.transitions[epsilonChar], from)

		for i := 2; i <= copyCount; i++ {
			s, e := makeNfa(p.token)
			to.transitions[epsilonChar] = append(to.transitions[epsilonChar], s)
			to = e
			if i > p.min {
				s.transitions[epsilonChar] = append(s.transitions[epsilonChar], end)
			}
		}

		to.transitions[epsilonChar] = append(to.transitions[epsilonChar], end)

		if p.max == repeatInfinity {
			end.transitions[epsilonChar] = append(end.transitions[epsilonChar], from)
		}

	default:
		panic("unknown token type")
	}

	return start, end
}

func tokensToNfa(tokens []token) (*State, *State) {
	start, end := tokenToNfa(&tokens[0])
	for i := 1; i < len(tokens); i++ {
		s, e := tokenToNfa(&tokens[i])
		end.transitions[epsilonChar] = append(end.transitions[epsilonChar], s)
		end = e
	}
	return start, end
}
