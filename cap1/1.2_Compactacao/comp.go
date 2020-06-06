package main

import (
	"fmt"
	"math/big"
	"os"
)

type gene struct {
	compressed *big.Int
}

func (g *gene) compressGene(str string) {
	g.compressed = big.NewInt(1)
	for _, nucleotide := range str {
		g.compressed.Lsh(g.compressed, 2)
		if nucleotide == 'A' {
			g.compressed.Or(g.compressed, big.NewInt(0))
		} else if nucleotide == 'C' {
			g.compressed.Or(g.compressed, big.NewInt(1))
		} else if nucleotide == 'T' {
			g.compressed.Or(g.compressed, big.NewInt(2))
		} else if nucleotide == 'G' {
			g.compressed.Or(g.compressed, big.NewInt(3))
		} else {
			os.Exit(1)
		}
	}
}

func (g *gene) deCompressGene() string {
	str := ""
	lCompressed := g.compressed.BitLen() - 1
	for i := lCompressed - 2; i >= 0; i -= 2 {
		res := g.compressed.Bit(i) + g.compressed.Bit(i+1)*2
		if res == 0 {
			str = str + "A"
		} else if res == 1 {
			str = str + "C"
		} else if res == 2 {
			str = str + "T"
		} else if res == 3 {
			str = str + "G"
		}
	}
	return str
}

func main() {
	gene := gene{}
	gene.compressGene("ACTGAACCTTGGACTGAACCTTGGACTGAACCTTGGACTGAACCTTGGACTGAACCTTGGACTGAACCTTGGACTGAACCTTGGACTGAACCTTGG")
	fmt.Println(gene.compressed)
	fmt.Println(gene.deCompressGene())
}
