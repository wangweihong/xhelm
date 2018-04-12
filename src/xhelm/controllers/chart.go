package controllers

import (
	"fmt"
	"strings"
	"xhelm/models/chart"
)

type ChartController struct {
	baseController
}

// Listcharts
// @Title chart
// @Description   获取仓库所有应用模板
// @Param repository path string true "仓库"
// @Success 201 {object} chart.ListResponse
// @Failure 500
// @router /repositories/:repository/charts [Get]
func (this *ChartController) ListCharts() {
	req := &chart.ListRequest{}
	repoName := this.Ctx.Input.Param(":repository")
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repository name"))
		return
	}
	req.Repo = repoName
	this.setOutput(chart.List(req))
}

// Getchart
// @Title chart
// @Description   获取应用模板元数据
// @Param repository path string true "仓库"
// @Param chart path string true "应用模板"
// @Param version path string true "版本"
// @Success 201 {object} chart.GetResponse
// @Failure 500
// @router /repositories/:repository/charts/:chart/version/:version [Get]
func (this *ChartController) GetChart() {
	req := &chart.GetRequest{}
	repoName := this.Ctx.Input.Param(":repository")
	chartName := this.Ctx.Input.Param(":chart")
	version := this.Ctx.Input.Param(":version")
	if strings.TrimSpace(chartName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild chart name"))
		return
	}
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repo name"))
		return
	}
	if strings.TrimSpace(version) == "" {
		this.setOutput(nil, fmt.Errorf("invaild version"))
		return
	}
	req.Name = chartName
	req.Repo = repoName
	req.Version = version
	this.setOutput(chart.Get(req))
}

// NewChart
// @Title chart
// @Description   添加应用模板
// @Param body body chart.CreateRequest true "应用模板创建参数"
// @Success 201 {string} create success!
// @Failure 500
// @router /charts [Post]
func (this *ChartController) NewChart() {
	req := &chart.CreateRequest{}
	if err := this.unmarshalRequestBodyJSON(req); err != nil {
		this.setOutput(nil, err)
		return
	}

	this.setOutput(chart.New(req))
}

// DeleteChart
// @Title chart
// @Description   删除应用模板
// @Param repository path string true "仓库"
// @Param chart path string true "应用模板"
// @Success 201 {string} create success!
// @Failure 500
// @router /repoistories/:repository/charts/:chart [Delete]
func (this *ChartController) DeleteChart() {
	req := &chart.DeleteRequest{}
	chartName := this.Ctx.Input.Param(":chart")
	repoName := this.Ctx.Input.Param(":repository")
	if strings.TrimSpace(chartName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild chart name"))
		return
	}
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repo name"))
		return
	}

	req.Name = chartName
	req.Repo = repoName

	this.setOutput(chart.Delete(req))
}

// DeleteChart
// @Title chart
// @Description   删除应用模板
// @Param repository path string true "仓库"
// @Param chart path string true "应用模板"
// @Param version path string true "版本"
// @Success 201 {string} create success!
// @Failure 500
// @router /repoistories/:repository/charts/:chart/version/:version [Delete]
func (this *ChartController) DeleteChartVersion() {
	req := &chart.DeleteRequest{}
	chartName := this.Ctx.Input.Param(":chart")
	repoName := this.Ctx.Input.Param(":repository")
	version := this.Ctx.Input.Param(":version")

	if strings.TrimSpace(chartName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild chart name"))
		return
	}
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repo name"))
		return
	}
	if strings.TrimSpace(version) == "" {
		this.setOutput(nil, fmt.Errorf("invaild version name"))
		return
	}
	req.Name = chartName
	req.Repo = repoName
	req.Version = &version

	this.setOutput(chart.Delete(req))
}
