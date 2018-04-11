package chart

import (
	"xhelm/charts"
)

type Chart struct {
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
}

func transfrom(v *charts.Chart) Chart {
	var r Chart
	return r
}

func List(req *ListRequest) (*ListResponse, error) {
	var resp ListResponse
	return &resp, nil
}

func Get(req *GetRequest) (*GetResponse, error) {
	var resp GetResponse
	return &resp, nil
}

func New(req *CreateRequest) (*CreateResponse, error) {
	var resp CreateResponse
	return &resp, nil
}

func Delete(req *DeleteRequest) (*DeleteResponse, error) {
	var resp DeleteResponse
	return &resp, nil
}
