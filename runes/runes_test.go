package runes

import (
	"testing"
)

func TestEqual(t *testing.T) {

	testCases := []struct {
		A    []rune
		B    []rune
		want bool
	}{
		{[]rune("abc"), []rune("abc"), true},
		{[]rune("abc"), []rune("def"), false},
	}

	for _, testCase := range testCases {
		if Equal(testCase.A, testCase.B) != testCase.want {
			t.Errorf("Mismatch: %v", testCase)
		}
	}
}

func TestIndex(t *testing.T) {
	testCases := []struct {
		S    []rune
		Sub  []rune
		Want int
	}{
		{[]rune("abc"), []rune("b"), 1},
		{[]rune("abc"), []rune("d"), -1},
		{[]rune("abc def ghi"), []rune("de"), 4},
	}

	for _, testCase := range testCases {
		if Index(testCase.S, testCase.Sub) != testCase.Want {
			t.Errorf("Mismatch: %s:%s", string(testCase.S), string(testCase.Sub))
		}
	}
}
