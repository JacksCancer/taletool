package talescript

import (
	"io"
	"fmt"
	"strings"
)

type Lexeme []byte
type terminalList []Lexeme
type parseraction func(p *parser, str Lexeme)
type parsermap map[int]parseraction
type ParserState int
type parser struct {
	listener ParserListener
	attribs  attributeStack
	stack    []ParserState
	unread   bool
}

type ParserError struct {
	token    int
	lex      Lexeme
	expected TokenSet
	stack    []ParserState
}

func (e ParserError) Error() string {
	expstr := make([]string, len(e.expected))
	for i, exp := range e.expected {
		expstr[i] = TokenNames[exp]
	}
	return fmt.Sprintf("error at <%s>: %s\n expected: %s\n",
		TokenNames[e.token],
		e.lex,
		strings.Join(expstr, ", "))
}

type ParserListener interface {
	Error(token int, str Lexeme, expected TokenSet, stack []ParserState)
	Shift(str Lexeme)
	AddVariable(identifier Lexeme)
	AddRule()
	NewLine()
	AddReference(identifier Lexeme)
	AddAction(op, identifier Lexeme)
	AddCondition(identifier Lexeme)
}

type attributeStack struct {
	identifier terminalList
	op         terminalList
}

func Parse(rd io.Reader, listener ParserListener) {
	tokenizer := createTokenizer(rd)
	parser := parser{
		listener: listener,
		attribs: attributeStack{
			identifier: make(terminalList, 0),
			op:         make(terminalList, 0),
		},
		stack:  make([]ParserState, 1),
		unread: false,
	}
	var token int
	var str Lexeme

ParserLoop:
	for len(parser.stack) > 0 {
		state := parsertable[parser.stack[len(parser.stack)-1]]
		// check for special case: single-reduce
		if len(state.parsermap) == 1 && state.parsermap[TokEND] != nil {
			unread := parser.unread
			state.parsermap[TokEND](&parser, nil)
			parser.unread = unread
			continue
		}

		if !parser.unread {
			token, str = tokenizer(state.context)
		} else {
			parser.unread = false
		}

		parser.traverse(token, str)
	}

	if token != TokEND {
		parser.stack = append(parser.stack, 0)
		goto ParserLoop
	}
}

func (p *parser) shift(state int, str Lexeme) {
	p.stack = append(p.stack, ParserState(state))
	if str != nil {
		p.listener.Shift(str)
	}
}

func (p *parser) push(state int, str Lexeme, lexemeList *terminalList) {
	s := make(Lexeme, len(str))
	copy(s, str)
	*lexemeList = append(*lexemeList, s)
	p.shift(state, str)
}

func (p *parser) traverse(token int, str Lexeme) {
	pmap := parsertable[p.stack[len(p.stack)-1]].parsermap
	action, ok := pmap[token]

	if !ok {
		// recheck for end token
		action, ok = pmap[TokEND]
	}

	if !ok {
		expected := make([]int, 0, len(pmap))
		for tok := range pmap {
			if tok != TokEND {
				expected = append(expected, tok)
			}
		}
		p.listener.Error(token, str, expected, p.stack)
	}

	action(p, str)
}

func (p *parser) reduce(head int, bodylen int) {
	p.remove(bodylen)
	p.traverse(head, nil)
}

func (p *parser) remove(bodylen int) {
	p.unread = true
	p.stack = p.stack[:len(p.stack)-bodylen]
}

func shift(nextstate int) parseraction {
	return func(p *parser, str Lexeme) {
		p.shift(nextstate, str)
	}
}

func push(nextstate int, lexemeList *terminalList) parseraction {
	return func(p *parser, str Lexeme) {
		s := make(Lexeme, len(str))
		copy(s, str)
		*lexemeList = append(*lexemeList, s)
		p.shift(nextstate, str)
	}
}

func reduce(head int, bodylen int) parseraction {
	if head == TokEND {
		return func(p *parser, str Lexeme) {
			p.remove(bodylen)
		}
	}
	return func(p *parser, str Lexeme) {
		p.reduce(head, bodylen)
	}
}

// rule -> _list_condition_comma condition colon _list_action_comma action
func reduce7(p *parser, str Lexeme) {
	p.listener.AddRule()
	p.reduce(TokRule, 5)
}

// condition -> identifier "(" _list_arg_comma arg ")"
func reduce10(p *parser, str Lexeme) {
	p.listener.AddCondition(p.attribs.identifier[len(p.attribs.identifier)-1])
	p.attribs.identifier = p.attribs.identifier[:len(p.attribs.identifier)-1]
	p.reduce(TokCondition, 5)
}

// arg -> identifier
func reduce13(p *parser, str Lexeme) {
	p.listener.AddVariable(p.attribs.identifier[len(p.attribs.identifier)-1])
	p.attribs.identifier = p.attribs.identifier[:len(p.attribs.identifier)-1]
	p.reduce(TokArg, 1)
}

// refarg -> identifier
func reduce14(p *parser, str Lexeme) {
	p.listener.AddReference(p.attribs.identifier[len(p.attribs.identifier)-1])
	p.attribs.identifier = p.attribs.identifier[:len(p.attribs.identifier)-1]
	p.reduce(TokRefarg, 1)
}

// END -> nl
func reduce15(p *parser, str Lexeme) {
	p.listener.NewLine()
	p.remove(1)
}

// action -> op identifier "(" _list_refarg_comma refarg ")"
func reduce17(p *parser, str Lexeme) {
	p.listener.AddAction(p.attribs.op[len(p.attribs.op)-1], p.attribs.identifier[len(p.attribs.identifier)-1])
	p.attribs.op = p.attribs.op[:len(p.attribs.op)-1]
	p.attribs.identifier = p.attribs.identifier[:len(p.attribs.identifier)-1]
	p.reduce(TokAction, 6)
}

// action -> identifier "(" _list_refarg_comma refarg ")"
func reduce18(p *parser, str Lexeme) {
	p.listener.AddAction(Lexeme{}, p.attribs.identifier[len(p.attribs.identifier)-1])
	p.attribs.identifier = p.attribs.identifier[:len(p.attribs.identifier)-1]
	p.reduce(TokAction, 5)
}

func pushIdentifier(state int) parseraction {
	return func(p *parser, str Lexeme) {
		p.push(state, str, &p.attribs.identifier)
	}
}
func pushOp(state int) parseraction {
	return func(p *parser, str Lexeme) {
		p.push(state, str, &p.attribs.op)
	}
}
