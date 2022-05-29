package csp

import "fmt"

// Constraint models a single constraint to be satisfied
// while attempting to find a valid solution for a Problem
type Constraint[V comparable] struct {
	Variables []V
}

// checks if the given Constraint is satisfied by the current candidate solution
type Satisfied[V comparable, D any] func(Constraint[V], map[V]D) bool

// Problem models a single instance of a constraint-satisfaction problem.
// Once instantiated and populated, it will brute-force a valid solution
type Problem[V comparable, D any] struct {
	Domain      map[V][]D
	Constraints map[V][]Constraint[V]
	SatFn       Satisfied[V, D]
}

// construct a Problem instance
func New[V comparable, D any](domain map[V][]D, satFn Satisfied[V, D]) Problem[V, D] {
	return Problem[V, D]{
		Domain:      domain,
		Constraints: map[V][]Constraint[V]{},
		SatFn:       satFn,
	}
}

// apply another Constraint to filter candidate solutions
func (p Problem[V, D]) AddConstraint(constraint Constraint[V]) {
	for _, constraintVar := range constraint.Variables {
		// ensure each constraint var is part of the problem space
		found := false
		for acceptableVar := range p.Domain {
			if acceptableVar == constraintVar {
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("error: constraint variable %+v not found in Problem", constraintVar))
		}

		// store valid constraint
		p.Constraints[constraintVar] = append(p.Constraints[constraintVar], constraint)
	}
}

// backtracking recursive search through the domain of problem
// variables and all their possible values. the first valid
// solution obtained in this brute-force effort is returned
func (p Problem[V, D]) Solve(assignment map[V]D) map[V]D {
	// base case: all variables are assigned, a solution has been found
	if len(assignment) == len(p.Domain) {
		return assignment
	}

	// enumerate all currently-unassigned variables
	var unassigned []V
	for acceptableVar := range p.Domain {
		if _, found := assignment[acceptableVar]; !found {
			unassigned = append(unassigned, acceptableVar)
		}
	}

	// test the current solution, augmented by the next
	// unassigned variable and a candidate value, against
	// all the constraints
	nextVar := unassigned[0]
	for _, candidateValue := range p.Domain[nextVar] {
		assignment[nextVar] = candidateValue
		if p.consistent(nextVar, assignment) {
			result := p.Solve(assignment)
			if result != nil {
				return result
			}
		} else {
			// the candidate value isn't a component of a
			// valid solution; ditch it and keep trying
			delete(assignment, nextVar)
		}
	}

	return nil
}

// determine if this variable and assignment satisfy the
// constraints applied to the problem space for that variable
func (p Problem[V, D]) consistent(variable V, assignment map[V]D) bool {
	for _, constraint := range p.Constraints[variable] {
		if !p.SatFn(constraint, assignment) {
			return false
		}
	}

	return true
}

// utility: copy the current candidate solution into a new map
func dup[V comparable, D any](assignment map[V]D) map[V]D {
	out := make(map[V]D, len(assignment))

	for k, v := range assignment {
		out[k] = v
	}

	return out
}
