package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

func generateRandomKey(l int) *big.Int {
	randomKey := big.NewInt(0)
	token := make([]byte, l)
	rand.Read(token)
	for i := 0; i < l; i++ {
		randomKey.Lsh(randomKey, 8)
		randomKey.Or(randomKey, big.NewInt(int64(token[i])))
	}
	return randomKey
}

func encodeMessage(s string) (*big.Int, *big.Int) {
	randomKey := generateRandomKey(len(s))
	encMessage := big.NewInt(0)
	for i := 0; i < len(s); i++ {
		encMessage.Lsh(encMessage, 8)
		encMessage.Or(encMessage, big.NewInt(int64(s[i])))
	}
	encMessage.Xor(encMessage, randomKey)
	return encMessage, randomKey
}

func decodeMessage(encMessage *big.Int, randomKey *big.Int) string {
	strMessage := ""
	message := big.NewInt(0)
	message.Xor(encMessage, randomKey)
	for _, val := range message.Bytes() {
		strMessage += string(val)
	}
	return strMessage
}

func main() {
	encMessage, randomKey := encodeMessage("ARLima")
	fmt.Println(randomKey)
	fmt.Println(encMessage)
	fmt.Println(decodeMessage(encMessage, randomKey))
}
