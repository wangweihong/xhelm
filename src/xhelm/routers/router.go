package routers

import (
	"xhelm/controllers"

	"github.com/astaxie/beego"
)

func init() {
	InitXhelmRouters()

}

// InitVephereRouters init vephere routers
func InitXhelmRouters() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/*",
			//Options 用于跨域复杂请求预检
			beego.NSRouter("/*", &controllers.RepoController{}, "options:Options"),
			beego.NSRouter("/*/*", &controllers.RepoController{}, "options:Options"),
			beego.NSRouter("/*/*/*", &controllers.RepoController{}, "options:Options"),
		),
		beego.NSNamespace("/xhelm",
			beego.NSInclude(
				&controllers.RepoController{},
				&controllers.ChartController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
