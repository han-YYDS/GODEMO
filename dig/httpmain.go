package main

import (
	"fmt"
	"net/http"

	"go.uber.org/dig"
)

type Handler struct {
	Greeting string
	Path     string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s from %s", h.Greeting, h.Path)
}
func NewHello1Handler() HandlerResult {
	// 输出到server组, 并被该组的Params收集
	return HandlerResult{
		Handler: Handler{
			Path:     "/hello1",
			Greeting: "welcome",
		},
	}
}
func NewHello2Handler() HandlerResult {
	return HandlerResult{
		Handler: Handler{
			Path:     "/hello2",
			Greeting: "😄",
		},
	}
}

type HandlerResult struct { // 输出了该实例的函数都是server的
	dig.Out
	Handler Handler `group:"server"` // Handler添加到server
}
type HandlerParams struct { // 输入结构体会自动收集Handlers,而且会筛选出server组的
	dig.In
	Handlers []Handler `group:"server"` // 接收server组的所有Handler
}

func RunServer(params HandlerParams) error {
	mux := http.NewServeMux() // http请求多路复用, 将不同url请求发送到不同handler中
	for _, h := range params.Handlers {
		// 遍历所有handler(dig注入,使用到了参数params)
		mux.Handle(h.Path, h)
	}
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
func main() {
	container := dig.New()
	container.Provide(NewHello1Handler) // handler1 -> server
	container.Provide(NewHello2Handler) // handler2 -> server
	container.Invoke(RunServer)
	// 运行时发现需要HandlerParams

}
