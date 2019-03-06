package main

import (
	"fmt"
)

type suitor struct {
	id          int
	preferences []int
}

type suited struct {
	id             int
	choice         int
	preferences    []int
	currentSuitors map[int]bool
}

//func propose

func main() {

	suitors := []suitor{
		{0, []int{3, 5, 4, 2, 1, 0}},
		{1, []int{2, 3, 1, 0, 4, 5}},
		{2, []int{5, 2, 1, 0, 3, 4}},
		{3, []int{0, 1, 2, 3, 4, 5}},
		{4, []int{4, 5, 1, 2, 0, 3}},
		{3, []int{0, 1, 2, 3, 4, 5}},
	}

	suiteds := []suited{
		{0, 99, []int{3, 5, 4, 2, 1, 0}, map[int]bool{}},
		{1, 99, []int{2, 3, 1, 0, 4, 5}, map[int]bool{}},
		{2, 99, []int{5, 2, 1, 0, 3, 4}, map[int]bool{}},
		{3, 99, []int{0, 1, 2, 3, 4, 5}, map[int]bool{}},
		{4, 99, []int{4, 5, 1, 2, 0, 3}, map[int]bool{}},
		{5, 99, []int{0, 1, 2, 3, 4, 5}, map[int]bool{}},
	}

	unpaired := make([]int, len(suitors))
	for n := range suitors {
		unpaired[n] = n
	}

	for len(unpaired) > 0 {
		fmt.Println("UP:", unpaired)

		//log.Fatal("Use delve here")

		fmt.Println("Proposing")
		for _, suitor := range unpaired {
			pref := suitors[suitor].preferences[0]
			suiteds[pref].currentSuitors[suitor] = true
		}
		unpaired = []int{}

		fmt.Println("Rejecting")
		// for each suited
		// if suitors choose the best and reject the rest
		for i, suited := range suiteds {

			proposed := suited.currentSuitors
			fmt.Println("Current suitors for: ", suited)
			//bestIdx := 99
			for _, suitor := range suited.preferences {
				if proposed[suitor] {
					suiteds[i].choice = suitor
					break
				}
			}

			//fmt.Println("idx", bestIdx)
			//suiteds[i].choice = suited.preferences[bestIdx]

			// reject everyone not the favourite
			rejected := []int{}
			for k := range proposed {
				if k != suiteds[i].choice {
					rejected = append(rejected, k)
				}
			}
			unpaired = append(unpaired, rejected...)

			suiteds[i].currentSuitors = map[int]bool{}
			if suiteds[i].choice != 99 {
				suiteds[i].currentSuitors[suiteds[i].choice] = true
			}
		}

		// everyone who was rejected needs to update their preferences
		for _, u := range unpaired {
			fmt.Println(unpaired)
			suitors[u].preferences = suitors[u].preferences[1:]
		}
	}

	for _, s := range suiteds {
		fmt.Println(s)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
