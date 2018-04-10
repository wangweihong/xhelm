package db

import (
	"db/etcd"
	"fmt"
)

var (
	RDB = NewRepoDb()
)

const (
	etcdRepoRootPath = "/xhelm/repos"
)

type RepoDb interface {
	CreateRepo(repo string, data interface{}) error
	DeleteRepo(repo string) error
	ListRepos(objs interface{}) error
	GetRepo(repo string, obj interface{}) error
	UpdateRepo(repo string, data interface{}) error
}

type etcdRepoDb struct {
	etcd *etcd.Etcd
}

func NewRepoDb() RepoDb {
	return &etcdRepoDb{etcd: etcd.GlobalEtcd}
}

func generateEtcdRepoKey(repo string) string {
	return etcdRepoRootPath + "/" + repo
}

//使用事务来移除chart?
func (erb *etcdRepoDb) DeleteRepo(repo string) error {
	key := generateEtcdRepoKey(repo)
	return erb.etcd.Delete(key)
}

//如果已经存在则报错.
//事务
func (erb *etcdRepoDb) CreateRepo(repo string, data interface{}) error {
	key := generateEtcdRepoKey(repo)
	if erb.etcd.IsExist(key) {
		return fmt.Errorf("repo has exist")
	}

	err := erb.etcd.PutMarshal(key, data, 0)
	return err
}

func (erb *etcdRepoDb) ListRepos(values interface{}) error {
	key := etcdRepoRootPath + "/"
	err := erb.etcd.ListUnmarshal(key, values)
	return err
}

func (erb *etcdRepoDb) GetRepo(repo string, obj interface{}) error {
	key := generateEtcdRepoKey(repo)
	err := erb.etcd.GetUnmarshal(key, obj)
	return err
}

func (erb *etcdRepoDb) UpdateRepo(repo string, data interface{}) error {
	key := generateEtcdRepoKey(repo)
	if !erb.etcd.IsExist(key) {
		return fmt.Errorf("repo not exist")
	}
	err := erb.etcd.PutMarshal(key, data, 0)
	return err
}
