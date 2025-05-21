package main

import (
	"fmt"
	"log"

	"github.com/Knetic/govaluate"
)

type User struct {
	FirstName string
	LastName  string
	Age       int
}

func (u User) Fullname() string {
	return u.FirstName + " " + u.LastName
}

func main() {
	test6()
}

// 测试 10 > 0 的判断, 这里是没有参数的
func test1() {
	expr, err := govaluate.NewEvaluableExpression("10 > 0")
	if err != nil {
		log.Fatal("syntax error:", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		log.Fatal("evaluate error:", err)
	}

	fmt.Println(result)
}

// string中留有变量, 同时参数通过map传入
func test2() {
	// 引入变量foo
	expr, _ := govaluate.NewEvaluableExpression("foo > 0")
	// 配置参数 map
	parameters := make(map[string]interface{})
	parameters["foo"] = -1

	// 传入map
	result, _ := expr.Evaluate(parameters)
	fmt.Println(result)

	// 同样地, 传入参数map
	expr, _ = govaluate.NewEvaluableExpression("(requests_made * requests_succeeded / 100) >= 90")
	parameters = make(map[string]interface{})
	parameters["requests_made"] = 100
	parameters["requests_succeeded"] = 80
	result, _ = expr.Evaluate(parameters)
	fmt.Println(result)

	// 同
	expr, _ = govaluate.NewEvaluableExpression("(mem_used / total_mem) * 100")
	parameters = make(map[string]interface{})
	parameters["total_mem"] = 1024
	parameters["mem_used"] = 512
	result, _ = expr.Evaluate(parameters)
	fmt.Println(result)
}

// 变量命名可以不合golang的规范
func test3() {
	// 对于goland中不允许的变量名(不允许有-), 也是能够进行替换的
	expr, _ := govaluate.NewEvaluableExpression("[response-time] < 100")
	parameters := make(map[string]interface{})
	parameters["response-time"] = 80
	result, _ := expr.Evaluate(parameters)
	fmt.Println(result)

	//
	expr, _ = govaluate.NewEvaluableExpression("response\\-time < 100")
	parameters = make(map[string]interface{})
	parameters["response-time"] = 80
	result, _ = expr.Evaluate(parameters)
	fmt.Println(result)
}

// expr的复用
func test4() {

	expr, _ := govaluate.NewEvaluableExpression("a + b")
	parameters := make(map[string]interface{})
	parameters["a"] = 1
	parameters["b"] = 2
	result, _ := expr.Evaluate(parameters)
	fmt.Println(result)

	parameters = make(map[string]interface{})
	parameters["a"] = 10
	parameters["b"] = 20
	result, _ = expr.Evaluate(parameters)
	fmt.Println(result)
}

// 自定义函数
func test5() {
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return length, nil
		},
	}

	exprString := "strlen('teststring')"

	// 传入functions map, 实现自定义的map
	expr, _ := govaluate.NewEvaluableExpressionWithFunctions(exprString, functions)
	result, _ := expr.Evaluate(nil)
	fmt.Println(result)
}

func test6() {
	u := User{FirstName: "li", LastName: "dajun", Age: 18}
	parameters := make(map[string]interface{})
	parameters["u"] = u

	expr, err := govaluate.NewEvaluableExpression("u.Fullname()")
	if err != nil {
		fmt.Println("expr: ", err)
		return
	}
	result, _ := expr.Evaluate(parameters)
	fmt.Println("user", result)

	// expr, _ = govaluate.NewEvaluableExpression("u.Age > 18")
	// result, _ = expr.Evaluate(parameters)
	// fmt.Println("age > 18?", result)
}
