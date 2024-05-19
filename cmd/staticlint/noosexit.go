// Package main предоставляет инструмент для статического анализа кода на Go,
// который использует пакет multichecker для организации и запуска различных анализаторов.
// Этот инструмент предназначен для поиска и предотвращения специфических паттернов или ошибок в коде,
// таких как прямой вызов os.Exit в функции main.
//
// Запуск:
//
//	Инструмент можно скомпилировать и запустить, передав директорию или файлы для анализа:
//	go build -o myanalyzer
//	./myanalyzer ./path/to/your/package
//
//	Вы можете также передавать флаги и настройки через командную строку для настройки поведения анализаторов.
//
// Анализаторы:
//
//	noosexit - Анализатор, который проверяет использование os.Exit в функции main.
//	Он предназначен для предотвращения практики прямого выхода из приложения,
//	что может усложнить тестирование и отладку приложений.
package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// noOsExitAnalyzer определяет анализатор для поиска прямых вызовов os.Exit в функции main.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "checks for os.Exit calls in the main function",
	Run:  run,
}

// run выполняет поиск вызовов os.Exit внутри функции main.
// Он итерируется по всем узлам AST и идентифицирует вызовы функций,
// проверяя, находятся ли они в функции main и принадлежат ли они пакету os.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		if strings.Contains(filename, "/go-build/") {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)

			if !ok {
				return true
			}

			if funcDecl.Name.Name == "main" {
				ast.Inspect(funcDecl, func(n ast.Node) bool {
					callExpr, ok := n.(*ast.CallExpr)

					if !ok {
						return true
					}

					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
							pass.Reportf(file.Name.NamePos, "call to os.Exit in main function is prohibited")
						}
					}

					return true
				})
			}

			return true
		})
	}
	return nil, nil
}
