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
// 视图函数签名，到时候可以封装成上下文的形式
// 这种方式会导致无法扩展，因此需要用上下文来进行抽象
//type HandleFunc func(w http.ResponseWriter, req *http.Request) //抽象一个处理函数
type HandleFunc func(*Context)
type server interface {
	http.Handler
	Start(addr string) error
	Stop() error
	// 注册路由，一个非常核心的API，不能给开发者乱用
	// 造一些衍生API给开发者使用
	addRoute(method string, path string, handleFunc HandleFunc)
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
	// routers 临时存在路由的位置
	router map[string]HandleFunc
}

/*
{
	"GET-login": HandleFunc1,
	"POST-login": HandleFunc2,

}
*/

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
	h := &HTTPServer{
		router: map[string]HandleFunc{},
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// 接收请求转发请求
// 接收前端传过来的请求
// 转发请求：转发前端发过来的请求到咱们的框架中
// ServerHTTP向前对接前端请求，向后对接框架
// 前端发请求给ServerHTTP， 后端处理后直接根据这个方法发送给前端
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 1.匹配路由
	ctx := newContext(w, req)
	key := req.Method + "-" + req.URL.Path
	if handler, ok := h.router[key]; ok { // 如果对应的key存在handler
		handler(ctx) //转发请求
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404 NOT FOUND\n"))
	}

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

// 注册路由的时机：项目启动的时候，后续就不能注册路由了
// 注册路由放在哪里？--->有前缀树放前缀树，没前缀树先放map里面，实现一个静态路由匹配
func (h *HTTPServer) addRouter(method string, pattern string, handleFunc HandleFunc) {
	key := method + "-" + pattern                        // "GET-login" 目的是要唯一
	fmt.Printf("add router %s - %s \n", method, pattern) // method表示的是方法GET PUT DELETE AND POST，
	//pattern表示自定义的匹配格式
	h.router[key] = handleFunc //注册完毕, 每个路由对应一个HandleFunc
}

// GET
func (h *HTTPServer) GET(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodGet, pattern, handleFunc)
}

// POST
func (h *HTTPServer) POST(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodPost, pattern, handleFunc)
}

// PUT
func (h *HTTPServer) PUT(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodPut, pattern, handleFunc)
}

// DELETE
func (h *HTTPServer) DELETE(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodDelete, pattern, handleFunc)
}

//func main() {
//	h := NewHTTP(WithHTTPServerStop(nil))
//	go func() {
//		if err := h.Start(":8080"); err != nil && err != http.ErrServerClosed {
//			//h.Fail()
//			panic("启动失败")
//		}
//	}()
//	err := h.Stop()
//	if err != nil {
//		panic("关闭失败")
//	}
//}
