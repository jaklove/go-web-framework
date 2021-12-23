package main

import (
	"fmt"
	"go-web/framework"
	"go-web/route"
	"net/http"
)


func main()  {
	core := framework.NewCore()
	route.RegisterRouter(core)
	server := &http.Server{
		Handler:core,
		Addr:":8080",
	}

	fmt.Println("http server success!")
	//设置路由
	server.ListenAndServe()
}


