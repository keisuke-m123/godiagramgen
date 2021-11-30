package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func LogDiff(t *testing.T, want interface{}, got interface{}) {
	t.Helper()

	t.Logf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
}
