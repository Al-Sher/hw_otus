package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString ошибка для некорректной строки.
var ErrInvalidString = errors.New("invalid string")

// Unpack функция распаковки строки.
func Unpack(str string) (string, error) {
	r := []rune(str)
	var result strings.Builder
	escaped := false

	if len(r) == 0 {
		return "", nil
	}

	if unicode.IsDigit(r[0]) {
		return "", ErrInvalidString
	}

	for i := 1; i < len(r); i++ {
		if r[i-1] == '\\' && !escaped {
			escaped = true
			continue
		}

		if err := escapedError(r[i-1]); escaped && err != nil {
			return "", err
		}

		if err := numberError(r[i-1], r[i]); !escaped && err != nil {
			return "", err
		}

		repeat, err := repeatIfNeed(&result, r[i-1], r[i])
		if err != nil {
			return "", err
		}
		if repeat {
			escaped = false
			continue
		}

		writeSymbol(&result, escaped, r[i-1])
		escaped = false
	}

	if err := escapedError(r[len(r)-1]); escaped && err != nil {
		return "", err
	}

	writeSymbol(&result, escaped, r[len(r)-1])

	return result.String(), nil
}

// numberError функция для проверки на ошибку, в случае если два символа подряд
// являются цифрами.
func numberError(prevSymbol, currentSymbol rune) error {
	if unicode.IsDigit(prevSymbol) && unicode.IsDigit(currentSymbol) {
		return ErrInvalidString
	}

	return nil
}

// repeatIfNeed функция для повтора символа, если текущий символ является числом.
func repeatIfNeed(buf *strings.Builder, prevSymbol, currentSymbol rune) (bool, error) {
	if unicode.IsDigit(currentSymbol) {
		n, err := strconv.Atoi(string(currentSymbol))
		if err != nil {
			return false, err
		}
		buf.WriteString(strings.Repeat(string(prevSymbol), n))

		return true, nil
	}

	return false, nil
}

// writeSymbolIsNotDigit функция записи в буфер, в случае если символ не является числом.
func writeSymbol(buf *strings.Builder, escaped bool, symbol rune) {
	if !unicode.IsDigit(symbol) || escaped {
		buf.WriteString(string(symbol))
	}
}

// escapedError функция проверки, что следующий символ после слэша является валидным.
func escapedError(symbol rune) error {
	if symbol != '\\' && !unicode.IsDigit(symbol) {
		return ErrInvalidString
	}

	return nil
}
