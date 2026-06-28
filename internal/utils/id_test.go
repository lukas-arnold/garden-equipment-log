package utils

import "testing"

func TestIdIsDifferent(t *testing.T) {

	first := Id()
	second := Id()

	if first == second {
		t.Fatalf(
			"ids are equal: %d",
			first,
		)
	}
}
