package main

import (
	"github.com/ykpythemind/apispec-handler-linter/internal/apirule"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(apirule.Analyzer) }

