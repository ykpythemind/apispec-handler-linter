package apirule

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "apirule",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "apirule is ..."

var once = new(sync.Once)
var operationIdSet = make(map[string]struct{})

func parseOperationIds() map[string]struct{} {
	// run(pass *analysis.Pass)は複数回呼ばれるので最初の一回だけパースを実行する
	once.Do(func() {
		f, err := os.Open("apispec/api.gen.go") // apispecから自動生成したファイル
		if err != nil {
			panic(err)
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		fset := token.NewFileSet()
		specfile, err := parser.ParseFile(fset, "", b, parser.Mode(0))
		if err != nil {
			panic(err)
		}

		ast.Inspect(specfile, func(n ast.Node) bool {
			if typespec, ok := n.(*ast.TypeSpec); ok {
				if typespec.Name.Name == "ServerInterface" {
					if iftype, ok := typespec.Type.(*ast.InterfaceType); ok {
						for _, id := range getMethodNamesOfInterfaceType(iftype) {
							operationIdSet[id] = struct{}{}
						}
					}
				}
			}
			return true
		})
	})

	return operationIdSet
}

func getMethodNamesOfInterfaceType(iftype *ast.InterfaceType) (list []string) {
	// いくつかのnilチェックを省略しているが、自動生成コードなので必ず存在するはずである
	for _, method := range iftype.Methods.List {
		name := method.Names[0].Name
		list = append(list, name)
	}
	return
}

func run(pass *analysis.Pass) (interface{}, error) {
	operationIds := parseOperationIds()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// oyakusoku
	if pass.Pkg.Name() != "handler" {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)

		_, ok := operationIds[fn.Name.Name]
		if !ok {
			// funcがoperationId名ではなかったらなにもしない
			return
		}
		opId := fn.Name.Name

		if fn.Body == nil {
			return
		}

		if err := analyze(fn.Body, opId); err != nil {
			// 違反を検出したらReport()を使って報告する
			pass.Report(analysis.Diagnostic{
				Pos:     fn.Body.Pos(),
				End:     fn.Body.End(),
				Message: err.Error(),
			})
		}
	})

	return nil, nil
}

// analyze は与えられた関数の中でoperation idと紐づくrequest, responseの型を使用しているかどうか確認し、使用していない場合ルール違反としてエラーを返す.
func analyze(stmt ast.Stmt, opId string) error {
	responseTypeUsed := false
	requestTypeUsed := false

	ast.Inspect(stmt, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if ident.Name == fmt.Sprintf("%sResponse", opId) {
				responseTypeUsed = true
			} else if ident.Name == fmt.Sprintf("%sRequest", opId) {
				requestTypeUsed = true
			}
		}
		return true
	})

	if responseTypeUsed && requestTypeUsed {
		return nil // ok
	}

	err := fmt.Errorf("func %s() で使う必要のあるstructを使っていません.", opId)

	if !responseTypeUsed {
		err = fmt.Errorf("%w %sResponse を使用していません", err, opId)
	}

	if !requestTypeUsed {
		err = fmt.Errorf("%w %sRequest を使用していません", err, opId)
	}

	return err
}
