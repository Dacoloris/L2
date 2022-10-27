package dev02

import "testing"

func TestPrimitiveExtract(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "test-0",
			arg:  "a4bc2d5e",
			want: "aaaabccddddde",
		},
		{
			name: "test-1",
			arg:  "abcd",
			want: "abcd",
		},
		{
			name: "test-2",
			arg:  "45",
			want: "некорректная строка",
		},
		{
			name: "test-3",
			arg:  "",
			want: "",
		},
		{
			name: "test-4",
			arg:  "a44",
			want: "некорректная строка",
		},
		{
			name: "test-5",
			arg:  "b56c",
			want: "некорректная строка",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := primitiveExtract(test.arg); got != test.want {
				t.Errorf("primitiveExtract() = %v, want = %v", got, test.want)
			}
		})
	}
}
