package main

import (
	"fmt"
	"sort"
)

type codon [3]byte
type gene []codon

func strCodon(c codon) string {
	return string(c[0]) + string(c[1]) + string(c[2])
}

func (g gene) Len() int { return len(g) }

func (g gene) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

func (g gene) Less(i, j int) bool {
	return strCodon(g[i]) < strCodon(g[j])
}

func (g gene) String() string {
	str := ""
	for _, val := range g {
		str += strCodon(val) + " "
	}
	return fmt.Sprintf("%s", str)
}

func (g *gene) loadGene(s string) {
	for i := 0; i < len(s); i += 3 {
		if (i + 2) >= len(s) {
			return
		}
		*g = append(*g, codon{s[i], s[i+1], s[i+2]})
	}
	return
}

func (g gene) linearContains(c codon) bool {
	for _, v := range g {
		if v == c {
			return true
		}
	}
	return false
}

func (g gene) binaryContains(c codon) bool {
	low := 0
	high := len(g) - 1
	for low <= high {
		mid := (low + high) / 2
		if strCodon(g[mid]) < strCodon(c) {
			low = mid + 1
		} else if strCodon(g[mid]) > strCodon(c) {
			high = mid - 1
		} else {
			return true
		}
	}
	return false
}

func main() {
	g := gene{}
	g.loadGene("ACTGACTGACTGACTGCGATCGATAAATTGGCGAGTCGAGCTAGCTAGCGGATGCGGATGAGCGCGCGCG")
	fmt.Println(g)
	gac := codon{'G', 'A', 'C'}
	gtc := codon{'G', 'T', 'C'}
	fmt.Println(g.linearContains(gac))
	fmt.Println(g.linearContains(gtc))
	sort.Sort(g)
	fmt.Println(g)
	fmt.Println(g.binaryContains(gac))
	fmt.Println(g.binaryContains(gtc))

}
