package dev04

import (
	"reflect"
	"testing"
)

func TestAnagramSearch(t *testing.T) {
	tests := []struct {
		name string
		args *[]string
		want map[string][]string
	}{
		{
			name: "test-0",
			args: &[]string{""},
			want: map[string][]string{},
		},
		{
			name: "test-1",
			args: &[]string{"пятка"},
			want: map[string][]string{},
		},
		{
			name: "test-2",
			args: &[]string{"пятак", "пятка"},
			want: map[string][]string{
				"пятак": {"пятка"},
			},
		},
		{
			name: "test-3",
			args: &[]string{"пятка", "пятак", "листок", "тяпка", "слиток", "столик"},
			want: map[string][]string{
				"пятка":  {"пятак", "тяпка"},
				"листок": {"слиток", "столик"},
			},
		},
		{
			name: "test-4",
			args: &[]string{"пятка", "пятак", "кеп", "листок", "тяпка", "впр", "слиток", "столик", "в"},
			want: map[string][]string{
				"пятка":  {"пятак", "тяпка"},
				"листок": {"слиток", "столик"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anagramSearch(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("anagramSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
