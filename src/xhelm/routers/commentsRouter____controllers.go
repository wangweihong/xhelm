package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["controllers:RepoController"] = append(beego.GlobalControllerRouter["controllers:RepoController"],
		beego.ControllerComments{
			Method: "ListRepos",
			Router: `/repositories`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["controllers:RepoController"] = append(beego.GlobalControllerRouter["controllers:RepoController"],
		beego.ControllerComments{
			Method: "GetRepo",
			Router: `/repositories/:repository`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["controllers:RepoController"] = append(beego.GlobalControllerRouter["controllers:RepoController"],
		beego.ControllerComments{
			Method: "NewRepo",
			Router: `/repositories`,
			AllowHTTPMethods: []string{"Post"},
			Params: nil})

	beego.GlobalControllerRouter["controllers:RepoController"] = append(beego.GlobalControllerRouter["controllers:RepoController"],
		beego.ControllerComments{
			Method: "DeleteRepo",
			Router: `/repository/:repository`,
			AllowHTTPMethods: []string{"Delete"},
			Params: nil})

}
