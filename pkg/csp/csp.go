package csp

type Variable interface {
	comparable
}

type Domain interface {
	comparable
}

// A CSP is state and a framework that models a single instance
// of a constraint-satisfaction problem. Once instantiated and
// populated, this framework can produce a valid solution if
// one exists.
type CSP[V Variable, D Domain] struct {
	Variables   []V
	Domains     map[V][]D
	Constraints map[V][]Constraint
}

func New(variables []V, domains map[V][]D) CSP {
	return CSP{
		Variables:   variables,
		Domains:     domains,
		Constraints: map[V][]Constraint{},
	}
}

// apply another constraint to the problem space
func (c CSP) AddConstraint(constraint Constraint) {
	for _, constraintVar := range constraint.Variables {
		// ensure each constraint var is part of the problem space
		found := false
		for _, cspVar := range c.Variables {
			if cspVar == constraintVar {
				found = true
				break
			}
		}
		if !found {
			panic("error: constraint variable %+v not found in CSP")
		}

		// store valid constraint
		c.Constraints[cspVar] = append(c.Constraints[cspVar], constraint)
	}
}

// backtracking recursive search through the problem space, attempting
// to fit a solution by testing one variable at a time through all
// related constraints until all variables are assigned to a candidate
// solution, meaning a complete and consistent solution has been found.
func (c CSP) Search(assignment map[V]D) map[V]D {
	// base case: all variables are assigned, a solution has been found
	if len(assignment) == len(c.Variables) {
		return assignment
	}

	// enumerate all currently-unassigned variables
	var unassigned []V
	for v := range c.Variables {
		if _, found := assignment[v]; !found {
			unassigned = append(unassigned, v)
		}
	}

	first := unassigned[0]
	for _, domains := range c.Domains[first] {
		// create a new candidate solution including the (novel)
		// first unassigned value from the input assignment
		candidate = dup(assignment)
		candidate[first] = domains
		// test if this augmented candidate assignment is still consistent
		if c.consistent(first, candidate) {
			result := c.Search(candidate)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

// determine if this variable and assignment satisfy the
// constraints applied to the problem space
func (c CSP) consistent(variable V, assignment map[V]D) bool {
	for _, constraint := range c.Constraints[variable] {
		if !constraint.Satisfied(assignment) {
			return false
		}
	}

	return true
}

// utility: copy the current assignment into a new map
func dup(currentAssignment map[V]D) map[V]D {
	out := make(map[V]D, len(assignment))

	for k, v := range assignment {
		out[k] = v
	}

	return out
}
