package db

import (
	"fmt"
	"xhelm/db/etcd"
)

var (
	CDB = NewChartDB()
)

//
const (
	etcdChartRootPath = "/xhelm/charts"
)

type ChartDB interface {
	ListAllChartsMetadata(repo string, metadatas interface{}) error
	RemoveChartVersion(repo string, chart string, version string) error
	RemoveChart(repo string, chart string) error
	GetChartVersionMetadata(repo string, chart string, version string, metadata interface{}) error
	UpdateChartVersionMetadata(repo string, chart string, version string, metadata interface{}) error
	UpdateChartVersionTemplate(repo string, chart string, version string, template interface{}) error
	GetChartAllVersionMetadata(repo string, chart string, metadatas interface{}) error
	GetChartVersionTemplate(repo string, chart string, version string, template interface{}) error
	CreateChart(repo string, chart string, version string, metadata interface{}, template interface{}) error
}

//考虑多方面原因, 应用商店并发操作量太小,出现同时读写的可能性太低,没必要使用事务.
//后续出现问题再改
//节点切换时,注意清理垃圾数据即可
/*
 |-metadata--\
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

func generateEtcdChartTemplateRepoKey(repo string) string {
	return etcdChartRootPath + "/" + "template" + "/" + repo
}

func generateEtcdChartTemplateKey(repo string, chart string) string {
	return etcdChartRootPath + "/" + "template" + "/" + repo + "/" + chart
}

func generateEtcdChartTemplateVersionKey(repo string, chart, version string) string {
	return etcdChartRootPath + "/" + "template" + "/" + repo + "/" + chart + "/" + version
}

type etcdChartDB struct {
	etcd *etcd.Etcd
}

func NewChartDB() ChartDB {
	return &etcdChartDB{etcd: etcd.GlobalEtcd}
}

func (cb *etcdChartDB) ListAllChartsMetadata(repo string, metadatas interface{}) error {
	metadatakey := generateEtcdChartMetadataRepoKey(repo)
	err := cb.etcd.ListUnmarshal(metadatakey, metadatas)
	return err
}

func (cb *etcdChartDB) RemoveChart(repo string, chart string) error {
	metadatakey := generateEtcdChartMetadataKey(repo, chart)
	templatekey := generateEtcdChartTemplateKey(repo, chart)
	err := cb.etcd.Delete(metadatakey)
	if err != nil {
		return err
	}
	err = cb.etcd.Delete(templatekey)
	if err != nil {
		return err
	}
	return nil
}

func (cb *etcdChartDB) RemoveChartVersion(repo string, chart string, version string) error {
	metadatakey := generateEtcdChartMetadataVersionKey(repo, chart, version)
	templatekey := generateEtcdChartMetadataVersionKey(repo, chart, version)

	err := cb.etcd.Delete(metadatakey)
	if err != nil {
		return err
	}
	err = cb.etcd.Delete(templatekey)
	return err
}

//这个要有能够通过前缀回去chart的功能
//返回chart的所有版本
func (cb *etcdChartDB) GetChartVersionMetadata(repo string, chart string, version string, metadata interface{}) error {
	metadatakey := generateEtcdChartMetadataVersionKey(repo, chart, version)
	err := cb.etcd.GetUnmarshal(metadatakey, metadata)
	return err
}

func (cb *etcdChartDB) GetChartVersionTemplate(repo string, chart string, version string, template interface{}) error {
	templatekey := generateEtcdChartTemplateVersionKey(repo, chart, version)
	err := cb.etcd.GetUnmarshal(templatekey, template)
	return err
}

func (cb *etcdChartDB) GetChartAllVersionMetadata(repo string, chart string, metadatas interface{}) error {
	metadatakey := generateEtcdChartTemplateKey(repo, chart)
	err := cb.etcd.ListUnmarshal(metadatakey, metadatas)
	return err
}

func (cb *etcdChartDB) UpdateChartVersionMetadata(repo string, chart string, version string, metadata interface{}) error {
	metadatakey := generateEtcdChartMetadataVersionKey(repo, chart, version)

	if !cb.etcd.IsExist(metadatakey) {
		return fmt.Errorf("chart:%v/%v.%v not exist", repo, chart, version)
	}
	err := cb.etcd.PutMarshal(metadatakey, metadata, 0)
	if err != nil {
		return err
	}
	return nil

}

func (cb *etcdChartDB) UpdateChartVersionTemplate(repo string, chart string, version string, template interface{}) error {
	templatekey := generateEtcdChartTemplateVersionKey(repo, chart, version)

	if !cb.etcd.IsExist(templatekey) {
		return fmt.Errorf("chart:%v/%v.%v not exist", repo, chart, version)
	}
	err := cb.etcd.PutMarshal(templatekey, template, 0)
	if err != nil {
		return err
	}
	return nil

}

func (cb *etcdChartDB) CreateChart(repo string, chart string, version string, metadata interface{}, template interface{}) error {
	metadatakey := generateEtcdChartMetadataVersionKey(repo, chart, version)
	templatekey := generateEtcdChartTemplateVersionKey(repo, chart, version)

	if cb.etcd.IsExist(metadatakey) {
		return fmt.Errorf("chart:%v/%v.%v has exist", repo, chart, version)
	}
	err := cb.etcd.PutMarshal(metadatakey, metadata, 0)
	if err != nil {
		return err
	}

	err = cb.etcd.PutMarshal(templatekey, template, 0)
	if err != nil {
		return err
	}

	return err
}
