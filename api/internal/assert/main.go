package assert

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
		t.FailNow()
	}
}

func Contains(t *testing.T, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Errorf("%s doesn't contain %s", str, substr)
		t.FailNow()
	}
}

func IsUUID(t *testing.T, str string) {
	t.Helper()

	_, err := uuid.Parse(str)
	if err != nil {
		t.Errorf("%s is not a UUID", str)
		t.FailNow()
	}
}
