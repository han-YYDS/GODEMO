//go:build wireinject
// +build wireinject

// package main 和 go build 条件编译中间需要有空行
package main

import (
	"github.com/google/wire"
	_ "github.com/google/wire"
)

// Mission的初始化函数
// wire build, 参数是构造函数
func InitMission(name string) Mission {
	wire.Build(NewMonster, NewPlayer, NewMission)
	return Mission{}
}
