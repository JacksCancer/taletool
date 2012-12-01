package ruledb

import (
	"os"
	ts "taletool/talescript"
	"testing"
	"fmt"
)

func prettyParameters(params []parameter) (s []byte) {
	for i, p := range params {
		if i > 0 {
			s = append(s, ", "...)
		}
		switch a := p.(type) {
		case *variable:
			s = append(append(s, "$"...), fmt.Sprint(a.ref+1)...)
		case *argument:
			s = append(s, a.string...)
		default:
			s = append(s, "unknown"...)
		}
	}
	return
}

func (c condition) pretty() []byte {
	return append(append(append(append(make([]byte, 0), c.predicate...), "("...), prettyParameters(c.params)...), ")"...)
}

func (a action) pretty() (s []byte) {
	switch a.op {
	case AddFact:
		s = append(s, "+"...)
	case DeleteFact:
		s = append(s, "-"...)
	}

	return append(append(append(append(s, a.predicate...), "("...), prettyParameters(a.params)...), ")"...)
}

func (r *rule) pretty() (s []byte) {
	for i, c := range r.conditions {
		if i > 0 {
			s = append(s, ", "...)
		}
		s = append(s, c.pretty()...)
	}

	s = append(s, ": "...)

	for i, a := range r.actions {
		if i > 0 {
			s = append(s, ", "...)
		}
		s = append(s, a.pretty()...)
	}
	return
}

func (db RuleDb) checkRule(r *rule, t *testing.T) {
	for _, c := range r.conditions {
		rules := db[c.predicate]
		if rules == nil {
			t.Error(db, "misses", c.predicate)
		} else {
			for _, s := range rules {
				if s.equals(r) {
					return
				}
			}
			t.Errorf("db[%s] misses '%s':", c.predicate, string(r.pretty()))

			for _, s := range rules {
				t.Error("\t\t", string(s.pretty()))
			}
		}
	}
}

func TestParser(t *testing.T) {
	file, err := os.Open("rules.ts")
	if err != nil {
		t.Error(err)
	} else {
		db := make(RuleDb)

		defer func() {
			e := recover()
			if e != nil {
				//err := e.(error)
				err := e.(ts.ParserError)
				t.Errorf(err.Error())
			}
		}()

		Parse(file, &db)

		db.checkRule(&rule{numvars: 1,
			conditions: []condition{
				condition{"take", []parameter{&variable{0}}},
				condition{"reach", []parameter{&variable{0}}}},
			actions: []action{
				action{AddFact, "inventar", []parameter{&variable{0}}},
			}}, t)

		db.checkRule(&rule{numvars: 1,
			conditions: []condition{
				condition{"drop", []parameter{&variable{0}}},
				condition{"inventar", []parameter{&variable{0}}}},
			actions: []action{
				action{DeleteFact, "inventar", []parameter{&variable{0}}},
			}}, t)

		file.Close()
	}
}
