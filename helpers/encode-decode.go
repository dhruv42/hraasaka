package helpers

import (
	"math"
	"strings"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var COUNTER = 1000

func Base62Encode(n int) string {
	var hashStr string
	for n > 0 {
		hashStr += string(chars[n%62])
		n = n / 62
	}
	return hashStr
}

func Base62Decode(hashStr string) int {
	var number = 0
	var power = 0
	strLength := len(hashStr)

	for strLength > 0 {
		number += strings.Index(chars, string(hashStr[strLength])) * powInt(62, power)
		power++
		strLength--
	}
	return number
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func GetCounter() int64 {
	COUNTER++
	return int64(COUNTER)
}
