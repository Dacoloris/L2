package dev02

import (
	"strings"
	"unicode"
)

func primitiveExtract(str string) string {
	if len(str) <= 1 {
		return str
	}

	var result []rune
	symbols := []rune(str)
	for len(symbols) > 0 {
		if len(symbols) == 1 {
			if unicode.IsDigit(symbols[0]) {
				return "некорректная строка"
			}
			result = append(result, symbols[0])
			break
		}

		fSymbol, sSymbol := symbols[0], symbols[1]
		if unicode.IsDigit(fSymbol) {
			return "некорректная строка"
		}
		if unicode.IsDigit(sSymbol) {
			result = append(result, []rune(strings.Repeat(string(fSymbol), int(sSymbol-'0')))...)
			symbols = symbols[2:]
		} else {
			result = append(result, fSymbol)
			symbols = symbols[1:]
		}
	}

	return string(result)
}
