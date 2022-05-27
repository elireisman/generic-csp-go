package csp

// A Constraint represents state and business logic
// that models a single constraint to bevsatisfied
// while attempting to find a valid solution for a
// CSP. The constraint must be based on a valid
// variable binding in the problem space.
//
// A Constraint must be customized to fit the
// problem being solved:
// 1. Embed Constraint in another struct
// 2. Implement an override for the Satisfied method
type Constraint[V, D comparable] struct {
	Variables []V
}

// check if this iteration's candidate solution (set of assingments
// of legal variables to a single legal domain value each) violates
// this constraint on the problem space or not.
func (c Constraint[V, D]) Satisfied(candidate map[V]D) bool {
	panic("Abstract method: implement me!")
}
