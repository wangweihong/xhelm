package db

import "encoding/json"

var (
	RDB = NewRepoDb()
)

const (
	etcdRepoRootPath = "/xhelm/repo"
)

type RepoDb interface {
	CreateRepo(repo string, data interface{}) error
	DeleteRepo(repo string) error
	ListRepos() (map[string][]byte, error)
	GetRepo(repo string) ([]byte, error)
	UpdateRepo(repo string, data interface{}) error
}

type etcdRepoDb struct {
}

func NewRepoDb() RepoDb {
	return &etcdRepoDb{}
}

func (erb *etcdRepoDb) DeleteRepo(repo string) error {
	return nil
}

//如果已经存在则报错.
//事务
func (erb *etcdRepoDb) CreateRepo(repo string, data interface{}) error {
	_, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return nil

}

func (erb *etcdRepoDb) ListRepos() (map[string][]byte, error) {
	return nil, nil
}

func (erb *etcdRepoDb) GetRepo(repo string) ([]byte, error) {
	return nil, nil
}

func (erb *etcdRepoDb) UpdateRepo(repo string, data interface{}) error {
	return nil
}
