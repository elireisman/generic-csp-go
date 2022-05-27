package main

import (
	"fmt"
	"github.com/elireisman/generic-csp-go/pkg/csp"
)

type Province string
type Color string

var (
	// CSP variables
	Canada []Province

	// CSP domains
	Colors []Color

	// CSP constraints
	Constraints []Border[Province, Color]
)

// WTF! This should work fine, but results in:
// `cmd/map_coloring.go:112:25: cannot use border (variable of type Border[Province, Color]) as type csp.Constraint[Province, Color] in argument to problem.AddConstraint`
//
// References:
// https://github.com/golang/go/issues/44689
// https://stackoverflow.com/questions/66118867/go-generics-is-it-possible-to-embed-generic-structs
type Border[V, D comparable] struct {
	csp.Constraint[V, D]
}

func NewBorder(sideOne, sideTwo Province) Border[Province, Color] {
	return Border[Province, Color]{
		csp.Constraint[Province, Color]{
			Variables: []Province{sideOne, sideTwo},
		},
	}
}

// constraint: Ensure pair of province borders represented here
// are not assigned the same color in the candidate solution
func (b Border[Province, Color]) Satisfied(candidate map[Province]Color) bool {
	colorP1, foundP1 := candidate[b.Variables[0]]
	colorP2, foundP2 := candidate[b.Variables[1]]

	// if both provinces are not yet present in the candidate
	// solution, then (for now) the constraint is satisfied
	if !foundP1 || !foundP2 {
		return true
	}

	// if both provinces are present in the candidate
	// solution, their colors must not be the same
	return colorP1 != colorP2
}

func init() {
	Canada = []Province{
		"Yukon",
		"British Columbia",
		"Northwest Territories",
		"Nunavut",
		"Alberta",
		"Saskatchewan",
		"Manitoba",
		"Ontario",
		"Quebec",
		"Newfoundland and Laborador",
		"New Brunswick",
		"Nova Scotia",
		"Prince Edward Island",
	}

	Colors = []Color{
		"Red",
		"Yellow",
		"Blue",
		"Green",
	}

	Constraints = []Border[Province, Color]{
		NewBorder("Yukon", "British Columbia"),
		NewBorder("Yukon", "Northwest Territories"),
		NewBorder("British Columbia", "Alberta"),
		NewBorder("British Columbia", "Northwest Territories"),
		NewBorder("Northwest Territories", "Alberta"),
		NewBorder("Alberta", "Saskatchewan"),
		NewBorder("Saskatchewan", "Northwest Territories"),
		NewBorder("Nunavut", "Northwest Territories"),
		NewBorder("Saskatchewan", "Manitoba"),
		NewBorder("Manitoba", "Nunavut"),
		NewBorder("Manitoba", "Ontario"),
		NewBorder("Ontario", "Quebec"),
		NewBorder("Newfoundland and Laborador", "Quebec"),
		NewBorder("Newfoundland and Laborador", "Prince Edward Island"),
		NewBorder("Newfoundland and Laborador", "New Brunswick"),
		NewBorder("Newfoundland and Laborador", "Nova Scotia"),
		NewBorder("New Brunswick", "Quebec"),
		NewBorder("Nova Scotia", "New Brunswick"),
		NewBorder("Prince Edward Island", "New Brunswick"),
		NewBorder("Nova Scotia", "Prince Edward Island"),
	}
}

// model the map-coloring problem using CSP framework + Go generics
func main() {
	// assemble valid variable domains
	var domains map[Province][]Color
	for _, p := range Canada {
		domains[p] = Colors
	}
	// init empty solution to begin search through problem space
	candidate := map[Province]Color{}

	// create CSP framework instance, populate
	problem := csp.New(Canada, domains)
	for _, border := range Constraints {
		problem.AddConstraint(border)
	}

	// find ONE possible solution, and display it, if it exists
	if result := problem.Solve(candidate); result != nil {
		fmt.Println("Solution:")
		for p, c := range result {
			fmt.Printf("%s => %s\n", p, c)
		}
		return
	}

	panic("No solution found")
}
