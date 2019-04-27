package promql

import (
	"path/filepath"
	"testing"
)

func TestEvaluations(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	files, err := filepath.Glob("testdata/*.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range files {
		test, err := newTestFromFile(t, fn)
		if err != nil {
			t.Errorf("error creating test for %s: %s", fn, err)
		}
		err = test.Run()
		if err != nil {
			t.Errorf("error running test %s: %s", fn, err)
		}
		test.Close()
	}
}
