package ruledb

import (
	"testing"
)

func makeCondition(predicate string, args ...parameter) condition {
	return condition{predicate, args}
}

func makeAction(op int, predicate string, args ...parameter) action {
	return action{op, predicate, args}
}

func checkApply(r *rule, db *FactDb, tmp *FactDb, t *testing.T) {
	if !r.apply(db, tmp) {
		t.Error("failed to apply", *r)
	}
}

func TestUnify(t *testing.T) {
	db := CreateFactDb()
	tmp := CreateFactDb()

	db.checkSize(t, 0)

	tmp.Add(makeFact("take", "bollo"))
	tmp.checkSize(t, 1)
	tmp.facts["take"][0].checkEquality(t, "take", "bollo")

	r := &rule{
		numvars: 1,
		conditions: []condition{
			makeCondition("take", &variable{0}),
		},
		actions: []action{
			makeAction(AddFact, "inventar", &variable{0}),
		},
	}

	ctx := context{
		args:     make([]*argument, r.numvars),
		bindings: make([]int, 0, r.numvars),
	}

	if !r.unify(&ctx, 0, tmp) {
		t.Error("conditions of", r, "do not hold")
	}

	checkApply(r, db, tmp, t)

	db.checkSize(t, 1)
	db.facts["inventar"][0].checkEquality(t, "inventar", "bollo")

	tmp.checkSize(t, 1)
	tmp.facts["take"][0].checkEquality(t, "take", "bollo")

	r2 := &rule{
		numvars: 1,
		conditions: []condition{
			makeCondition("take", &variable{0}),
		},
		actions: []action{
			makeAction(SetFact, "hold", &variable{0}),
		},
	}

	checkApply(r2, db, tmp, t)

	db.checkSize(t, 1)
	db.facts["inventar"][0].checkEquality(t, "inventar", "bollo")

	tmp.checkSize(t, 2)
	tmp.facts["take"][0].checkEquality(t, "take", "bollo")
	tmp.facts["hold"][0].checkEquality(t, "hold", "bollo")
}
