package main

import (
	"mychecker/exitchecker"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var mychecks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	mychecks = append(mychecks, exitchecker.ExitCheckAnalyzer)
	mychecks = append(mychecks, defers.Analyzer)       // проверяет наличие ошибок в операторах defer
	mychecks = append(mychecks, nilfunc.Analyzer)      // проверяет бесполезность сравнений с nil
	mychecks = append(mychecks, nilness.Analyzer)      // проверяет разыменование указателей nil и сравнения указателей nil
	mychecks = append(mychecks, printf.Analyzer)       // проверяет согласованность строк и аргументов формата printf
	mychecks = append(mychecks, assign.Analyzer)       // выявляет бесполезные назначения
	mychecks = append(mychecks, shift.Analyzer)        // проверяет наличие сдвигов, превыщающих ширину целого числа
	mychecks = append(mychecks, structtag.Analyzer)    // проверяет правильность формирования тегов полей struct
	mychecks = append(mychecks, timeformat.Analyzer)   // проверяет использование вызовов time.Format или time.Parse с плохим форматом
	mychecks = append(mychecks, unmarshal.Analyzer)    // проверяет передачу в функции unmarshal и decode типов, не являющихся указателями или не являющихся интерфейсами
	mychecks = append(mychecks, unusedresult.Analyzer) // проверяет наличие неиспользованных результатов вызовов некоторых чистых функций.

	multichecker.Main(
		mychecks...,
	)
}
