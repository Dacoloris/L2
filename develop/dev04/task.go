package dev04

import (
	"sort"
	"strings"
)

func anagramSearch(source *[]string) map[string][]string {
	//Результирующая мапа
	result := make(map[string][]string)
	//Мапа, хранящая первые вхождения новых анаграмм
	firstValuesHolder := make(map[string]string)

	for _, word := range *source {
		//Если слово состоит из одного символа, то пропускаем его
		if len(word) < 2 {
			continue
		}

		processWord(word, &firstValuesHolder, &result)
	}

	//удаляем одиночные слова
	for k, v := range result {
		if len(v) == 0 {
			delete(result, k)
		}
	}

	return result
}

func processWord(word string, h *map[string]string, r *map[string][]string) {
	lowerCaseWord := strings.ToLower(word)

	sortedChars := []rune(lowerCaseWord)
	//Сортируем символы в слове, чтобы понять, есть ли ему анаграмма
	sort.Slice(sortedChars, func(i, j int) bool {
		return sortedChars[i] < sortedChars[j]
	})

	//Если уже имеется анаграмма к данному слову, то добавляем в результат
	//Иначе - устанавливаем слово в качестве ключа для новой группы анаграмм
	if v, ok := (*h)[string(sortedChars)]; ok {
		(*r)[v] = append((*r)[v], lowerCaseWord)
	} else {
		(*h)[string(sortedChars)] = lowerCaseWord
		(*r)[lowerCaseWord] = []string{}
	}
}
