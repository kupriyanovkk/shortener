package main

import (
	"go/ast"
	"strings"

	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

// Custom analyzer for os.Exit call in main func
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit call in main func",
	Run:  run,
}

// run parses the source file and defines a custom visitor to check for os.Exit in the main function.
//
// Parameter: pass *analysis.Pass
// Return: (interface{}, error)
func run(pass *analysis.Pass) (interface{}, error) {
	// Get the main package
	pkg := pass.Pkg
	if pkg.Name() != "main" {
		return nil, nil // Skip if it's not the main package
	}
	for _, file := range pass.Files {
		// Check if it's the main function
		for _, decl := range file.Decls {
			fn, isFunc := decl.(*ast.FuncDecl)
			if isFunc && fn.Name.Name == "main" {
				// Inspect the function body for os.Exit calls
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					callExpr, isCallExpr := n.(*ast.CallExpr)
					if isCallExpr {
						selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
						if isSelector {
							ident, isIdent := selectorExpr.X.(*ast.Ident)
							if isIdent && ident.Name == "os" && selectorExpr.Sel.Name == "Exit" {
								// Report the use of os.Exit in the main function
								pass.Reportf(callExpr.Pos(), "Direct os.Exit call detected in main function")
							}
						}
					}
					return true
				})
			}
		}
	}
	return nil, nil
}

func main() {
	// standard static analyzers from package golang.org/x/tools/go/analysis/passes
	analyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		stringintconv.Analyzer,
		ifaceassert.Analyzer,

		bodyclose.Analyzer,
		rowserr.NewAnalyzer(
			"github.com/jmoiron/sqlx",
		),

		OsExitAnalyzer,
	}

	for _, v := range simple.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	// analyzers of class SA from package staticcheck.io
	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(analyzers...)
}
