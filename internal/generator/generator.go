package generator

import (
	"math/rand"
)

const urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"
const maxLen = len(urlAlphabet)

func GetRandomStr(size int) string {
	var result string

	for i := 0; i < size; i++ {
		result += string(urlAlphabet[rand.Intn(maxLen)])
	}

	return result
}
