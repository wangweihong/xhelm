package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "ListCharts",
			Router: `/repositories/:repository/charts`,
			AllowHTTPMethods: []string{"Get"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "GetChart",
			Router: `/repositories/:repository/charts/:chart/version/:version`,
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
			Router: `/repoistories/:repository/charts/:chart`,
			AllowHTTPMethods: []string{"Delete"},
			Params: nil})

	beego.GlobalControllerRouter["xhelm/controllers:ChartController"] = append(beego.GlobalControllerRouter["xhelm/controllers:ChartController"],
		beego.ControllerComments{
			Method: "DeleteChartVersion",
			Router: `/repoistories/:repository/charts/:chart/version/:version`,
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
			Router: `/repositories/:repository`,
			AllowHTTPMethods: []string{"Delete"},
			Params: nil})

}
