package controller

import "go-web/framework"

func UsersLoginController(c *framework.Context)error  {
	c.Json(200,"ok,UserLoginController")
	return nil
}