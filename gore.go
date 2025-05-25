package gore

func Match(regex, input string) bool {
	ctx := parse(regex)
	nfa := toNfa(ctx)
	return nfa.check(input, 0)
}
