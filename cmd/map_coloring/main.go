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
	Constraints []csp.Constraint[Province]
)

func NewBorder(us, them Province) csp.Constraint[Province] {
	return csp.Constraint[Province]{
		Variables: []Province{us, them},
	}
}

// constraint: Ensure pair of province borders represented here
// are not assigned the same color in the candidate solution
func Satisfied[V, D comparable](border csp.Constraint[V], candidate map[V]D) bool {
	colorP1, foundP1 := candidate[border.Variables[0]]
	colorP2, foundP2 := candidate[border.Variables[1]]

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

	Constraints = []csp.Constraint[Province]{
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
	// assemble mapping of variables to a set of possible
	// values to search for a valid solution
	domain := map[Province][]Color{}
	for _, p := range Canada {
		domain[p] = Colors
	}

	// create CSP framework instance, populate
	problem := csp.New(domain, Satisfied[Province, Color])
	for _, border := range Constraints {
		problem.AddConstraint(border)
	}

	// init empty solution to begin search through problem space
	candidate := map[Province]Color{}

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
