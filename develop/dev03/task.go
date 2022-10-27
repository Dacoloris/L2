package dev03

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type IndexWord struct {
	index int
	word  string
}

func sortWords(src []IndexWord, asc bool) []int {
	indexes := make(map[string]int)
	var toSort []string
	for _, v := range src {
		indexes[v.word] = v.index
		toSort = append(toSort, v.word)
	}
	sort.Strings(toSort)
	var resultIndexes []int
	switch asc {
	case true:
		for i := 0; i < len(toSort); i++ {
			resultIndexes = append(resultIndexes, indexes[toSort[i]])
		}
	case false:
		for i := len(toSort) - 1; i >= 0; i-- {
			resultIndexes = append(resultIndexes, indexes[toSort[i]])
		}
	}

	return resultIndexes
}

func Task() {
	var column int
	var n, r, u bool
	flag.IntVar(&column, "k", 0, "указание колонки для сортировки")
	flag.BoolVar(&n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&r, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&u, "u", false, "не выводить повторяющиеся строки")
	flag.Parse()

	var in io.Reader
	if fileName := flag.Arg(0); fileName == "" {
		fmt.Printf("Не указано название файла\n")
		os.Exit(1)
	} else {
		f, err := os.Open(fileName)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Printf("Error closing file: %v", err)
			}
		}(f)
		if err != nil {
			fmt.Printf("Error opening file: %v", err)
			os.Exit(1)
		}
		in = f
	}

	var lines [][]string
	buf := bufio.NewScanner(in)
	for buf.Scan() {
		line := buf.Text()
		lines = append(lines, strings.Split(line, " "))
	}
	var wordsWithIndexes []IndexWord
	for i, line := range lines {
		wordsWithIndexes = append(wordsWithIndexes, IndexWord{
			index: i,
			word:  line[column],
		})
	}
	result := sortWords(wordsWithIndexes, !r)
	if u {
		var prevLineIndex int
		for i, index := range result {
			if i > 0 &&
				strings.Join(lines[prevLineIndex], " ") == strings.Join(lines[index], " ") {
				continue
			}
			fmt.Println(lines[index])
			prevLineIndex = index
		}
	} else {
		for _, index := range result {
			fmt.Println(lines[index])
		}
	}
}
