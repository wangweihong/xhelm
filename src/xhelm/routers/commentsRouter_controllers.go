package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "ListCharts",
			Router: `/charts`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "GetChart",
			Router: `/charts/:chart`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "NewChart",
			Router: `/charts`,
			AllowHTTPMethods: []string{"Post"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "DeleteChart",
			Router: `/charts/:chart`,
			AllowHTTPMethods: []string{"Delete"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:RepoController"] = append(beego.GlobalControllerRouter["xhelm/controllers:RepoController"],
		beego.ControllerComments{
			Method: "ListRepos",
			Router: `/repositories`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:RepoController"] = append(beego.GlobalControllerRouter["xhelm/controllers:RepoController"],
		beego.ControllerComments{
			Method: "GetRepo",
			Router: `/repositories/:repository`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:RepoController"] = append(beego.GlobalControllerRouter["xhelm/controllers:RepoController"],
		beego.ControllerComments{
			Method: "NewRepo",
			Router: `/repositories`,
			AllowHTTPMethods: []string{"Post"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:RepoController"] = append(beego.GlobalControllerRouter["xhelm/controllers:RepoController"],
		beego.ControllerComments{
			Method: "DeleteRepo",
			Router: `/repository/:repository`,
			AllowHTTPMethods: []string{"Delete"},
			Params: nil})

}
