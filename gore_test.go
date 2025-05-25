package gore

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		regex   string
		input   string
		matches bool
	}{
		// Basic literals
		{"a", "a", true},
		{"a", "b", false},

		// Concatenation
		{"ab", "ab", true},
		{"abc", "abc", true},
		{"abc", "abcd", false},

		// Kleene star
		{"a*", "", true},
		{"a*", "a", true},
		{"a*", "aaaa", true},
		{"a*b", "aaab", true},
		{"a*b", "b", true},

		// Plus
		{"a+", "", false},
		{"a+", "a", true},
		{"a+", "aaaa", true},

		// Optional
		{"a?", "", true},
		{"a?", "a", true},
		{"a?", "aa", false},

		// Exact repetitions
		{"a{3}", "aaa", true},
		{"a{3}", "aa", false},

		// Ranged repetitions
		{"a{2,4}", "a", false},
		{"a{2,4}", "aa", true},
		{"a{2,4}", "aaaa", true},
		{"a{2,4}", "aaaaa", false},

		// Infinite upper bound
		{"a{2,}", "aa", true},
		{"a{2,}", "aaaaa", true},
		{"a{2,}", "a", false},

		// Character classes
		{"[abc]", "a", true},
		{"[a-c]", "b", true},
		{"[a-c]", "d", false},

		// Grouping
		{"(ab)*", "", true},
		{"(ab)*", "ab", true},
		{"(ab)*", "abab", true},
		{"(ab)*", "aba", false},

		// Alternation
		{"a|b", "a", true},
		{"a|b", "b", true},
		{"a|b", "c", false},
		{"ab|cd", "ab", true},
		{"ab|cd", "cd", true},
		{"ab|cd", "ad", false},
	}

	for _, test := range tests {
		t.Run(test.regex+"_"+test.input, func(t *testing.T) {
			got := Match(test.regex, test.input)
			if got != test.matches {
				t.Errorf("Match(%q, %q) = %v; want %v", test.regex, test.input, got, test.matches)
			}
		})
	}
}
