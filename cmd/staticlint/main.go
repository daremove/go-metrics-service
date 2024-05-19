// Package main предоставляет мощный инструмент для статического анализа кода на языке Go.
// Этот инструмент интегрирует широкий спектр анализаторов из пакетов golang.org/x/tools/go/analysis/passes
// и honnef.co/go/tools, предоставляя комплексную проверку качества кода и соблюдения лучших практик разработки.
//
// Список анализаторов включает в себя, но не ограничивается следующими:
// - printf: проверка форматных строк функций fmt.Printf и подобных.
// - shadow: обнаружение переменных, которые могут быть непреднамеренно затенены во время присваивания.
// - structtag: проверка тегов в структурах на соответствие соглашениям Go.
// - asmdecl: проверка согласованности между комментарием Go ассемблера и его кодом.
// - assign: обнаружение неправильного использования оператора присваивания.
// - и многие другие, каждый из которых нацелен на определенный аспект кода или практику программирования.
//
// Запуск анализатора:
//
//	Инструмент может быть скомпилирован и запущен с помощью стандартных команд Go:
//	  go build -o myanalyzer
//	  ./myanalyzer ./path/to/your/code
//
//	Это выполнит все включенные анализаторы на указанных файлах или директориях.
//
// Примечание:
//
//	Разработчики могут добавлять или удалять анализаторы из списка в зависимости от требований проекта,
//	настраивая инструмент для обеспечения наиболее релевантного и эффективного анализа.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/unused"
)

func main() {
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		reflectvaluecompare.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		usesgenerics.Analyzer,
		noOsExitAnalyzer,
	}

	AddAnalyzers := func(analyzers ...*lint.Analyzer) {
		for _, v := range analyzers {
			checks = append(checks, v.Analyzer)
		}
	}

	AddAnalyzers(simple.Analyzers...)
	AddAnalyzers(staticcheck.Analyzers...)
	AddAnalyzers(stylecheck.Analyzers...)
	AddAnalyzers(unused.Analyzer)

	multichecker.Main(
		checks...,
	)
}
