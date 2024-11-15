package assert

import "testing"

func Eq[T comparable](t *testing.T, x, y T) {
	t.Helper()
	if x != y {
		t.Errorf("expected %v, got %v", y, x)
	}
}
