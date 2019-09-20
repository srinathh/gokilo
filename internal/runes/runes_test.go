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
