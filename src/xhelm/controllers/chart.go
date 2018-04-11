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
// @Description   获取示图列表
// @Success 201 {object} chart.ListResponse
// @Failure 500
// @router /charts [Get]
func (this *ChartController) ListCharts() {
	req := &chart.ListRequest{}
	this.setOutput(chart.List(req))
}

// Getchart
// @Title chart
// @Description   获取示图列表
// @Param chart path string true "vsphere集群 uuid"
// @Success 201 {object} chart.GetResponse
// @Failure 500
// @router /charts/:chart [Get]
func (this *ChartController) GetChart() {
	req := &chart.GetRequest{}
	chartName := this.Ctx.Input.Param(":chart")
	if strings.TrimSpace(chartName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild chart name"))
		return
	}
	req.Name = chartName
	this.setOutput(chart.Get(req))
}

// NewChart
// @Title chart
// @Description   添加新示图
// @Param body body chart.CreateRequest true "示图数据"
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
// @Description   删除新示图
// @Param body body chart.CreateRequest true "示图数据"
// @Success 201 {string} create success!
// @Failure 500
// @router /charts/:chart [Delete]
func (this *ChartController) DeleteChart() {
	req := &chart.DeleteRequest{}
	chartName := this.Ctx.Input.Param(":chart")
	if strings.TrimSpace(chartName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild chart name"))
		return
	}
	req.Name = chartName

	this.setOutput(chart.Delete(req))
}
