package ruledb

type RuleDb map[string][]*rule

type rule struct {
	numvars    int
	conditions []condition
	actions    []action
}

type condition struct {
	predicate string
	params    []parameter
}

type context struct {
	args     []*argument
	bindings []int
}

const (
	DeleteFact = iota
	SetFact
	AddFact
)

type action struct {
	op        int
	predicate string
	params    []parameter
}

type parameter interface {
	unify(c *context, arg argument) bool
	deref(c *context) *argument
	equals(p parameter) bool
}

type variable struct {
	ref int
}

type argument struct {
	string
}

func (db *RuleDb) add(r *rule) {
	for _, c := range r.conditions {
		(*db)[c.predicate] = append((*db)[c.predicate], r)
	}
}

func (v *variable) unify(c *context, arg argument) bool {
	if c.args[v.ref] == nil {
		c.args[v.ref] = new(argument)
		*c.args[v.ref] = arg
		c.bindings = append(c.bindings, v.ref)
		return true
	}
	return *c.args[v.ref] == arg
}

func (v *variable) deref(c *context) *argument {
	return c.args[v.ref]
}

func (v *variable) equals(p parameter) bool {
	v2, ok := p.(*variable)
	return ok && *v == *v2
}

func (a *argument) unify(c *context, arg argument) bool {
	return *a == arg
}

func (a *argument) deref(c *context) *argument {
	return a
}

func (a *argument) equals(p parameter) bool {
	a2, ok := p.(*argument)
	return ok && *a == *a2
}

func (c *context) reset(n int) {
	for _, v := range c.bindings[n:] {
		c.args[v] = nil
	}
	c.bindings = c.bindings[:n]
}

func (r *rule) unify(ctx *context, icond int, db *FactDb) bool {
	if icond < len(r.conditions) {
		bindings := len(ctx.bindings)
		cond := r.conditions[icond]
		iter := db.iterate(cond.predicate)
	Facts:
		for f := iter.next(); f != nil; f = iter.next() {
			for i, p := range cond.params {
				if !p.unify(ctx, f.args[i]) {
					ctx.reset(bindings)
					continue Facts
				}
			}
			if r.unify(ctx, icond+1, db) {
				return true
			}
			ctx.reset(bindings)
		}
		return false
	}
	return true
}

func (a *action) makeFact(c *context) (f *fact) {
	f = &fact{a.predicate, make([]argument, len(a.params))}
	for i, p := range a.params {
		f.args[i] = *p.deref(c)
	}
	return
}

func (r *rule) apply(db *FactDb, tmp *FactDb) bool {
	ctx := context{
		args:     make([]*argument, r.numvars),
		bindings: make([]int, 0, r.numvars),
	}

	if r.unify(&ctx, 0, tmp) || r.unify(&ctx, 0, db) {
		for _, a := range r.actions {
			switch a.op {
			case AddFact:
				db.Add(a.makeFact(&ctx))
			case DeleteFact:
				db.Delete(a.makeFact(&ctx))
			case SetFact:
				tmp.Add(a.makeFact(&ctx))
			}
		}
		return true
	}
	return false
}

func (c *condition) equals(other condition) bool {
	if c.predicate != other.predicate || len(c.params) != len(other.params) {
		return false
	}

	for i, p := range c.params {
		if !p.equals(other.params[i]) {
			return false
		}

	}
	return true
}

func (a *action) equals(other action) bool {
	if a.predicate != other.predicate || a.op != other.op || len(a.params) != len(other.params) {
		return false
	}

	for i, p := range a.params {
		if !p.equals(other.params[i]) {
			return false
		}
	}
	return true
}

func (r *rule) equals(other *rule) bool {
	if r.numvars != other.numvars || len(r.conditions) != len(other.conditions) || len(r.actions) != len(other.actions) {
		return false
	}

	for i, c := range r.conditions {
		if !c.equals(other.conditions[i]) {
			return false
		}
	}

	for i, a := range r.actions {
		if !a.equals(other.actions[i]) {
			return false
		}
	}
	return true
}
