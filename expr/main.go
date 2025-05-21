package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

// 表达式计算
func test1() {
	// env设置变量申明
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	// Eval进行计算
	out, err := expr.Eval("foo + bar", env)

	if err != nil {
		fmt.Println("err: ", err)
		panic(err)
	}

	// 获得结果
	fmt.Println(out)
}

// 代码片段
func test2() {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",             // greet是一个fmt格式的string
		"names":   []string{"world", "you"}, // names是一个slice
		"sprintf": fmt.Sprintf,              // You can pass any functions. // sprintf是一个函数,没有括号
	}

	code := `sprintf(greet, names[0])` // 如果直接执行env的字符串替换倒也是合理

	// Compile code into bytecode. This step can be done once and program may be reused.
	// Specify environment for type check.
	// Complie将string -> exe 相当于提前编译成二进制, 在之后可以对该程序进行Run复用
	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	// 复用program
	env["names"] = []string{"WORLD", "you"}
	output, err = expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

type Env struct {
	Tweets []Tweet
}

// Methods defined on such struct will be functions.
func (Env) Format(t time.Time) string { return t.Format(time.RFC822) }

// 一条推特
type Tweet struct {
	Text string
	Date time.Time
}

// 自定义结构体, 调用对象方法, 运用Expr自带的函数
func test3() {
	// 这里的map和filter都是函数
	// map相当于range
	// filter用于筛选
	// 这里的.Text是对Tweet结构体的调用, 隐含了句柄
	code := `map(
			filter(Tweets, {len(.Text) > 0}), 
			{.Text + Format(.Date)}
	)`

	// We can use an empty instance of the struct as an environment.
	// ENV用空的
	program, err := expr.Compile(code, expr.Env(Env{}))
	if err != nil {
		panic(err)
	}

	// env赋值
	env := Env{
		Tweets: []Tweet{{"Oh My God!", time.Now()}, {"How you doin?", time.Now()}, {"Could I be wearing any more clothes?", time.Now()}},
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

// 自定义函数的传入
func test4() {
	// interface{} 可以是一个 func,
	env := map[string]interface{}{
		"foo":    1,
		"double": func(i int) int { return i * 2 },
	}

	out, err := expr.Eval("double(foo)", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}

type Env1 map[string]interface{}

// 只能是外部函数才可以
func (Env1) FastMethod(...interface{}) interface{} {
	return "Hello, "
}

// 快速函数可以不使用反射调用, 但是只支持 func(...interface{}) interface{} 这样的函数声明, 返回值可以多一个error
// 快速函数能够提升性能, 但是丢失了类型
// 快速函数通过指针调用, 绕过了反射,
func test5() {
	env := Env1{
		"fast_func": func(...interface{}) interface{} { return "world" },
	}

	// 这里的Fast是 env的成员函数, 这个是需要大写外部的
	// fast是env的成员
	out, err := expr.Eval("FastMethod() + fast_func()", env)

	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}

func test6() {
	env := Env1{
		"fast_func": func(...interface{}) interface{} { return "world" },
	}

	// 通过Function能够弥补快速函数的参数类型丢失的缺陷
	atoi := expr.Function(
		"atoi", // 函数名称
		func(params ...any) (any, error) { // 匿名函数声明
			return strconv.Atoi(params[0].(string))
		},
		new(func(string) int),
	)

	program, err := expr.Compile(`atoi("42")`, atoi)
	if err != nil {
		panic(err)
	}
	expr.Run(program, env)
}

// 需要返回错误
func test7() {
	env := map[string]interface{}{
		"foo": -1,
		"Double": func(i int) (int, error) { // 返回错误的函数
			if i < 0 {
				return 0, errors.New("[DEMO]value cannot be less than zero")
			}
			return i * 2, nil
		},
	}

	// 也是用两个值接住返回值
	out, err := expr.Eval("Double(foo)", env)

	// This `err` will be the one returned from `double` function.
	// err.Error() == "value cannot be less than zero"
	fmt.Println("err: ", err)
	fmt.Print(out)
}

// Env->struct->func的调用路径
type Env2 struct {
	datetime  // 组合
	CreatedAt time.Time
}

// Functions may be defined on embedded structs as well.
type datetime struct{}

func (datetime) Now() time.Time                   { return time.Now() }
func (datetime) Sub(a, b time.Time) time.Duration { return a.Sub(b) }

func test8() {
	code := `(Now() - CreatedAt).Hours() / 24 / 365` // 这里用的 - 是负载后的

	// We can define options before compiling.
	options := []expr.Option{
		expr.Env(Env2{}),          // Expr 编译选项之一，用于指定表达式求值时的环境。通过传入 Env2{}，告诉 Expr 在求值时可以使用 Env2 结构体中的字段和方法。
		expr.Operator("-", "Sub"), // 重载了 operator- 运算符, 这里的Sub是成员函数
	}

	program, err := expr.Compile(code, options...)
	if err != nil {
		panic(err)
	}

	// 编译之后才开始传参
	env := Env2{
		CreatedAt: time.Date(1987, time.November, 24, 20, 0, 0, 0, time.UTC),
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

type visitor struct {
	identifiers []string
}

func (v *visitor) Visit(node *ast.Node) {
	if n, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers = append(v.identifiers, n.Value)
	}
}

// 通过AST对code进行遍历
func test9() {
	tree, err := parser.Parse("foo + bar")
	if err != nil {
		panic(err)
	}

	visitor := &visitor{}
	ast.Walk(&tree.Node, visitor)

	fmt.Printf("%v", visitor.identifiers) // outputs [foo bar]
}

type patcher struct{}

func (p *patcher) Enter(_ *ast.Node) {}
func (p *patcher) Exit(node *ast.Node) {
	n, ok := (*node).(*ast.IndexNode)
	if !ok {
		return
	}
	unary, ok := n.Index.(*ast.UnaryNode)
	if !ok {
		return
	}
	if unary.Operator == "-" {
		ast.Patch(&n.Index, &ast.BinaryNode{
			Operator: "-",
			Left:     &ast.BuiltinNode{Name: "len", Arguments: []ast.Node{n.Node}},
			Right:    unary.Node,
		})
	}

}

func test10() {
	env := map[string]interface{}{
		"list": []int{1, 2, 3},
	}

	code := `list[-1]` // will output 3

	program, err := expr.Compile(code, expr.Env(env), expr.Patch(&patcher{}))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}
	fmt.Print(output)
}

func main() {
	test10()
}
