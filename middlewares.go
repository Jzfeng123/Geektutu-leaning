package Geektutu_learning

// MiddlewareHandleFunc 中间件的函数签名
// 参数next表示下一次执行的中间件的逻辑
// 返回值表示当前中间件的逻辑
// Gin框架和Geektutu中都是将中间件和视图函数视为是一个整体，而并没有像这样在抽象出来一个,这是一种责任链的实现机制，他们是一种洋葱机制
type MiddlewareHandleFunc func(next HandleFunc) HandleFunc
