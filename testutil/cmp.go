package testutil

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Diff(t *testing.T, want interface{}, got interface{}) string {
	t.Helper()

	return fmt.Sprintf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
}
