package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// 自定义 Context
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handle         ControllerHandler

	// 是否超时标记位
	writerMux *sync.Mutex
	// 是否超时标记位
	hasTimeout bool
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
	}
}


func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}


func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.hasTimeout {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	bytes, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(bytes)
	return nil
}
