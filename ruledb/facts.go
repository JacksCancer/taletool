package ruledb

type FactDb struct {
	facts map[string][]fact
}

type fact struct {
	predicate string
	args      []argument
}

type factIterator struct {
	facts []fact
}

func CreateFactDb() (db *FactDb) {
	db = new(FactDb)
	db.facts = make(map[string][]fact)
	return
}

func (iter *factIterator) next() (next *fact) {
	if len(iter.facts) > 0 {
		next = &iter.facts[0]
		iter.facts = iter.facts[1:]
	} else {
		next = nil
	}
	return
}

func (db *FactDb) iterate(predicate string) (iter *factIterator) {
	iter = new(factIterator)
	iter.facts = db.facts[predicate]
	return
}

func (f *fact) equals(other *fact) bool {
	if f.predicate == other.predicate && len(f.args) == len(other.args) {
		for i, a := range f.args {
			if a != other.args[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (db *FactDb) Add(other *fact) {
	facts := db.facts[other.predicate]
	for _, f := range facts {
		if f.equals(other) {
			return
		}
	}
	if facts != nil {
		db.facts[other.predicate] = append(facts, *other)
	} else {
		db.facts[other.predicate] = []fact{*other}
	}

}

func (db *FactDb) Delete(other *fact) {
	facts := db.facts[other.predicate]
	for _, f := range facts {
		if f.equals(other) {
			newlen := len(facts) - 1
			if newlen > 0 {
				f = facts[newlen]
				db.facts[other.predicate] = facts[:newlen]
			} else {
				delete(db.facts, other.predicate)
			}
			break
		}
	}
}
