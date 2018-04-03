package db

import "encoding/json"

var (
	RDB = NewRepoDb()
)

type RepoDb interface {
	CreateRepo(repo string, data interface{}) error
	DeleteRepo(repo string) error
	ListRepos() (map[string][]byte, error)
}

type etcdRepoDb struct {
}

func NewRepoDb() RepoDb {
	return &etcdRepoDb{}
}

func (erb *etcdRepoDb) DeleteRepo(repo string) error {
	return nil
}

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
