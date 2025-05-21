package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// 20250211-10:27 reflect demo

func testKind() {
	for _, v := range []any{"hi", 42, func() {}} {
		switch v := reflect.ValueOf(v); v.Kind() {
		case reflect.String:
			fmt.Println(v.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Println(v.Int())
		default:
			fmt.Printf("unhandled kind %s", v.Kind())
		}
	}

}

func test2() {
	type S struct {
		F0 string `alias:"field_0"`
		F1 string `alias:""`
		F2 string
	}

	s := S{}
	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if alias, ok := field.Tag.Lookup("alias"); ok {
			if alias == "" {
				fmt.Println("(blank)")
			} else {
				fmt.Println(alias)
			}
		} else {
			fmt.Println("(not specified)")
		}
	}

}

func test3() {
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Height",
			Type: reflect.TypeOf(float64(0)),
			Tag:  `json:"height"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(int(0)),
			Tag:  `json:"age"`,
		},
	})

	v := reflect.New(typ).Elem()
	v.Field(0).SetFloat(0.4)
	v.Field(1).SetInt(2)
	s := v.Addr().Interface()

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

	fmt.Printf("value: %+v\n", s)
	fmt.Printf("json:  %s", w.Bytes())

	r := bytes.NewReader([]byte(`{"height":1.5,"age":10}`))
	if err := json.NewDecoder(r).Decode(s); err != nil {
		panic(err)
	}
	fmt.Printf("value: %+v\n", s)

}

type ms string

type cat struct {
}

func test4() {
	var a ms
	var b cat

	typeOfA := reflect.TypeOf(a)         // 获取类型对象
	fmt.Println(typeOfA, typeOfA.Kind()) // 类型名称与类型种类

	typeOfCat := reflect.TypeOf(b)           // 获取结构体实例的反射类型对象
	fmt.Println(typeOfCat, typeOfCat.Kind()) // 显示反射类型对象的名称和种类
}

// 用reflect包对一个Employee类型的值进行遍历，要求输出字段的名称、类型和值。
type Employee struct {
	Id   int32
	Name string
}

func test5() {
	var e = Employee{
		Id:   1,
		Name: "demo",
	}
	v := reflect.ValueOf(e)
	t := reflect.TypeOf(e)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := range t.NumField() {
		tf := t.Field(i)
		vf := v.Field(i)
		fmt.Println("\n字段名称: ", tf.Name)
		fmt.Println("字段类型: ", tf.Type)
		if tf.PkgPath != "" {
			fmt.Println("字段值: ", vf.Interface())
		}

	}
}

type Tree struct {
	E1   Employee
	name string
	e2   *Employee
}

func test6() {
	var e = Employee{
		Id:   1,
		Name: "demo",
	}

	var tree = &Tree{
		E1:   e,
		name: "treee",
		e2:   &e,
	}

	PrintStruct(tree)
}

func PrintStruct(i interface{}) {
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	fmt.Println("结构体名称: ", t.Name())
	fmt.Println("结构体类型: ", t.String())

	for i := range t.NumField() {
		fmt.Println("\n field : ", i)
		tf := t.Field(i)
		vf := v.Field(i)

		fmt.Println("字段名称: ", tf.Name)
		fmt.Println("字段类型: ", tf.Type)
		if tf.PkgPath != "" { // 未导出
			fmt.Println("未导出")
			continue
		}
		if vf.CanInterface() {
			fmt.Println("字段值:", vf.Interface())
		} else {
			fmt.Println("字段值: [无法获取]")
		}

		switch tf.Type.Kind() {
		case reflect.Ptr:
			if !vf.IsNil() { // 空指针检查
				PrintStruct(vf.Interface())
			}
		case reflect.Struct:
			PrintStruct(vf.Interface())
		}

	}
}

type Em struct {
	Name   string  `key:"name"`
	Role   string  `key:"role"`
	Salary float64 `key:"salary"`
}

// func rebuiltStruct(mapData map[string]interface{}, target interface{}) {
// 	sVal := reflect.ValueOf(target)
// 	sType := reflect.TypeOf(target)
// 	if sType.Kind() == reflect.Ptr {
// 		sVal = sVal.Elem()
// 		sType = sType.Elem()
// 	}
// 	num := sVal.NumField()
// 	for i := 0; i < num; i++ {
// 		f := sType.Field(i)
// 		val := sVal.Field(i)
// 		key := f.Tag.Get("key")
// 		if dataVal, ok := mapData[key]; ok {
// 			//类型判断与转换
// 			dataType := reflect.TypeOf(dataVal)
// 			fieldType := val.Type()
// 			if dataType == fieldType {
// 				val.Set(reflect.ValueOf(dataVal))
// 			} else {
// 				if dataType.ConvertibleTo(fieldType) {
// 					val.Set(reflect.ValueOf(dataVal).Convert(fieldType))
// 				} else {
// 					panic(fmt.Sprintf("failed to convert from %s to %s \n", dataType, fieldType))
// 				}
// 			}
// 		} else {
// 			fmt.Printf("key %s not found in struct definition! \n", key)
// 		}
// 	}
// 	traverse2(target)
// }

func traverse2(target interface{}) {
	sVal := reflect.ValueOf(target)
	sType := reflect.TypeOf(target)
	if sType.Kind() == reflect.Ptr {
		sVal = sVal.Elem()
		sType = sType.Elem()
	}
	num := sVal.NumField()
	for i := 0; i < num; i++ {
		//判断字段是否为结构体类型，或者是否为指向结构体的指针类型
		if sVal.Field(i).Kind() == reflect.Struct || (sVal.Field(i).Kind() == reflect.Ptr && sVal.Field(i).Elem().Kind() == reflect.Struct) {
			traverse2(sVal.Field(i).Interface())
		} else {
			f := sType.Field(i)
			val := sVal.Field(i).Interface()
			fmt.Printf("%5s %v = %v\n", f.Name, f.Type, val)
		}
	}
}

// 给定一个map值，我们根据该map提供的信息，恢复构建出一个Employee类型的值
func rebuiltStruct(mapData map[string]interface{}, target interface{}) {
	v := reflect.ValueOf(target)
	t := reflect.TypeOf(target)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := range v.NumField() {
		vf := v.Field(i)
		tf := t.Field(i)
		key := tf.Tag.Get("key") // kv存储在map中
		mapvalue, ok := mapData[key]
		if !ok {
			continue
		}
		// 判断类型是否一致
		if reflect.TypeOf(mapvalue) == tf.Type {
			vf.Set(reflect.ValueOf(mapvalue))
		} else {
			fmt.Println("///// convert ", tf.Name, " from ", reflect.TypeOf(mapvalue).Name(), " to ", tf.Type)
			if reflect.TypeOf(mapvalue).ConvertibleTo(tf.Type) {
				vf.Set(reflect.ValueOf(mapvalue).Convert(tf.Type))
			}
		}
		// vf.Set(reflect.ValueOf(value))
	}

	PrintStruct(target)
}

var employeeData = map[string]interface{}{
	"name":   "laozhang",
	"role":   "annother glory engineer",
	"salary": 1,
}

func test7() {
	rebuiltStruct(employeeData, &Em{})
}

// 调用struct的方法
func (e Employee) Code() {
	fmt.Printf("I like to code \n")
}

// func (e Employee) Debug(i int) error {
// 	fmt.Printf("I dislike to debug \n")
// 	return nil
// }

func (e Employee) raiseSalary() {
	fmt.Printf("I want to raise my salasy \n")
}

func callMethodWithReflect(x interface{}) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("num: ", v.NumMethod())
	for i := range v.NumMethod() {
		vm := v.Method(i) // method value
		tm := t.Method(i) // method type
		fmt.Println("Name: ", tm.Name)
		fmt.Println("type: ", tm.Type)

		args := []reflect.Value{}
		vm.Call(args)
	}

}

func test8() {
	callMethodWithReflect(&Employee{})
}

// 按方法名调用
func (e Employee) Work(i int) (error, int) {
	fmt.Printf("I work for %v hours per day. \n", i)
	return nil, 123
}

func callByMethodName(x interface{}, methodName string) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := v.MethodByName(methodName)
	if m.IsValid() { // 通过valid判断是否找到method
		res := m.Call([]reflect.Value{reflect.ValueOf(10)})
		for _, re := range res {
			fmt.Println("re: ", re.Interface())
		}
	}
}

func test9() {
	callByMethodName(Employee{}, "Work")
}

func main() {
	test9()
}
