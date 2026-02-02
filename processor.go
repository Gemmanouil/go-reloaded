package main

// ---------------------------
// FSM STATE DEFINITIONS
// ---------------------------

type State int

const (
	Normal State = iota
	MarkerStart
	MarkerBody
	MarkerEnd
	QuoteStart
	QuoteBody
	QuoteEnd
)

// ---------------------------
// PUBLIC ENTRY POINT
// ---------------------------

func ProcessLineFSM(input string) string {
	tokens := tokenize(input)

	state := Normal
	var out []string

	var markerBuf []string
	var quoteBuf []string

	for _, tok := range tokens {
		switch state {

		case Normal:
			next, outTok := fsmNormal(tok, &markerBuf, &quoteBuf)
			state = next
			if outTok != "" {
				out = append(out, outTok)
			}

		case MarkerStart:
			state = fsmMarkerStart(tok, &markerBuf)

		case MarkerBody:
			state = fsmMarkerBody(tok, &markerBuf)

		case MarkerEnd:
			next, newOut := fsmMarkerEnd(&markerBuf, out)
			state = next
			out = newOut

		case QuoteStart:
			state = fsmQuoteStart(tok, &quoteBuf)

		case QuoteBody:
			state = fsmQuoteBody(tok, &quoteBuf)

		case QuoteEnd:
			next, newOut := fsmQuoteEnd(&quoteBuf, out)
			state = next
			out = newOut
		}
	}

	return reconstruct(out)
}

// ---------------------------
// FSM HANDLERS (EMPTY STUBS)
// ---------------------------

func fsmNormal(tok string, markerBuf *[]string, quoteBuf *[]string) (State, string) {
	return Normal, tok
}

func fsmMarkerStart(tok string, markerBuf *[]string) State {
	return MarkerBody
}

func fsmMarkerBody(tok string, markerBuf *[]string) State {
	return MarkerEnd
}

func fsmMarkerEnd(markerBuf *[]string, out []string) (State, []string) {
	return Normal, out
}

func fsmQuoteStart(tok string, quoteBuf *[]string) State {
	return QuoteBody
}

func fsmQuoteBody(tok string, quoteBuf *[]string) State {
	return QuoteBody
}

func fsmQuoteEnd(quoteBuf *[]string, out []string) (State, []string) {
	return Normal, out
}

// ---------------------------
// TOKENIZER (TEMPORARY STUB)
// ---------------------------

func tokenize(input string) []string {
	return splitBySpaces(input)
}

func splitBySpaces(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ' ' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
		} else {
			cur += string(r)
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

// ---------------------------
// RECONSTRUCTION (TEMPORARY)
// ---------------------------

func reconstruct(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	out := tokens[0]
	for i := 1; i < len(tokens); i++ {
		out += " " + tokens[i]
	}
	return out
}
