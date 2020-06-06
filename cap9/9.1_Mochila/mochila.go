package main

import "fmt"

type item struct {
	name   string
	weight int
	value  float64
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func knapsack(items []item, maxCapacity int) []item {
	table := [][]float64{}
	for i := 0; i <= len(items); i++ {
		row := make([]float64, maxCapacity+1)
		table = append(table, row)
	}

	for i, it := range items {
		for capacity := 1; capacity <= maxCapacity; capacity++ {
			previousItemValue := table[i][capacity]
			if capacity >= it.weight {
				valueFreeingWeightForItem := table[i][capacity-it.weight]
				table[i+1][capacity] = max(valueFreeingWeightForItem+it.value, previousItemValue)
			} else {
				table[i+1][capacity] = previousItemValue
			}
		}
	}
	solution := []item{}
	capacity := maxCapacity

	for i := len(items); i > 0; i-- {
		if table[i-1][capacity] != table[i][capacity] {
			solution = append(solution, items[i-1])
			capacity -= items[i-1].weight
		}
	}
	return solution
}

func main() {
	items := []item{{"television", 50, 500},
		{"candlestick", 2, 300},
		{"stereo", 35, 400},
		{"laptop", 3, 1000},
		{"food", 15, 50},
		{"clothing", 20, 800},
		{"jewelry", 1, 4000},
		{"books", 100, 300},
		{"printer", 18, 30},
		{"refrigerator", 200, 700},
		{"painting", 10, 1000}}
	fmt.Println(knapsack(items, 75))
}
