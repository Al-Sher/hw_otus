package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ErrInvalidString ошибка для некорректной строки.
var ErrInvalidString = errors.New("invalid string")

// Unpack функция распаковки строки.
func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	prevRune, firstRuneLen := utf8.DecodeRuneInString(str)
	if unicode.IsDigit(prevRune) || (prevRune == utf8.RuneError && firstRuneLen == 1) {
		return "", ErrInvalidString
	}

	var result strings.Builder
	escaped := false
	str = str[firstRuneLen:]

	for _, currentRune := range str {
		if prevRune == '\\' && !escaped {
			escaped = true
			prevRune = currentRune
			continue
		}

		if escaped && prevRune != '\\' && !unicode.IsDigit(prevRune) {
			return "", ErrInvalidString
		}

		if !escaped && unicode.IsDigit(prevRune) && unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}

		repeat, err := repeatIfNeed(&result, prevRune, currentRune)
		if err != nil {
			return "", err
		}
		if repeat {
			escaped = false
			prevRune = currentRune
			continue
		}

		writeSymbol(&result, escaped, prevRune)
		escaped = false
		prevRune = currentRune
	}

	lastRune, _ := utf8.DecodeLastRuneInString(str)

	if lastRune != '\\' && !unicode.IsDigit(lastRune) && escaped {
		return "", ErrInvalidString
	}

	writeSymbol(&result, escaped, lastRune)

	return result.String(), nil
}

// repeatIfNeed функция для повтора символа, если текущий символ является числом.
func repeatIfNeed(buf *strings.Builder, prevSymbol, currentSymbol rune) (bool, error) {
	if !unicode.IsDigit(currentSymbol) {
		return false, nil
	}

	n, err := strconv.Atoi(string(currentSymbol))
	if err != nil {
		return false, err
	}
	buf.WriteString(strings.Repeat(string(prevSymbol), n))
	return true, nil
}

// writeSymbolIsNotDigit функция записи в буфер, в случае если символ не является числом.
func writeSymbol(buf *strings.Builder, escaped bool, symbol rune) {
	if !unicode.IsDigit(symbol) || escaped {
		buf.WriteString(string(symbol))
	}
}
