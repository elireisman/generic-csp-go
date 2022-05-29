package main

import (
	"fmt"
	"github.com/elireisman/generic-csp-go/pkg/csp"
)

type Row int
type Column int

var (
	// CSP variables
	Queens []Row

	// CSP domains
	Columns []Column

	// CSP constraints
	Constraints []csp.Constraint[Row]
)

func NewQueen(row Row) csp.Constraint[Row] {
	return csp.Constraint[Row]{
		Variables: []Row{row},
	}
}

// constraint: ensure no newly-placed queen occupies a row and column
// that can be threatened by any other already-placed queen
func SatisfiesConstraint(queen csp.Constraint[Row], candidate map[Row]Column) bool {
	rowOccupied := queen.Variables[0]
	colOccupied, found := candidate[rowOccupied]

	// if no Queen has been assigned to the Row in the
	// Constraint yet, the Constraint is satisfied
	if !found {
		return true
	}

	// check if the proposed new queen satisfies the constraint that
	// no other already placed queen can threaten it
	return checkPosition(rowOccupied, colOccupied, candidate, 1, 0) &&
		checkPosition(rowOccupied, colOccupied, candidate, -1, 0) &&
		checkPosition(rowOccupied, colOccupied, candidate, 0, 1) &&
		checkPosition(rowOccupied, colOccupied, candidate, 0, -1) &&
		checkPosition(rowOccupied, colOccupied, candidate, 1, 1) &&
		checkPosition(rowOccupied, colOccupied, candidate, -1, -1) &&
		checkPosition(rowOccupied, colOccupied, candidate, 1, -1) &&
		checkPosition(rowOccupied, colOccupied, candidate, -1, 1)
}

func checkPosition(qRow Row, qCol Column, queens map[Row]Column, diffRow Row, diffCol Column) bool {
	qNextRow := qRow + diffRow
	qNextCol := qCol + diffCol

	if qNextRow < 1 || qNextRow > 8 {
		return true
	}
	if qNextCol < 1 || qNextCol > 8 {
		return true
	}

	qActualCol, otherQFound := queens[qNextRow]
	if !otherQFound || qActualCol != qNextCol {
		return checkPosition(qNextRow, qNextCol, queens, diffRow, diffCol)
	}

	return false
}

func init() {
	Queens = []Row{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
	}

	Columns = []Column{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
	}

	// TODO
	Constraints = []csp.Constraint[Row]{
		csp.Constraint[Row]{Variables: []Row{1}},
		csp.Constraint[Row]{Variables: []Row{2}},
		csp.Constraint[Row]{Variables: []Row{3}},
		csp.Constraint[Row]{Variables: []Row{4}},
		csp.Constraint[Row]{Variables: []Row{5}},
		csp.Constraint[Row]{Variables: []Row{6}},
		csp.Constraint[Row]{Variables: []Row{7}},
		csp.Constraint[Row]{Variables: []Row{8}},
	}
}

func renderBoard(result map[Row]Column) {
	for row := 1; row <= 8; row++ {
		for col := 1; col <= 8; col++ {
			elem := "\x1b[0;38m.\x1b[0;0m"
			if result[Row(row)] == Column(col) {
				elem = "\x1b[0;46mQ\x1b[0;0m"
			}
			fmt.Printf("%s ", elem)
		}
		fmt.Println()
	}
}

// model the 8 Queens problem using CSP framework + Go generics
func main() {
	// assemble mapping of variables to a set of possible
	// values to search for a valid solution
	domain := map[Row][]Column{}
	for _, q := range Queens {
		domain[q] = Columns
	}

	// create CSP framework instance, populate
	problem := csp.New(domain, SatisfiesConstraint)
	for _, takeable := range Constraints {
		problem.AddConstraint(takeable)
	}

	// init empty solution to begin search through problem space
	candidate := map[Row]Column{}

	// find ONE possible solution, and display it, if it exists
	if result := problem.Solve(candidate); result != nil {
		fmt.Println("Solution:")
		renderBoard(result)
		return
	}

	panic("No solution found")
}
