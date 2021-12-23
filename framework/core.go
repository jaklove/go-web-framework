package framework

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//框架核心结构
type Core struct {
	router map[string]*Tree
}

//初始框架代码
func NewCore()*Core  {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"]  = NewTree()
	router["POST"] = NewTree()
	router["PUT"]  = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

//框架核心结构实现handler接口
func (c *Core)ServeHTTP(response http.ResponseWriter,request *http.Request)  {
	//todo
	log.Println("core.serveHTTP")
	ctx := NewContext(request,response)

	//  寻找路由
    router := c.FindRouteByRequest(request)
    if router == nil{
		// 如果没有找到，这里打印日志
		ctx.Json(404,"not found")
		return
	}

    // 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
    if err := router(ctx);err != nil{
    	ctx.Json(500,"inner error")
		return
	}
}


// 对应 Method = Get
func (c *Core)Get(url string,handler ControllerHandler)  {
	fmt.Println("注册get路由")
	if err := c.router["GET"].AddRouter(url,handler);err != nil{
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = POST
func (c *Core)Post(url string,handler ControllerHandler)  {
	if err := c.router["POST"].AddRouter(url,handler);err != nil{
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = PUT
func (c *Core)PUT(url string,handler ControllerHandler)  {
	if err := c.router["PUT"].AddRouter(url,handler);err != nil{
		log.Fatal("add router error: ", err)
	}
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handler ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url,handler);err != nil{
		log.Fatal("add router error: ", err)
	}
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core)FindRouteByRequest(request *http.Request)ControllerHandler{
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri         := request.URL.Path
	method      := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// 构建一个group组
func (c *Core)Group(prefix string)IGroup {
	return NewGroup(c,prefix)
}

