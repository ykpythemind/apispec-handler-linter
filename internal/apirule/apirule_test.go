package apirule_test

import (
	"testing"

	"github.com/ykpythemind/apispec-handler-linter/internal/apirule"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, apirule.Analyzer, "a")
}

