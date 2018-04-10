package db

import (
	"db/etcd"
	"fmt"
)

var (
	CDB = NewChartDB()
)

//
const (
	etcdChartRootPath = "/xhelm/charts"
)

type ChartDB interface {
	CreateChartVersion(repo string, chart string, version string, obj interface{}) error
	UpdateChartVersion(repo string, chart string, version string, obj interface{}) error
	ListCharts(repo string, objs interface{}) error
	RemoveChartVersion(repo string, chart string, version string) error
	RemoveChart(repo string, chart string) error
	GetChartVersion(repo string, chart string, version string, obj interface{}) error
	GetChartAllVersion(repo string, chart string, objs interface{}) error
}

/*
	 |-metadata--\
	  \- versions
						\- digest
						\- v1
		\- name
		\- createTime
	 |-package/v1
*/
func generateEtcdChartMetadataRepoKey(repo string) string {
	return etcdChartRootPath + "/" + "metadata" + "/" + repo
}

func generateEtcdChartMetadataKey(repo string, chart string) string {
	return etcdChartRootPath + "/" + "metadata" + "/" + repo + "/" + chart
}
func generateEtcdChartMetadataVersionKey(repo string, chart, version string) string {
	return etcdChartRootPath + "/" + "metadata" + "/" + repo + "/" + chart + "/" + version
}

func generateEtcdChartPackageRepoKey(repo string) string {
	return etcdChartRootPath + "/" + "package" + "/" + repo
}

func generateEtcdChartPackageKey(repo string, chart string) string {
	return etcdChartRootPath + "/" + "package" + "/" + repo + "/" + chart
}

func generateEtcdChartPackageVersionKey(repo string, chart, version string) string {
	return etcdChartRootPath + "/" + "package" + "/" + repo + "/" + chart + "/" + version
}

type etcdChartDB struct {
	etcd *etcd.Etcd
}

func NewChartDB() ChartDB {
	return &etcdChartDB{etcd: etcd.GlobalEtcd}
}

func (cb *etcdChartDB) ListCharts(repo string, objs interface{}) error {
	key := generateEtcdChartPackageRepoKey(repo)
	err := cb.etcd.ListUnmarshal(key, objs)
	return err

}

//移除
//TODO:事务
func (cb *etcdChartDB) RemoveChart(repo string, chart string) error {
	key := generateEtcdChartPackageKey(repo, chart)
	err := cb.etcd.Delete(key)
	return err
}

func (cb *etcdChartDB) RemoveChartVersion(repo string, chart string, version string) error {
	key := generateEtcdChartPackageVersionKey(repo, chart, version)
	err := cb.etcd.Delete(key)
	return err
}

//这个要有能够通过前缀回去chart的功能
//返回chart的所有版本
func (cb *etcdChartDB) GetChartVersion(repo string, chart string, version string, obj interface{}) error {
	key := generateEtcdChartPackageVersionKey(repo, chart, version)
	err := cb.etcd.GetUnmarshal(key, obj)
	return err
}

func (cb *etcdChartDB) GetChartAllVersion(repo string, chart string, objs interface{}) error {
	key := generateEtcdChartPackageKey(repo, chart)
	err := cb.etcd.ListUnmarshal(key, objs)
	return err
}

func (cb *etcdChartDB) UpdateChartVersion(repo string, chart string, version string, obj interface{}) error {
	key := generateEtcdChartPackageVersionKey(repo, chart, version)
	if !cb.etcd.IsExist(key) {
		return fmt.Errorf("chart:%v/%v.%v not exist", repo, chart, version)
	}
	err := cb.etcd.PutMarshal(key, obj, 0)
	return err
}

func (cb *etcdChartDB) CreateChartVersion(repo string, chart string, version string, obj interface{}) error {
	key := generateEtcdChartPackageVersionKey(repo, chart, version)
	if cb.etcd.IsExist(key) {
		return fmt.Errorf("chart:%v/%v.%v has exist", repo, chart, version)
	}
	err := cb.etcd.PutMarshal(key, obj, 0)
	return err
}
