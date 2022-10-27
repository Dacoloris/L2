package dev06

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

/*
=== Утилита cut ===
Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные
Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var fields string
var delim string
var separated bool

func init() {
	testing.Init()
	flag.StringVar(&fields, "f", "", "")
	flag.StringVar(&delim, "d", "\t", "")
	flag.BoolVar(&separated, "s", false, "")
	flag.Parse()
}

// регулярные выражения для parseFields
var re0 = regexp.MustCompile(`^\d+$`)
var re1 = regexp.MustCompile(`^\d+\-\d+$`)
var re2 = regexp.MustCompile(`^\-\d+$`)
var re3 = regexp.MustCompile(`^\d+\-$`)

var errFieldsOne = errors.New("fields are numbered from 1")
var errFieldsDecr = errors.New("invalid decreasing range")

// Функция parseFields парсит строку, переданную через флаг -f
// и возвращает слайс интервалов, отображаемых полей
func parseFields(s string) (fs [][]int, err error) {
	rsStrings := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })

	// не указано ни одного поля: ","
	if len(rsStrings) == 0 {
		return nil, errFieldsOne
	}

	// указано одно поле (диапазон) перед или после запятой: ",2", "2-4,"
	contaisComma, err := regexp.MatchString(",", s)
	if err != nil {
		return nil, err
	}
	if len(rsStrings) == 1 && contaisComma {
		return nil, errFieldsOne
	}
	for _, rString := range rsStrings {

		// если подстрока соответствует регулярному выражению, диапазон заносится в результат
		switch {
		case re0.MatchString(rString):
			a, _ := strconv.Atoi(rString)
			if a < 1 {
				return nil, errFieldsOne
			}
			fs = append(fs, []int{a, a})
		case re1.MatchString(rString):
			ds := strings.FieldsFunc(rString, func(r rune) bool { return r == '-' })
			a, _ := strconv.Atoi(ds[0])
			if a < 1 {
				return nil, errFieldsOne
			}
			b, _ := strconv.Atoi(ds[1])
			if a > b {
				return nil, errFieldsDecr
			}
			fs = append(fs, []int{a, b})
		case re2.MatchString(rString):
			b, _ := strconv.Atoi(rString[len("-"):])
			if b < 1 {
				return nil, errFieldsDecr
			}
			fs = append(fs, []int{1, b})
		case re3.MatchString(rString):
			a, _ := strconv.Atoi(rString[:len(rString)-len("-")])
			if a < 1 {
				return nil, errFieldsOne
			}
			fs = append(fs, []int{a, math.MaxInt})
		default:
			return nil, errors.New("invalid fields value")
		}
	}
	return fs, nil
}

// Функция cutAndPrint печатает поля из диапазонов fs строки s.
func cutAndPrint(s string, fs [][]int, dr rune) {

	// поиск подстрок, разделенных символом dr
	var sf []string
	rs := []rune(s)
	st := 0
	match := false
	for i, r := range rs {
		if r == dr {
			match = true
			if i == st {
				sf = append(sf, "")
			} else {
				sf = append(sf, string(rs[st:i]))
			}
			st = i + 1
		} else if i == len(rs)-1 {
			sf = append(sf, string(rs[st:]))
		}
	}

	// строка не содержит разделителя
	if !match {
		if separated {
			return
		}
		fmt.Println(s)
		return
	}

	// случай, когда строка оканчивается символом-разделителем
	if lr, _ := utf8.DecodeLastRuneInString(s); lr == dr {
		sf = append(sf, "")
	}

	// выбор подстрок (полей), входящих в диапазоны fs
	var outSlice []string
	for i, f := range sf {
		if indexInSegments(i+1, fs) {
			outSlice = append(outSlice, f)
		}
	}
	fmt.Println(strings.Join(outSlice, string(dr)))
}

// Функция indexInSegments возвращает true, если i принадлежит одному из отрезков в segs
func indexInSegments(i int, segs [][]int) bool {
	res := segs[0][0] <= i && i <= segs[0][1]
	for j := 1; j < len(segs); j++ {
		res = res || (segs[j][0] <= i && i <= segs[j][1])
	}
	return res
}

func main() {
	if fields == "" {
		fmt.Fprintln(os.Stderr, "fields are not specified")
		return
	}

	fs, err := parseFields(fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// проверка того, что в -d передан один символ-руна
	if dn := utf8.RuneCountInString(delim); dn != 1 {
		fmt.Fprintln(os.Stderr, "the delimiter must be a single character")
		return
	}
	dr, _ := utf8.DecodeRuneInString(delim)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		cutAndPrint(sc.Text(), fs, dr)
	}
}
