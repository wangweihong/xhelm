package db

import (
	"db/etcd"
)

var (
	CDB = NewChartDB()
)

//
const (
	etcdChartRootPath = "/xhelm/charts"
)

type ChartDB interface {
	ListCharts(repo string) (map[string][]byte, error)
	RemoveChartVersion(repo string, chart string, version string) error
	RemoveChart(repo string, chart string) error
	GetChartVersion(repo string, chart string, version string) ([]byte, error)
	GetChartAllVersion(repo string, chart string) (map[string][]byte, error)
}

func generateEtcdChartKey(repo string, chart string) string {
	return etcdChartRootPath + "/" + repo + "_" + chart
}

type etcdChartDB struct {
	etcd *etcd.Etcd
}

func NewChartDB() ChartDB {
	return &etcdChartDB{etcd: etcd.GlobalEtcd}
}

func (cb *etcdChartDB) ListCharts(repo string) (map[string][]byte, error) {
	return nil, nil
}

//移除
func (cb *etcdChartDB) RemoveChart(repo string, chart string) error {
	return nil
}

func (cb *etcdChartDB) RemoveChartVersion(repo string, chart string, version string) error {
	return nil
}

//这个要有能够通过前缀回去chart的功能
//返回chart的所有版本
func (cb *etcdChartDB) GetChartVersion(repo string, chart string, version string) ([]byte, error) {
	return nil, nil
}

func (cb *etcdChartDB) GetChartAllVersion(repo string, chart string) (map[string][]byte, error) {
	return nil, nil
}
