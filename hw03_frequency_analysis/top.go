package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// word структура для хранения слов и их частоты использования.
type word struct {
	word      string
	frequency int
}

// CountStringForResult количество слов в результате выполнения функции.
const CountStringForResult = 10

// excludedWords слова для исключения.
var excludedWords = map[string]struct{}{"-": {}}

// regular регулярное выражение для получения слов.
var regular = regexp.MustCompile(`[a-zA-Z0-9а-яА-Я\-]+`)

// Top10 функция для получения 10 часто повторяющихся слов.
func Top10(str string) []string {
	sl := wordsWithFrequency(str)

	sort.Slice(sl, func(a, b int) bool {
		switch {
		case sl[a].frequency > sl[b].frequency:
			return true
		case sl[a].frequency == sl[b].frequency:
			return sl[a].word < sl[b].word
		default:
			return false
		}
	})

	resultSlice := make([]string, 0, CountStringForResult)
	for i, s := range sl {
		resultSlice = append(resultSlice, s.word)
		if (i + 1) == CountStringForResult {
			break
		}
	}

	return resultSlice
}

// wordsWithFrequency получить слайс слов с их частотой использования.
func wordsWithFrequency(str string) []word {
	sl := regular.FindAllString(str, -1)
	sl = toLowerStringsAndExcludeSpecialWords(sl)
	frequencyMap := make(map[string]int, len(sl))

	for _, v := range sl {
		frequencyMap[v]++
	}

	resultSlice := make([]word, 0, len(frequencyMap))
	for w, frequency := range frequencyMap {
		resultSlice = append(resultSlice, word{w, frequency})
	}

	return resultSlice
}

// toLowerStringsAndExcludeSpecialWords функция для приведения слов к единому регистру и исключения специальных слов.
func toLowerStringsAndExcludeSpecialWords(sl []string) []string {
	result := make([]string, 0, len(sl))
	for _, v := range sl {
		if _, ok := excludedWords[v]; ok {
			continue
		}
		result = append(result, strings.ToLower(v))
	}

	return result
}
