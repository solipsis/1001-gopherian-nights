package main

import (
	"fmt"
	"sort"
)

// simple struct
type puppy struct {
	age  int
	name string
}

// two ways to sort
type byAge []puppy
type byName []puppy

// implement the sort.Interface interface for both sorting types
func (b byAge) Less(i, j int) bool { return b[i].age < b[j].age }
func (b byAge) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byAge) Len() int           { return len(b) }

func (b byName) Less(i, j int) bool { return b[i].name < b[j].name }
func (b byName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byName) Len() int           { return len(b) }

func main() {
	puppies := []puppy{
		{1, "barkley"},
		{4, "princess"},
		{99, "golden"},
		{10, "grey"},
	}

	fmt.Printf("Original: %v\n", puppies)
	sort.Sort(byAge(puppies))
	fmt.Printf("ByAge: %v\n", puppies)
	sort.Sort(byName(puppies))
	fmt.Printf("ByName: %v\n", puppies)
	sort.Sort(sort.Reverse(byAge(puppies)))
	fmt.Printf("Reverse ByAge: %v\n", puppies)
}
