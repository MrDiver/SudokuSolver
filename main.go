package main

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var ctx, cancel = context.WithTimeout(context.TODO(), time.Minute)
var shouldStop = false
var solutions = make(chan [9][9]int, 1000)

func printField(feld [9][9]int) {
	fmt.Println("┌───────┬───────┬───────┐")
	for i, zeile := range feld {
		chars := []string{}
		for _, v := range zeile {
			if v == 0 {
				chars = append(chars, ".")
			} else {
				chars = append(chars, strconv.Itoa(v))
			}
		}
		fmt.Printf("│ %v %v %v │ %v %v %v │ %v %v %v │\n", chars[0], chars[1], chars[2], chars[3], chars[4], chars[5], chars[6], chars[7], chars[8])
		if i == 2 || i == 5 {
			fmt.Println("├───────┼───────┼───────┤")
		}

	}
	fmt.Println("└───────┴───────┴───────┘")
}

type Position struct {
	x int
	y int
}

func findNextFree(feld [9][9]int) (Position, bool) {
	for y, zeile := range feld {
		for x, v := range zeile {
			if v == 0 {
				return Position{x: x, y: y}, false
			}
		}
	}

	return Position{}, true
}

func filter(l []int, x int) []int {
	var entries = []int{}
	for _, e := range l {
		if e != x {
			entries = append(entries, e)
		}
	}
	return entries
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

// Returns true if remove
func checkNeighbor(nval int, val int) bool {
	if nval != 0 && Abs(nval-val) < 2 {
		return true
	}
	return false
}

func getValidEntries(feld [9][9]int, p Position) []int {
	var contains = []int{}

	// Prüfe zeilen und spalten
	for y, zeile := range feld {
		for x, v := range zeile {
			if x == p.x && v != 0 {
				contains = append(contains, v)
			}
			if y == p.y && v != 0 {
				contains = append(contains, v)
			}
		}
	}

	// Prüfe block
	a := (p.x / 3) * 3
	b := (p.y / 3) * 3
	for y := range 3 {
		for x := range 3 {
			ix := a + x
			iy := b + y
			if feld[iy][ix] != 0 {
				contains = append(contains, feld[iy][ix])
			}
		}
	}

	// generate uncontained numbers
	var entries = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, x := range contains {
		entries = filter(entries, int(x))
	}

	var toRemove = []int{}

	for _, v := range entries {
		// Links
		if p.x-1 >= 0 {
			if checkNeighbor(feld[p.y][p.x-1], v) {
				toRemove = append(toRemove, v)
			}
		}

		// Rechts
		if p.x+1 <= 8 {
			if checkNeighbor(feld[p.y][p.x+1], v) {
				toRemove = append(toRemove, v)
			}
		}

		// Oben
		if p.y-1 >= 0 {
			if checkNeighbor(feld[p.y-1][p.x], v) {
				toRemove = append(toRemove, v)
			}
		}

		// Unten
		if p.y+1 <= 8 {
			if checkNeighbor(feld[p.y+1][p.x], v) {
				toRemove = append(toRemove, v)
			}
		}
	}
	for _, x := range toRemove {
		entries = filter(entries, x)
	}

	return entries
}

func tryNextStep(feld [9][9]int) {
	p, done := findNextFree(feld)
	if done {
		// if feld[0][5] == 7 && feld[1][4] == 1 && feld[2][3] == 2 {
		if validSolution(feld) {
			solutions <- feld
		}
		// }
		return
	}

	tryEntries := getValidEntries(feld, p)
	for _, t := range tryEntries {
		feld[p.y][p.x] = t
		// printField(feld)
		if shouldStop {
			return
		}
		go tryNextStep(feld)
	}
}

func validSolution(feld [9][9]int) bool {
	for y := range 9 {
		for x := range 9 {
			// Links
			if x-1 >= 0 {
				if checkNeighbor(feld[y][x-1], feld[y][x]) {
					return false
				}
			}

			// Rechts
			if x+1 <= 8 {
				if checkNeighbor(feld[y][x+1], feld[y][x]) {
					return false
				}
			}

			// Oben
			if y-1 >= 0 {
				if checkNeighbor(feld[y-1][x], feld[y][x]) {
					return false
				}
			}

			// Unten
			if y+1 <= 8 {
				if checkNeighbor(feld[y+1][x], feld[y][x]) {
					return false
				}
			}
		}
	}
	return true
}

func main() {
	feld := [9][9]int{
		{0, 0, 0, 0, 4, 0, 0, 0, 0},
		{0, 9, 0, 0, 0, 0, 0, 6, 0},
		{0, 0, 4, 0, 9, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 3, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 3, 0, 6, 0, 2, 0, 0},
		{0, 2, 0, 0, 0, 0, 0, 4, 0},
		{0, 0, 0, 0, 7, 0, 0, 0, 0}}

	// feld := [9][9]int{
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 0, 0}}

	printField(feld)

	go tryNextStep(feld)

	go func() {
		for {
			if runtime.NumGoroutine() < 3 {
				cancel()
				return
			}
		}
	}()

END:
	for {
		select {
		case s := <-solutions:
			printField(s)
		case <-ctx.Done():
			break END
		}
	}
	cancel()
	shouldStop = true
}
