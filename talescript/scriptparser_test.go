package talescript

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.Open("sample.ts")
	if err != nil {
		t.Error(err)
	} else {
		handler := newScriptHandler(file)

		defer func() {
			e := recover()
			if e != nil {
				//err := e.(error)
				err := e.(ParserError)
				t.Errorf(err.Error())
			}
		}()

		Parse(file, handler)

		file.Close()
	}
}

type scriptHandler struct {
	file *os.File

	offset int64
	col    int
	line   int

	args []Lexeme
}

func newScriptHandler(file *os.File) *scriptHandler {
	offset, _ := file.Seek(0, 1)
	return &scriptHandler{file: file, offset: offset, line: 1}
}

func (s *scriptHandler) Error(token int, str Lexeme, expected TokenSet, stack []ParserState) {
	expstr := make([]string, len(expected))
	for i, exp := range expected {
		expstr[i] = TokenNames[exp]
	}

	text := make([]byte, s.col+len(str))
	s.file.ReadAt(text, s.offset)

	panic(errors.New(fmt.Sprintf("%s:%d:%d: error at <%s>: '%s'\n%s\nexpected: %s\ncontext: %v",
		s.file.Name(),
		s.line,
		s.col,
		TokenNames[token],
		str,
		text,
		strings.Join(expstr, ", "),
		stack)))
}

func (s *scriptHandler) Shift(str Lexeme) {
	//fmt.Printf("shift %s\n", str)
	s.col += len(str)
}

func (s *scriptHandler) AddCondition(predicateIdentifier Lexeme) {
	fmt.Printf("condition %s\n", predicateIdentifier)
}

func (s *scriptHandler) AddRule() {
	fmt.Println("rule", s.line)
}

func (s *scriptHandler) AddVariable(identifier Lexeme) {
	fmt.Printf("identifier %s\n", identifier)
}

func (s *scriptHandler) AddReference(identifier Lexeme) {
	fmt.Printf("ref %s\n", identifier)
}

func (s *scriptHandler) NewLine() {
	s.offset += int64(s.col)
	s.col = 0
	s.line += 1
	fmt.Println("newline", s.line)
}

func (s *scriptHandler) AddAction(op, predicateIdentifier Lexeme) {
	fmt.Printf("action %s%s\n", op, predicateIdentifier)
}
