package csp

// A Problem is state and a framework that models a single instance
// of a constraint-satisfaction problem. Once instantiated and
// populated, this framework can produce a (not every!) valid
// solution, if one exists.
type Problem[V, D comparable] struct {
	Variables   []V
	Domains     map[V][]D
	Constraints map[V][]Constraint[V, D]
}

func New[V, D comparable](variables []V, domains map[V][]D) Problem[V, D] {
	return Problem[V, D]{
		Variables:   variables,
		Domains:     domains,
		Constraints: map[V][]Constraint[V, D]{},
	}
}

// apply another constraint to the problem space
func (p Problem[V, D]) AddConstraint(constraint Constraint[V, D]) {
	for _, constraintVar := range constraint.Variables {
		// ensure each constraint var is part of the problem space
		found := false
		for _, cspVar := range p.Variables {
			if cspVar == constraintVar {
				found = true
				break
			}
		}
		if !found {
			panic("error: constraint variable %+v not found in Problem")
		}

		// store valid constraint
		p.Constraints[constraintVar] = append(p.Constraints[constraintVar], constraint)
	}
}

// backtracking recursive search through the problem space, attempting
// to fit a solution by testing one variable at a time through all
// related constraints until all variables are assigned to a candidate
// solution, meaning a complete, valid solution has been found.
func (p Problem[V, D]) Solve(assignment map[V]D) map[V]D {
	// base case: all variables are assigned, a solution has been found
	if len(assignment) == len(p.Variables) {
		return assignment
	}

	// enumerate all currently-unassigned variables
	var unassigned []V
	for _, v := range p.Variables {
		if _, found := assignment[v]; !found {
			unassigned = append(unassigned, v)
		}
	}

	next := unassigned[0]
	for _, domains := range p.Domains[next] {
		// create a new candidate solution including the (novel)
		// next unassigned value from the input assignment
		candidate := dup(assignment)
		candidate[next] = domains
		// test if this augmented candidate assignment is still consistent
		if p.consistent(next, candidate) {
			result := p.Solve(candidate)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

// determine if this variable and assignment satisfy the
// constraints applied to the problem space
func (p Problem[V, D]) consistent(variable V, assignment map[V]D) bool {
	for _, constraint := range p.Constraints[variable] {
		if !constraint.Satisfied(assignment) {
			return false
		}
	}

	return true
}

// utility: copy the current candidate solution into a new map
func dup[V, D comparable](assignment map[V]D) map[V]D {
	out := make(map[V]D, len(assignment))

	for k, v := range assignment {
		out[k] = v
	}

	return out
}
