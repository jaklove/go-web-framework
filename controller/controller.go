package controller

import (
	"context"
	"fmt"
	"go-web/framework"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context)error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	fmt.Println("foo controller")

	durationCtx, cancelFunc := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancelFunc()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		time.Sleep(10 * time.Second)
		c.Json(200, "ok")

		//完成
		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}

func UserLoginController(c *framework.Context)error  {
	c.Json(200, "ok, UserLoginController")
	return nil
}