package main

import (
	"strconv"
	"strings"
	"unicode"
)

// ---------------------------
// FSM STATE DEFINITIONS
// ---------------------------

type State int

const (
	Normal State = iota
	QuoteBody
)

// ---------------------------
// PUBLIC ENTRY POINT
// ---------------------------

func ProcessLineFSM(input string) string {
	tokens := tokenize(input)

	state := Normal
	var out []string

	var quoteBuf []string

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		switch state {
		case Normal:
			if isMarkerToken(tok) {
				nextIsWord := false
				if i+1 < len(tokens) {
					nt := tokens[i+1]
					if !isMarkerToken(nt) && nt != "'" && !isPunct(nt) {
						nextIsWord = true
					}
				}
				applyMarker(tok, &out, nextIsWord)
				continue
			}

			if tok == "'" {
				quoteBuf = nil
				state = QuoteBody
				continue
			}

			if isPunct(tok) {
				if len(out) > 0 {
					out[len(out)-1] = out[len(out)-1] + tok
				} else {
					out = append(out, tok)
				}
				continue
			}

			out = append(out, tok)

		case QuoteBody:
			if tok == "'" {
				q := joinAndTrim(quoteBuf)
				out = append(out, "'"+q+"'")
				state = Normal
			} else {
				quoteBuf = append(quoteBuf, tok)
			}
		}
	}

	return reconstruct(out)
}

// ---------------------------
// TOKENIZER
// ---------------------------

func tokenize(input string) []string {
	var out []string
	runes := []rune(input)
	i := 0
	for i < len(runes) {
		r := runes[i]
		if unicode.IsSpace(r) {
			i++
			continue
		}

		if r == '(' {
			j := i
			for j < len(runes) && runes[j] != ')' {
				j++
			}
			if j < len(runes) && runes[j] == ')' {
				out = append(out, string(runes[i:j+1]))
				i = j + 1
				continue
			}
			// fallthrough to word handling if no closing paren
		}

		if r == '\'' {
			out = append(out, "'")
			i++
			continue
		}

		if isPunctRune(r) {
			j := i
			for j < len(runes) && isPunctRune(runes[j]) {
				j++
			}
			out = append(out, string(runes[i:j]))
			i = j
			continue
		}

		j := i
		for j < len(runes) && !unicode.IsSpace(runes[j]) && runes[j] != '(' && runes[j] != ')' && runes[j] != '\'' && !isPunctRune(runes[j]) {
			j++
		}
		out = append(out, string(runes[i:j]))
		i = j
	}
	return out
}

func isPunctRune(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';':
		return true
	}
	return false
}

func isPunct(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !isPunctRune(r) {
			return false
		}
	}
	return true
}

// ---------------------------
// MARKER / QUOTE / CASE HANDLING
// ---------------------------

func isMarkerToken(s string) bool {
	return len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')'
}

func applyMarker(marker string, out *[]string, nextIsWord bool) {
	inner := strings.TrimSpace(marker[1 : len(marker)-1])
	lower := strings.ToLower(inner)

	if lower == "hex" || lower == "bin" {
		if len(*out) == 0 {
			// no previous word: preserve marker literally per PRD
			*out = append(*out, marker)
			return
		}
		last := (*out)[len(*out)-1]
		// detach trailing punctuation
		trail := ""
		for len(last) > 0 && isPunctRune(rune(last[len(last)-1])) {
			trail = string(last[len(last)-1]) + trail
			last = last[:len(last)-1]
		}
		if last == "" {
			return
		}
		var v int64
		var err error
		if lower == "hex" {
			v, err = strconv.ParseInt(last, 16, 64)
		} else {
			v, err = strconv.ParseInt(last, 2, 64)
		}
		if err != nil {
			return
		}
		(*out)[len(*out)-1] = strconv.FormatInt(v, 10) + trail
		return
	}

	// case transforms: command [, count]
	parts := strings.Split(inner, ",")
	cmd := strings.ToLower(strings.TrimSpace(parts[0]))
	n := 1
	if len(parts) >= 2 {
		c := strings.TrimSpace(parts[1])
		if c == "" {
			return
		}
		vi, err := strconv.Atoi(c)
		if err != nil || vi <= 0 {
			return
		}
		n = vi
	}

	// if there are no previous words to act on, preserve marker literally
	if len(*out) == 0 {
		*out = append(*out, marker)
		return
	}

	switch cmd {
	case "up":
		applyCase(out, n, strings.ToUpper, nextIsWord)
	case "low":
		applyCase(out, n, strings.ToLower, nextIsWord)
	case "cap":
		applyCase(out, n, stringsToCap, nextIsWord)
	default:
		// unknown marker -> ignore
	}
}

func applyCase(out *[]string, n int, fn func(string) string, nextIsWord bool) {
	if n <= 0 {
		return
	}
	// choose starting index based on context
	i := len(*out) - 1
	if n == 1 && nextIsWord && len(*out) >= 2 {
		// marker is followed by a word: target second-to-last
		i = len(*out) - 2
	}
	applied := 0
	for i >= 0 && applied < n {
		s := (*out)[i]
		if isPunct(s) {
			i--
			continue
		}
		// separate trailing punctuation
		trail := ""
		for len(s) > 0 && isPunctRune(rune(s[len(s)-1])) {
			trail = string(s[len(s)-1]) + trail
			s = s[:len(s)-1]
		}
		s = fn(s)
		(*out)[i] = s + trail
		applied++
		i--
	}
}

func stringsToCap(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}

func joinAndTrim(parts []string) string {
	// collapse internal whitespace and trim
	combined := strings.Join(parts, " ")
	fields := strings.Fields(combined)
	return strings.Join(fields, " ")
}

// ---------------------------
// RECONSTRUCTION + ARTICLE RULE
// ---------------------------

func reconstruct(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}

	var out []string
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		// look ahead for article rule
		if (tok == "a" || tok == "A") && i+1 < len(tokens) {
			next := stripLeadingPunct(tokens[i+1])
			if next != "" {
				r0 := rune(next[0])
				if isVowelOrH(r0) {
					if tok == "A" {
						tok = "An"
					} else {
						tok = "an"
					}
				}
			}
		}
		out = append(out, tok)
	}

	return strings.Join(out, " ")
}

func stripLeadingPunct(s string) string {
	i := 0
	for i < len(s) && isPunctRune(rune(s[i])) {
		i++
	}
	return s[i:]
}

func isVowelOrH(r rune) bool {
	r = unicode.ToLower(r)
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'h'
}
