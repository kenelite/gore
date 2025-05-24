package gore

const (
	epsilonChar uint8 = 0
)

type state struct {
	start       bool
	terminal    bool
	transitions map[uint8][]*state
}

// empty character

func toNfa(ctx *parseContext) *state {
	startState, endState := tokenToNfa(&ctx.tokens[0]) // <1>

	for i := 1; i < len(ctx.tokens); i++ { // <2>
		startNext, endNext := tokenToNfa(&ctx.tokens[i]) // <3>
		endState.transitions[epsilonChar] = append(
			endState.transitions[epsilonChar],
			startNext,
		)                  // <4>
		endState = endNext // <5>
	}

	start := &state{ // <6>
		transitions: map[uint8][]*state{
			epsilonChar: {startState},
		},
		start: true,
	}
	end := &state{ // 7
		transitions: map[uint8][]*state{},
		terminal:    true,
	}

	endState.transitions[epsilonChar] = append(
		endState.transitions[epsilonChar],
		end,
	) // <8>

	return start // <9>
}

func tokenToNfa(t *token) (*state, *state) {
	start := &state{
		transitions: map[uint8][]*state{},
	}
	end := &state{
		transitions: map[uint8][]*state{},
	}

	switch t.tokenType {
	case literal:
		ch := t.value.(uint8)
		start.transitions[ch] = []*state{end}
		
	case or:
		values := t.value.([]token)
		left := values[0]
		right := values[1]

		s1, e1 := tokenToNfa(&left)  // <1>
		s2, e2 := tokenToNfa(&right) // <1>

		start.transitions[epsilonChar] = []*state{s1, s2} // <2>
		e1.transitions[epsilonChar] = []*state{end}       // <3>
		e2.transitions[epsilonChar] = []*state{end}       // <3>

	case bracket:
		literals := t.value.(map[uint8]bool)

		for l := range literals { // <1>
			start.transitions[l] = []*state{end} // <2>
		}

	case group, groupUncaptured:
		tokens := t.value.([]token)
		start, end = tokenToNfa(&tokens[0])
		for i := 1; i < len(tokens); i++ {
			ts, te := tokenToNfa(&tokens[i])
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				ts,
			)
			end = te
		}
	case repeat:
		p := t.value.(repeatPayload)

		if p.min == 0 { // <1>
			start.transitions[epsilonChar] = []*state{end}
		}

		var copyCount int // <2>

		if p.max == repeatInfinity {
			if p.min == 0 {
				copyCount = 1
			} else {
				copyCount = p.min
			}
		} else {
			copyCount = p.max
		}

		from, to := tokenToNfa(&p.token) // <3>
		start.transitions[epsilonChar] = append( // <4>
			start.transitions[epsilonChar],
			from,
		)

		for i := 2; i <= copyCount; i++ { // <5>
			s, e := tokenToNfa(&p.token)

			// connect the end of the previous one
			// to the start of this one
			to.transitions[epsilonChar] = append( // <6>
				to.transitions[epsilonChar],
				s,
			)

			// keep track of the previous NFA's entry and exit states
			from = s // <7>
			to = e   // <7>

			// after the minimum required amount of repetitions
			// the rest must be optional, thus we add an
			// epsilon transition to the start of each NFA
			// so that we can skip them if needed
			if i > p.min { // <8>
				s.transitions[epsilonChar] = append(
					s.transitions[epsilonChar],
					end,
				)
			}
		}

		to.transitions[epsilonChar] = append( // <9>
			to.transitions[epsilonChar],
			end,
		)

		if p.max == repeatInfinity { // <10>
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				from,
			)
		}
	default:
		panic("unknown type of token")
	}

	return start, end
}
