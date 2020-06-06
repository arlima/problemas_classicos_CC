package main

import "fmt"

func possibleMnemonics(phoneNumber string, phoneMapping map[string][]string) [][]string {
	letterTuples := [][]string{}
	for _, digit := range phoneNumber {
		letterTuples = append(letterTuples, phoneMapping[string(digit)])
	}
	return letterTuples
}

func product2(setA [][]string, setB []string) [][]string {
	res := [][]string{}
	for _, a := range setA {
		temp := append([]string{}, a...)
		for _, b := range setB {
			res = append(res, append(temp, b))
		}
	}
	return res
}

func product(set [][]string) [][]string {
	res := [][]string{}
	for _, a := range set[0] {
		for _, b := range set[1] {
			res = append(res, []string{a, b})
		}
	}

	for i := 2; i < len(set); i++ {
		res = product2(res, set[i])
	}
	return res
}

func main() {

	var phoneNumber string

	phoneMapping := make(map[string][]string)
	phoneMapping["1"] = []string{"1"}
	phoneMapping["2"] = []string{"a", "b", "c"}
	phoneMapping["3"] = []string{"d", "e", "f"}
	phoneMapping["4"] = []string{"g", "h", "i"}
	phoneMapping["5"] = []string{"j", "k", "l"}
	phoneMapping["6"] = []string{"m", "n", "o"}
	phoneMapping["7"] = []string{"p", "q", "r", "s"}
	phoneMapping["8"] = []string{"t", "u", "v"}
	phoneMapping["9"] = []string{"w", "x", "y", "z"}
	phoneMapping["0"] = []string{"0"}

	fmt.Printf("Enter a phone number: ")
	fmt.Scan(&phoneNumber)

	pm := possibleMnemonics(phoneNumber, phoneMapping)
	res := product(pm)

	fmt.Printf("Here are the potential mnemonics: \n")
	for _, i := range res {
		fmt.Println(i)
	}
}
