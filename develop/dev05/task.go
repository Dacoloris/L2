package dev05

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"testing"
)

var after uint64
var before uint64
var context uint64
var count bool
var ignore bool
var invert bool
var fixed bool
var numerate bool

func init() {
	testing.Init()
	flag.Uint64Var(&after, "A", 0, "")
	flag.Uint64Var(&before, "B", 0, "")
	flag.Uint64Var(&context, "C", 0, "")
	flag.BoolVar(&count, "c", false, "")
	flag.BoolVar(&ignore, "i", false, "")
	flag.BoolVar(&invert, "v", false, "")
	flag.BoolVar(&fixed, "F", false, "")
	flag.BoolVar(&numerate, "n", false, "")
	flag.Parse()
}

func parseArgs(args []string) (re *regexp.Regexp, file string, err error) {
	switch {
	case len(args) == 0:
		return nil, "", errors.New("no pattern specified")
	case len(args) == 1:
		ps := args[0]
		if fixed {
			ps = `\Q` + ps + `\E`
		}
		if ignore {
			ps = `(?i)(` + ps + `)`
		}
		re, err := regexp.Compile(ps)
		if err != nil {
			return nil, "", err
		}
		return re, "", nil
	default:
		st, en := `(`, `)`
		if fixed {
			st, en = `(\Q`, `\E)`
		}
		ps := st + args[0] + en
		for i := 1; i < len(args)-1; i++ {
			ps += `|` + st + args[i] + en
		}
		if ignore {
			ps = `(?i)(` + ps + `)`
		}
		re, err := regexp.Compile(ps)
		if err != nil {
			return nil, "", err
		}
		return re, args[len(args)-1], nil
	}
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type scanner struct {
	sc *bufio.Scanner
	f  *os.File
}

func newScanner(file string) (*scanner, error) {
	if file == "" {
		sc := bufio.NewScanner(os.Stdin)
		return &scanner{sc: sc}, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)
	return &scanner{sc, f}, nil
}

func readLines(scan *scanner) (lines [][]byte) {
	for scan.sc.Scan() {
		lines = append(lines, scan.sc.Bytes())
	}
	return lines
}

func filter(lines [][]byte, re *regexp.Regexp) (s []int) {
	for i, line := range lines {
		if re.Match(line) != invert {
			s = append(s, i)
		}
	}
	return s
}

func printLine(num int, line []byte, match bool) {
	switch {
	case numerate && match:
		fmt.Printf("%d:%s\n", num, line)
	case numerate && !match:
		fmt.Printf("%d-%s\n", num, line)
	default:
		fmt.Println(string(line))
	}
}

func printLines(lines [][]byte, n []int) {
	h := 0
	for i := 0; i < len(n)-1; i++ {
		b := maxInt(h, n[i]-int(before))
		a := minInt(n[i]+int(after), n[i+1]-1)
		for j := b; j <= a; j++ {
			printLine(j+1, lines[j], n[i] == j)
		}
		h = a + 1
	}
	b := maxInt(h, n[len(n)-1]-int(before))
	a := minInt(n[len(n)-1]+int(after), len(lines)-1)
	for j := b; j <= a; j++ {
		printLine(j+1, lines[j], n[len(n)-1] == j)
	}
}

func main() {
	after = maxUint64(after, context)
	before = maxUint64(before, context)

	re, file, err := parseArgs(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	scan, err := newScanner(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer scan.f.Close()

	lines := readLines(scan)
	filtered := filter(lines, re)

	if count {
		fmt.Println(len(filtered))
	} else if len(filtered) > 0 {
		printLines(lines, filtered)
	}
}
