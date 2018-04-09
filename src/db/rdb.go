package db

import (
	"db/etcd"
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
	ListRepos() (map[string][]byte, error)
	GetRepo(repo string) ([]byte, error)
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
	err := erb.etcd.PutMarshal(key, data, 0)
	return err
}

func (erb *etcdRepoDb) ListRepos() (map[string][]byte, error) {
	key := etcdRepoRootPath
	erb.etcd.ListUnmarshal()
}

func (erb *etcdRepoDb) GetRepo(repo string) ([]byte, error) {
	key := generateEtcdRepoKey(repo)
	return nil, nil
}

func (erb *etcdRepoDb) UpdateRepo(repo string, data interface{}) error {
	key := generateEtcdRepoKey(repo)
	erb.etcd.PutMarshal(key, data, 0)
	return nil
}
