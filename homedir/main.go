// 在 Darwin 系统上，标准库os/user的使用需要 cgo。所以，任何使用os/user的代码都不能交叉编译。
// 但是，大多数人使用os/user的目的仅仅只是想获取主目录。因此，go-homedir库出现了。

package main

import (
	"fmt"
	"os/user"

	"github.com/mitchellh/go-homedir"
)

// os/user
func test1() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("home: ", u.HomeDir)
}

func test2() {
	homedir.DisableCache = false
	homedir.Reset() // 调用缓存时, git reset 获得的不一定是最新目录, 需要reset重置
	path, _ := homedir.Dir() // Dir 获取用户主目录
	fmt.Println("path: ", path)

	dir := "~/golang/src"
	ex, _ := homedir.Expand(dir)
	fmt.Println("ex: ", ex)
}

func main() {
	test2()
}
