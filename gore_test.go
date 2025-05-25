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

func TestFindFunctions(t *testing.T) {
	t.Run("FindString", func(t *testing.T) {
		tests := []struct {
			regex, input, want string
			found              bool
		}{
			{"(ab)*", "xxababyy", "abab", true},
			{"a+", "aa bb aaa", "aa", true},
			{"z+", "xxxyyy", "", false},
		}

		for _, tt := range tests {
			got, ok := FindString(tt.regex, tt.input)
			if got != tt.want || ok != tt.found {
				t.Errorf("FindString(%q, %q) = %q, %v; want %q, %v",
					tt.regex, tt.input, got, ok, tt.want, tt.found)
			}
		}
	})

	t.Run("FindAllString", func(t *testing.T) {
		tests := []struct {
			regex, input string
			want         []string
		}{
			{"a+", "aa bb aaa c", []string{"aa", "aaa"}},
			{"b+", "aabbbccbb", []string{"bbb", "bb"}},
			{"x+", "abcdef", []string{}},
		}

		for _, tt := range tests {
			got := FindAllString(tt.regex, tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("FindAllString(%q, %q) = %v; want %v",
					tt.regex, tt.input, got, tt.want)
				continue
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("FindAllString(%q, %q)[%d] = %q; want %q",
						tt.regex, tt.input, i, got[i], tt.want[i])
				}
			}
		}
	})

	t.Run("FindStringIndex", func(t *testing.T) {
		tests := []struct {
			regex, input string
			wantStart    int
			wantEnd      int
		}{
			{"a+", "xxaaab", 2, 5},
			{"b+", "aaabbbbccc", 3, 7},
			{"z+", "abcdef", -1, -1},
		}

		for _, tt := range tests {
			start, end := FindStringIndex(tt.regex, tt.input)
			if start != tt.wantStart || end != tt.wantEnd {
				t.Errorf("FindStringIndex(%q, %q) = (%d, %d); want (%d, %d)",
					tt.regex, tt.input, start, end, tt.wantStart, tt.wantEnd)
			}
		}
	})
}
