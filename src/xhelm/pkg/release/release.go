package release

import (
	"xhelm/util/request"
)

//这里需要吗???
type Generator interface {
	Generate(yaml string) error //用于将解析后的yaml进行部署
}

type xfleetGenerator struct {
	url  string
	name string
}

type XfleetCreateApplicationRequest struct {
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	App       string `json:"name"`
	User      string `json:"user"`
	Comment   string `json:"comment"`
	Data      string `json:"data"`
}

func NewGenerator(endpoint string) (Generator, error) {

	gtor := &xfleetGenerator{}
	gtor.url = endpoint
	gtor.name = "xfleet"

	err := gtor.Healthy()
	if err != nil {
		return nli, err.Errorf("release generator '%v' is unhealthy :%v", gtor.name, err)
	}
	return gtor, nil
}

func (x *xfleetGenerator) Healthy() error {
	url := x.url + "/test"
	_, err := request.Get(url, "")
	if err != nil {
		return err
	}
	return nil
}

func (x *xfleetGenerator) Generate(yaml string) error {
	return nil
}
