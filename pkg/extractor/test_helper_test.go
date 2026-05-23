package extractor

import (
	"os"
	"testing"
)

func skipIfNoDir(t *testing.T, dir string) {
	t.Helper()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skipf("testdata directory %s not available", dir)
	}
}
