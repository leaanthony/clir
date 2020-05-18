package main

import "fmt"

func findMissing(s, r []string) string {
	rmap := make(map[string]bool, len(r))
	for _, kr := range r {
		rmap[kr] = true
	}
	for _, ks := range s {
		if !rmap[ks] {
			return ks
		}
	}
	return ""
}

func main() {
	// These are the strings that are required to be supplied
	required := []string{"a", "b", "c"}
	// These are the strings that were supplied (missing c)
	supplied := []string{"d", "a", "b", "f", "g", "c"}

	result := findMissing(required, supplied)
	fmt.Println("Result: ", result)

}
