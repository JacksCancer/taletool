package ruledb

import (
	"errors"
	"fmt"
	"os"
	"strings"
	ts "taletool/talescript"
)

type ruleBuilder struct {
	db *RuleDb

	file   *os.File
	offset int64
	col    int
	line   int

	rule *rule

	args []parameter
	vars []string
}

func Parse(file *os.File, db *RuleDb) {
	offset, _ := file.Seek(0, 1)
	ts.Parse(file, &ruleBuilder{db: db, file: file, offset: offset, line: 1, rule: new(rule), args: make([]parameter, 0)})
}

func (s *ruleBuilder) Error(token int, str ts.Lexeme, expected ts.TokenSet, stack []ts.ParserState) {
	expstr := make([]string, len(expected))
	for i, exp := range expected {
		expstr[i] = ts.TokenNames[exp]
	}

	text := make([]byte, s.col+len(str))
	s.file.ReadAt(text, s.offset)

	panic(errors.New(fmt.Sprintf("%s:%d:%d: error at <%s>: '%s'\n%s\nexpected: %s\ncontext: %v",
		s.file.Name(),
		s.line,
		s.col,
		ts.TokenNames[token],
		str,
		text,
		strings.Join(expstr, ", "),
		stack)))
}

func (s *ruleBuilder) Shift(str ts.Lexeme) {
	fmt.Println("shift", string(str))
	s.col += len(str)
}

func (s *ruleBuilder) AddCondition(predicateIdentifier ts.Lexeme) {
	fmt.Println("condition", string(predicateIdentifier))
	s.rule.conditions = append(s.rule.conditions, condition{string(predicateIdentifier), s.args})
	s.args = make([]parameter, 0)
}

func (s *ruleBuilder) AddRule() {
	//fmt.Println("rule", string(predicateIdentifier))
	s.rule.numvars = len(s.vars)
	s.db.add(s.rule)
	fmt.Println("rule", s.rule)
	s.rule = new(rule)
	s.vars = s.vars[:0]
}

func (s *ruleBuilder) findVariable(identifier string) *variable {
	for i, v := range s.vars {
		if v == identifier {
			return &variable{ref: i}
		}
	}
	return nil
}

func (s *ruleBuilder) AddVariable(identifier ts.Lexeme) {
	id := string(identifier)
	v := s.findVariable(id)

	if v == nil {
		v = &variable{ref: len(s.vars)}
		s.vars = append(s.vars, id)
	}

	s.args = append(s.args, v)
}

func (s *ruleBuilder) AddReference(identifier ts.Lexeme) {
	id := string(identifier)
	v := s.findVariable(id)
	if v == nil {
		panic(errors.New(fmt.Sprintf("%s:%d:%d: unknown identifer '%s'\n",
			s.file.Name(), s.line, s.col, id)))
	}

	s.args = append(s.args, v)
}

func (s *ruleBuilder) NewLine() {
	s.offset += int64(s.col)
	s.col = 0
	s.line += 1
}

func (s *ruleBuilder) AddAction(opLex, predicateIdentifier ts.Lexeme) {
	var op int = SetFact
	switch string(opLex) {
	case "+":
		op = AddFact
	case "-":
		op = DeleteFact
	case "":
		op = SetFact
	default:
		panic("unknown op")
	}
	s.rule.actions = append(s.rule.actions, action{op, string(predicateIdentifier), s.args})
	s.args = make([]parameter, 0)
}
