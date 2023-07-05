package Geektutu_leaning

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//	func indexHandler(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintf(w, "欢迎啊， 处理器为:indexhandler, 路由为%q\n", req.URL.Path)
//	}
//
//	func helloHandler(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintf(w, "欢迎啊， 处理器为:hellohandler, 路由为%q\n", req.URL.Path)
//	}
//
//	func main() {
//		http.HandleFunc("/", indexHandler)
//		http.HandleFunc("/hello", helloHandler)
//		if err := http.ListenAndServe(":8080", nil); err != nil {
//			log.Fatal(err)
//		}
//	}
/*
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
type server interface {
	http.Handler
	Start(addr string) error
	Stop() error
}

/*
	一个Server需要的功能是

开启和关闭,这是为了兼容性，因为不能就只写http，后续可能还有https
*/
// Option设计模式
type HTTPOption func(h *HTTPServer)

type HTTPServer struct {
	srv  *http.Server
	stop func() error
}

func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil { //实现一个优雅关闭逻辑
			fn = func() error {
				fmt.Println("程序正常启动")
				quit := make(chan os.Signal)
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit // 阻塞住
				log.Println("Shutdown Server ...")

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				// 关闭之前需要做某些操作
				if err := h.srv.Shutdown(ctx); err != nil {
					log.Fatal("Server Shutdown", err)
				}
				// 关闭之后需要做的操作
				select {
				case <-ctx.Done():
					log.Println("timeout of 5 seconds")
				}
				return nil
			}
		}
		h.stop = fn

	}
}
func NewHTTP(opts ...HTTPOption) *HTTPServer {
	h := &HTTPServer{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// 接收请求转发请求
// 接收前端传过来的请求
// 转发请求：转发前端发过来的请求到咱们的框架中
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	panic("implement me")
}

// 启动服务
func (h *HTTPServer) Start(addr string) error {
	h.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return h.srv.ListenAndServe()

}

// 停止服务
func (h *HTTPServer) Stop() error {
	return h.stop()
}

func main() {
	h := NewHTTP(WithHTTPServerStop(nil))
	go func() {
		if err := h.Start(":8080"); err != nil && err != http.ErrServerClosed {
			//h.Fail()
			panic("启动失败")
		}
	}()
	err := h.Stop()
	if err != nil {
		panic("关闭失败")
	}
}
