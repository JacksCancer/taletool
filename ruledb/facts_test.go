package ruledb

import (
	"testing"
)

func makeFact(predicate string, args ...string) (f *fact) {
	f = &fact{predicate, make([]argument, len(args))}
	for i, a := range args {
		f.args[i].string = a
	}
	return
}

func (f *fact) checkEquality(t *testing.T, predicate string, args ...string) {
	if !f.equals(makeFact(predicate, args...)) {
		t.Error("fact", *f, "does not equal", makeFact(predicate, args...))
	}
}

func (db *FactDb) checkSize(t *testing.T, size int) {
	if len(db.facts) != size {
		t.Error("size of", db.facts, "does not equal", size)
	}
}

func (db *FactDb) checkPredicateSize(t *testing.T, predicate string, size int) {
	if len(db.facts[predicate]) != size {
		t.Error("size of", db.facts[predicate], "does not equal", size)
	}
}

func TestFactDbAdd(t *testing.T) {
	db := CreateFactDb()
	db.checkSize(t, 0)

	db.Add(&fact{"P", []argument{}})
	db.checkSize(t, 1)
	db.facts["P"][0].checkEquality(t, "P")

	db.Add(&fact{"P", []argument{argument{"param1"}, argument{"param2"}}})
	db.checkSize(t, 1)
	db.facts["P"][0].checkEquality(t, "P")
	db.facts["P"][1].checkEquality(t, "P", "param1", "param2")

	//test double add
	db.Add(&fact{"Q", []argument{argument{"param1"}, argument{"param2"}}})
	db.Add(&fact{"Q", []argument{argument{"param1"}, argument{"param2"}}})
	db.checkSize(t, 2)
	db.facts["P"][0].checkEquality(t, "P")
	db.facts["P"][1].checkEquality(t, "P", "param1", "param2")
	db.facts["Q"][0].checkEquality(t, "Q", "param1", "param2")

	db.Add(&fact{"Q", []argument{argument{"param1"}, argument{"param3"}}})
	db.checkSize(t, 2)
	db.facts["P"][0].checkEquality(t, "P")
	db.facts["P"][1].checkEquality(t, "P", "param1", "param2")
	db.facts["Q"][0].checkEquality(t, "Q", "param1", "param2")
	db.facts["Q"][1].checkEquality(t, "Q", "param1", "param3")
}

func TestFactDbDelete(t *testing.T) {
	db := CreateFactDb()

	db.Add(&fact{"P", []argument{}})
	db.Add(&fact{"P", []argument{argument{"param1"}, argument{"param2"}}})
	db.Add(&fact{"Q", []argument{argument{"param1"}, argument{"param2"}}})
	db.Add(&fact{"Q", []argument{argument{"param1"}, argument{"param3"}}})
	db.Add(&fact{"Q", []argument{argument{"param1"}}})
	db.checkSize(t, 2)
	db.checkPredicateSize(t, "P", 2)
	db.checkPredicateSize(t, "Q", 3)
	db.facts["P"][0].checkEquality(t, "P")
	db.facts["P"][1].checkEquality(t, "P", "param1", "param2")
	db.facts["Q"][0].checkEquality(t, "Q", "param1", "param2")
	db.facts["Q"][1].checkEquality(t, "Q", "param1", "param3")
	db.facts["Q"][2].checkEquality(t, "Q", "param1")

	db.Delete(makeFact("P", "param1", "param2"))
	db.checkSize(t, 2)
	db.checkPredicateSize(t, "P", 1)
	db.checkPredicateSize(t, "Q", 3)
	db.facts["P"][0].checkEquality(t, "P")
	db.facts["Q"][0].checkEquality(t, "Q", "param1", "param2")
	db.facts["Q"][1].checkEquality(t, "Q", "param1", "param3")
	db.facts["Q"][2].checkEquality(t, "Q", "param1")

	db.Delete(makeFact("P"))
	db.checkPredicateSize(t, "P", 0)
	db.checkPredicateSize(t, "Q", 3)
	db.checkSize(t, 1)
	db.facts["Q"][0].checkEquality(t, "Q", "param1", "param2")
	db.facts["Q"][1].checkEquality(t, "Q", "param1", "param3")
	db.facts["Q"][2].checkEquality(t, "Q", "param1")
}
