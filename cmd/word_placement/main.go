package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/elireisman/generic-csp-go/pkg/csp"
)

const GridSize = 16

type Word string

type Point struct {
	Row int
	Col int
}

type Placement struct {
	Points []Point
}

type Letter struct {
	Rune  rune
	Color LetterColor
}

type LetterColor string

var (
	None   LetterColor = "\x1b[0;0m"
	Green  LetterColor = "\x1b[1;42m"
	Yellow LetterColor = "\x1b[1;41m"
)

var (
	// CSP variables
	Words []Word

	// CSP domains
	Orientations []Point
	Placements   map[Word][]Placement

	// CSP constraints
	Constraints []csp.Constraint[Word]
)

func init() {
	rand.Seed(time.Now().Unix())

	Words = []Word{
		"ANNA",
		"BRANDYN",
		"COURTNEY",
		"ELI",
		"FEDERICO",
		"HENRI",
		"LANE",
		"LORENA",
		"JUSTIN",
		"PATRICK",
		"SARAH",
	}

	// when building the Placements (CSP domains)
	// we need to add one candidate placement per
	// word _and_ per orientation on the puzzle
	// grid that the word could take on!
	Orientations = []Point{
		// VerticalDown
		Point{Row: 1, Col: 0},
		// VerticalUp
		Point{Row: -1, Col: 0},
		// HorizontalLeft
		Point{Row: 0, Col: -1},
		// HorizontalRight
		Point{Row: 0, Col: 1},
		// Diagonal Left Down
		Point{Row: 1, Col: -1},
		// Diagonal Left Up
		Point{Row: -1, Col: -1},
		// Diagonal Right Down
		Point{Row: 1, Col: 1},
		// Diagonal Right Up
		Point{Row: -1, Col: 1},
	}

	// create a domain of candidate placements and
	// a constraint per word we need to place on
	// the board
	Placements = map[Word][]Placement{}
	Constraints = []csp.Constraint[Word]{}
	for _, word := range Words {
		Placements[word] = generatePlacements(word)
		Constraints = append(Constraints, NewWord(word))
	}
}

func NewWord(word Word) csp.Constraint[Word] {
	return csp.Constraint[Word]{
		Variables: []Word{word},
	}
}

// check each existing placement in the candidate assingments for conflicts
// with the new (proposed) placement named in the constraint
func SatisfiesConstraint(wordConstraint csp.Constraint[Word], candidate map[Word]Placement) bool {
	nextWord := wordConstraint.Variables[0]
	nextPlacement := candidate[nextWord]

	for word, placement := range candidate {
		// this is the new placement we're testing the
		// other word placements against; skip it!
		if nextWord == word {
			continue
		}

		for i, point := range placement.Points {
			for j, nextPoint := range nextPlacement.Points {
				if point.Row == nextPoint.Row &&
					point.Col == nextPoint.Col &&
					nextWord[j] != word[i] {
					return false
				}
			}
		}
	}

	return true
}

func generatePlacements(word Word) []Placement {
	out := []Placement{}

	for row := 0; row < GridSize; row++ {
		for col := 0; col < GridSize; col++ {
			start := Point{Row: row, Col: col}
			for _, diff := range Orientations {
				if result := generatePlacement(start, word, diff); result != nil {
					out = append(out, Placement{Points: result})
				}
			}
		}
	}

	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

func generatePlacement(start Point, word Word, diff Point) []Point {
	endRow := start.Row + (diff.Row * (len(word) - 1))
	endCol := start.Col + (diff.Col * (len(word) - 1))

	if endRow < 0 || endRow >= GridSize {
		return nil
	}
	if endCol < 0 || endCol >= GridSize {
		return nil
	}

	out := []Point{}
	for ndx := 0; ndx < len(word); ndx++ {
		nextRow := start.Row + (diff.Row * ndx)
		nextCol := start.Col + (diff.Col * ndx)

		out = append(out,
			Point{
				Row: nextRow,
				Col: nextCol,
			})
	}

	return out
}

func renderGrid[V Word, D Placement](candidate map[Word]Placement) {
	// init puzzle board
	puzzle := [GridSize][]Letter{}
	for row := 0; row < GridSize; row++ {
		puzzle[row] = make([]Letter, GridSize, GridSize)
		for col := 0; col < GridSize; col++ {
			puzzle[row][col] = Letter{
				Rune:  'A' + rune(rand.Intn(26)),
				Color: Green,
			}
		}
	}

	// apply resolved placements to puzzle
	for word, placement := range candidate {
		for ndx := 0; ndx < len(word); ndx++ {
			row := placement.Points[ndx].Row
			col := placement.Points[ndx].Col
			letter := word[ndx]
			puzzle[row][col] = Letter{
				Rune:  rune(letter),
				Color: Yellow,
			}
		}
	}

	// render the puzzle with all placements
	for row := 0; row < GridSize; row++ {
		for col := 0; col < GridSize; col++ {
			letter := puzzle[row][col]
			fmt.Printf("%s%c%s ", letter.Color, letter.Rune, None)
		}
		fmt.Println()
	}
}

// model puzzle the word placement problem using CSP framework + Go generics
func main() {
	// create CSP framework instance, populate
	problem := csp.New(Placements, SatisfiesConstraint)
	for _, wordToPlace := range Constraints {
		problem.AddConstraint(wordToPlace)
	}

	// init empty solution to begin search through problem space
	candidate := map[Word]Placement{}

	// find ONE possible solution, and display it, if it exists
	if result := problem.Solve(candidate); result != nil {
		fmt.Println("Solution:")
		renderGrid(result)
		return
	}

	panic("No solution found")
}
