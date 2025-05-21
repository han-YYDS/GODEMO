package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	// 创建文件集用于定位信息
	fset := token.NewFileSet()

	// 示例Go源代码
	src := `
package main

import "fmt"

func greet(name string) {
	fmt.Println("Hello", name)
}

func main() {
	x := 10
	y := x * 2
	greet("Gopher")
}
`

	// 解析源代码生成AST
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		log.Fatal("解析错误:", err)
	}

	// 示例1：使用ast.Print打印完整AST结构
	log.Println("==== 完整AST结构 ====")
	ast.Print(fset, f)

	// 示例2：使用ast.Inspect遍历函数声明
	log.Println("\n==== 函数分析 ====")
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			log.Printf("发现函数: %s (位置: %s)",
				x.Name.Name,
				fset.Position(x.Pos()))
		}
		return true
	})

	// 示例3：使用ast.Inspect分析变量声明
	log.Println("\n==== 变量分析 ====")
	ast.Inspect(f, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					log.Printf("变量声明: %s (位置: %s)",
						ident.Name,
						fset.Position(ident.Pos()))
				}
			}
		}
		return true
	})

	// 示例4：使用ast.Fprint自定义输出
	log.Println("\n==== 自定义AST输出 ====")
	ast.Fprint(log.Writer(), fset, f, ast.NotNilFilter)
}
