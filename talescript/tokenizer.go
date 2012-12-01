package talescript

import (
	"io"
)

type state int
type TokenSet []int

type transition struct {
	symbol byte
	next   state
	tokens TokenSet
}

func createTokenizer(reader io.Reader) func(int) (int, Lexeme) {

	buffer := make(Lexeme, 0, 1)
	pos := 0
	state := state(0)

	reset := func() {
		copy(buffer, buffer[pos:])
		buffer = buffer[:len(buffer)-pos]
		pos = 0
		state = 0
	}

	read := func() (sym byte) {
		if pos >= len(buffer) {
			if len(buffer) == cap(buffer) {
				newlen := 2 * len(buffer)
				newbuffer := make(Lexeme, newlen)
				copy(newbuffer, buffer)
				buffer = newbuffer
			}
			n, err := reader.Read(buffer[pos:cap(buffer)])

			if n == 0 && err != nil {
				if err == io.EOF {
					return 0
				} else {
					panic(err)
				}
			}
			buffer = buffer[:pos+n]
		}
		sym = buffer[pos]
		pos++
		return
	}

	traverse := func(sym byte) TokenSet {
		trans := transitionTable[state].transitions
		lower, upper := 0, len(trans)
		if upper > 0 {
			for lower < upper-1 {
				index := (upper + lower) / 2
				if sym < trans[index].symbol {
					upper = index
				} else {
					lower = index
				}
			}

			state = trans[lower].next

			if sym >= trans[lower].symbol && state >= 0 {
				return trans[lower].tokens
			}
		}
		return nil
	}

	updateContext := func(tokens, context TokenSet) (TokenSet, bool) {
		for i, j := 0, 0; i < len(tokens); i++ {
			t := tokens[i]
			for ; j < len(context) && t > context[j]; j++ {
			}
			if t == context[j] {
				return context[:j+1], true
			}
		}

		potential := transitionTable[state].potential
		for i, j := 0, 0; i < len(context); i++ {
			t := context[i]
			for ; j < len(potential) && t > potential[j]; j++ {
			}
			if t == potential[j] {
				return context[i:], false
			}
		}

		return nil, false
	}

	return func(ictx int) (token int, str Lexeme) {

		context := tokenizerContexts[ictx]
		str = nil
		reset()

		for len(context) > 0 {
			tokens := traverse(read())
			if tokens == nil {
				break
			}
			var valid bool
			context, valid = updateContext(tokens, context)
			if valid {
				token = context[len(context)-1]
				str = buffer[:pos]
			}
		}

		pos = len(str)

		return
	}
}
