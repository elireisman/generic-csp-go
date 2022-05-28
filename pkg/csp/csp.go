package csp

import "fmt"

// Constraint models a single constraint to be satisfied
// while attempting to find a valid solution for a Problem.
// The constraint must be based on a valid variable declared
// in the problem setup.
//
// A Constraint must be customized to fit the problem being solved:
// 1. Write a function to construct properly-typed Constraint[V, D]s
// 2. Implement a Satisfied method to pass to the Problem[V, D] class
type Constraint[V comparable] struct {
	Variables []V
}

// check if this iteration's candidate solution (set of assingments
// of legal variables to a single legal domain value each) violates
// this constraint on the problem space or not.
//
// WORKAROUND for broken embedded struct + generics functionality :(
type Satisfied[V, D comparable] func(Constraint[V], map[V]D) bool

// Problem models a single instance of a constraint-satisfaction problem.
// Once instantiated and populated, this framework can produce a valid
// solution, if one exists.
type Problem[V, D comparable] struct {
	Domain      map[V][]D
	Constraints map[V][]Constraint[V]
	SatFn       Satisfied[V, D]
}

func New[V, D comparable](domain map[V][]D, satFn Satisfied[V, D]) Problem[V, D] {
	return Problem[V, D]{
		Domain:      domain,
		Constraints: map[V][]Constraint[V]{},
		SatFn:       satFn,
	}
}

// apply another constraint to the problem space
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

// backtracking recursive search through the problem space, attempting
// to fit a solution by testing one variable at a time through all
// related constraints until all variables are assigned to a candidate
// solution, meaning a complete, valid solution has been found.
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

	next := unassigned[0]
	for _, acceptableValues := range p.Domain[next] {
		// create a new candidate solution including the (novel)
		// next unassigned value from the input assignment
		candidate := dup(assignment)
		candidate[next] = acceptableValues
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
func dup[V, D comparable](assignment map[V]D) map[V]D {
	out := make(map[V]D, len(assignment))

	for k, v := range assignment {
		out[k] = v
	}

	return out
}
