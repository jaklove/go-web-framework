package route

import (
	"fmt"
	"go-web/controller"
	"go-web/framework"
)

func RegisterRouter(core *framework.Core)  {
	fmt.Println("注册路由")
	// 需求1+2:HTTP方法+静态路由匹配

	fmt.Println("controller:",controller.UsersLoginController)
	core.Get("/user/login",controller.UsersLoginController)
	//core.Get("foo",controller.FooControllerHandler)
	////
	//// 需求3:批量通用前缀
	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", controller.SubjectDelController)
		subjectApi.Put("/:id", controller.SubjectUpdateController)
		subjectApi.Get("/:id", controller.SubjectGetController)
		subjectApi.Get("/list/all", controller.SubjectListController)
	}
}