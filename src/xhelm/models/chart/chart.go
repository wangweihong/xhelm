package chart

import (
	"time"
	"xhelm/charts"
	"xhelm/repository"
	"xhelm/xlog"
)

type Chart struct {
	Name       string   `json:"name"`
	CreateTime int64    `json:"create_time"` //最新版本的创建时间
	Versions   []string `json:"versions"`    //版本列表
	Latest     string   `json:"latest"`      //最新版本
	//Icon       string   `json:"icon"`   //最新版本图标
}

type Detail struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AppVersion  string `json:"appVersion"`
	Template    string `json:"template"`
	Values      string `json:"string"`
}

func addToChartsIfNotExist(cs *[]Chart, t *charts.Metadata) {
	for k, v := range *cs {
		if v.Name == t.Name {
			for _, j := range v.Versions {
				if j == t.Version {
					//已存在,什么都不做
					return
				}
			}
			//chart相同,但版本不存在,则添加
			v.Versions = append(v.Versions, t.Version)
			//根据时间获取最新的版本
			if time.Unix(v.CreateTime, 0).Before(t.CreateTime) {
				v.Latest = t.Version
				v.CreateTime = t.CreateTime.Unix()
			}
			(*cs)[k] = v
			return
		}
	}

	var c Chart
	c.Name = t.Name
	//	c.Icon = t.Icon
	c.CreateTime = t.CreateTime.Unix()
	c.Latest = t.Version
	c.Versions = make([]string, 0)
	c.Versions = append(c.Versions, t.Version)

	*cs = append(*cs, c)
	return
}

func List(req *ListRequest) (*ListResponse, error) {
	var resp ListResponse

	metas, err := repository.RM.ListCharts(req.Repo)
	if err != nil {
		return nil, err
	}

	resp.Charts = make([]Chart, 0)
	for _, v := range metas {
		addToChartsIfNotExist(&resp.Charts, &v)
	}
	resp.Total = len(resp.Charts)

	return &resp, nil
}

func Get(req *GetRequest) (*GetResponse, error) {
	var resp GetResponse

	c, err := repository.RM.GetChartVersion(req.Repo, req.Name, req.Version)
	if err != nil {
		return nil, err
	}
	xlog.Logger.Info()
	var d Detail
	d.Name = c.Metadata.Name
	d.AppVersion = c.Metadata.AppVersion
	d.Description = c.Metadata.Description
	var t string
	for _, v := range c.Templates {
		t = t + string(v.Data)
		t = t + "\n---\n"
	}
	d.Template = t
	d.Values = c.Values.String()
	resp.Detail = d

	return &resp, nil
}

func New(req *CreateRequest) (*CreateResponse, error) {
	var resp CreateResponse
	opt := repository.ChartCreateOption{}
	opt.DefaultValues = make([]byte, len(req.DefaultValues))
	opt.Template = make([]byte, len(req.Template))
	opt.Version = req.Version
	opt.Description = req.Description
	copy(opt.DefaultValues, req.DefaultValues)
	copy(opt.Template, req.Template)

	err := repository.RM.CreateChart(req.Repo, req.Name, opt)
	return &resp, err
}

func Delete(req *DeleteRequest) (*DeleteResponse, error) {
	var resp DeleteResponse
	err := repository.RM.DeleteChart(req.Repo, req.Name, req.Version)
	return &resp, err
}
