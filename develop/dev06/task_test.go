package dev06

import (
	"math"
	"reflect"
	"testing"
)

func TestParseFieldsCorrectInput(t *testing.T) {
	s := "-2,5,7-9,12-"
	expected := [][]int{{1, 2}, {5, 5}, {7, 9}, {12, math.MaxInt}}

	parsed, err := parseFields(s)

	if err != nil || !reflect.DeepEqual(parsed, expected) {
		t.Logf("parseFields(%s) = %v, %v, expected: %v, nil", s, parsed, err, expected)
	}
}

func TestParseFieldsIncorrectInputDecreasingRange(t *testing.T) {
	s := "5-2"

	parsed, err := parseFields(s)

	if err == nil {
		t.Logf("parseFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestParseFieldsIncorrectInputFormat(t *testing.T) {
	s := "hello"

	parsed, err := parseFields(s)

	if err == nil {
		t.Logf("parseFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestParseFieldsIncorrectInputComma(t *testing.T) {
	s := ","

	parsed, err := parseFields(s)

	if err == nil {
		t.Logf("parseFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestParseFieldsIncorrectInputPreComma(t *testing.T) {
	s := ",2"

	parsed, err := parseFields(s)

	if err == nil {
		t.Logf("parseFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestParseFieldsIncorrectInputPostComma(t *testing.T) {
	s := "2,"

	parsed, err := parseFields(s)

	if err == nil {
		t.Logf("parseFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestIndexInSegmentsTrue(t *testing.T) {
	i := 8
	segs := [][]int{{7, 9}}

	res := indexInSegments(i, segs)

	if !res {
		t.Logf("indexInSegments(%d, %v) = %v, expected: true", i, segs, res)
	}
}

func TestIndexInSegmentsFalse(t *testing.T) {
	i := 5
	segs := [][]int{{7, 9}}

	res := indexInSegments(i, segs)

	if res {
		t.Logf("indexInSegments(%d, %v) = %v, expected: true", i, segs, res)
	}
}
