package repository

import (
	"db"
	"fmt"
	"os"
	"setting"
	"sync"
	"time"
	"xlog"

	helmgetter "k8s.io/helm/pkg/getter"
	helmrepo "k8s.io/helm/pkg/repo"
)

const (
	//localRepoRootPath = "/var/local/xhelm"

	etcdRepoRootPath = "/xhelm/repo"

	StateInitilizing  = "initializing"  //用于表明repo仍然在创建过程中,不接受chart创建
	StateInitComplete = "initcompleted" //repo完成创建过程,接收正常操作
)

var (
	RM = &RepositoryManager{Repos: make(map[string]Repository)}
)

type RepositoryManager struct {
	Locker sync.RWMutex
	Repos  map[string]Repository
}

type Repository struct {
	helmrepo.Entry

	Remote     bool
	State      string //不保留到etcd中,每次master节点更改,都需要重新初始化
	CreateTime int64
}

type CreateOption struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func (rm *RepositoryManager) getRepo(Name string) (*Repository, error) {
	repo, ok := rm.Repos[Name]
	if !ok {
		return nil, fmt.Errorf("repo not found")
	}
	return &repo, nil
}

//TODO: file lock
func cleanRepoLocalDir(repo string) error {
	rp := setting.LocalRepoPath(repo)
	err := os.RemoveAll(rp)
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func createRepoLocalDir(repo string) error {
	//
	//强制删除以前可能残留的目录架构
	err := cleanRepoLocalDir(repo)
	if err != nil {
		return err
	}

	rp := setting.LocalRepoPath(repo)
	err = os.MkdirAll(rp, 0755)
	if err != nil {
		return err
	}

	chartDir := setting.LocalRepoChartsRootPath(repo)
	cacheDir := setting.LocalRepoCacheRootPath(repo)

	err = os.Mkdir(chartDir, 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir(cacheDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func isNameValid(name string) error {
	if len(name) == 0 || len(name) > 50 {
		return fmt.Errorf("The name is between 0 and 50 characters.")
	}
	return nil
}

// newHTTPGetter constructs a valid http/https client as Getter
func newHTTPGetter(URL, CertFile, KeyFile, CAFile string) (helmgetter.Getter, error) {
	return helmgetter.NewHTTPGetter(URL, CertFile, KeyFile, CAFile)
}

func downloadRemoteRepoIndex(repo *Repository) error {
	c := repo.Entry
	getters := helmgetter.Providers{
		{
			Schemes: []string{"http", "https"},
			New:     newHTTPGetter,
		},
	}

	r, err := helmrepo.NewChartRepository(&c, getters)
	if err != nil {
		return err
	}

	cachePath := setting.LocalRepoIndexFile(repo.Name)
	if err := r.DownloadIndexFile(cachePath); err != nil {
		return fmt.Errorf("Looks like %q is not a valid chart repository or cannot be reached: %s", repo.URL, err.Error())
	}
	return nil
}

func (rm *RepositoryManager) AddRemoteRepo(name string, opt CreateOption) error {
	if err := isNameValid(name); err != nil {
		return err
	}

	rm.Locker.Lock()

	if _, err := rm.getRepo(name); err == nil {
		rm.Locker.Unlock()
		return fmt.Errorf("repo has exist")
	}

	var repo Repository
	repo.CreateTime = time.Now().Unix()
	repo.Name = name
	repo.Remote = true
	repo.CAFile = opt.CertFile
	repo.KeyFile = opt.KeyFile
	repo.URL = opt.URL
	repo.Username = opt.Username
	repo.Password = opt.Password

	repo.State = StateInitilizing
	rm.Repos[name] = repo

	rm.Locker.Unlock()

	//创建本地目录
	//TODO: 处理可能会残留的旧的repo的目录
	err := createRepoLocalDir(name)
	if err != nil {
		rm.Locker.Lock()
		defer rm.Locker.Unlock()
		delete(rm.Repos, name)
		return err
	}

	err = downloadRemoteRepoIndex(&repo)
	if err != nil {
		rm.Locker.Lock()
		defer rm.Locker.Unlock()
		delete(rm.Repos, name)
		return err
	}

	//避免创建文件阻塞, etcd网络阻塞, 导致锁长期被占用
	err = db.RDB.CreateRepo(name, repo)
	if err != nil {
		err2 := cleanRepoLocalDir(name)
		if err2 != nil {
			xlog.Logger.Errorf("clean repo local '%v' dir fail while creating: %v", setting.LocalRepoPath(name), err2)
		}
		rm.Locker.Lock()
		defer rm.Locker.Unlock()
		delete(rm.Repos, name)
		return err
	}

	rm.Locker.Lock()
	defer rm.Locker.Unlock()

	repo.State = StateInitComplete
	rm.Repos[name] = repo
	return nil

}

func (rm *RepositoryManager) CreateRepo(name string) error {
	if err := isNameValid(name); err != nil {
		return err
	}

	rm.Locker.Lock()

	if _, err := rm.getRepo(name); err == nil {
		rm.Locker.Unlock()
		return fmt.Errorf("repo has exist")
	}

	var repo Repository
	repo.CreateTime = time.Now().Unix()
	repo.Name = name
	repo.State = StateInitilizing
	rm.Repos[name] = repo

	rm.Locker.Unlock()

	//创建本地目录
	//TODO: 处理可能会残留的旧的repo的目录
	err := createRepoLocalDir(name)
	if err != nil {
		rm.Locker.Lock()
		defer rm.Locker.Unlock()
		delete(rm.Repos, name)
		return err
	}

	//避免创建文件阻塞, etcd网络阻塞, 导致锁长期被占用
	err = db.RDB.CreateRepo(name, repo)
	if err != nil {
		err2 := cleanRepoLocalDir(name)
		if err2 != nil {
			xlog.Logger.Errorf("clean repo local '%v' dir fail while creating: %v", setting.LocalRepoPath(name), err2)
		}
		rm.Locker.Lock()
		defer rm.Locker.Unlock()
		delete(rm.Repos, name)
		return err
	}

	rm.Locker.Lock()
	defer rm.Locker.Unlock()

	repo.State = StateInitComplete
	rm.Repos[name] = repo
	return nil
}

func isRepoLocal(repo *Repository) bool {
	if repo.URL == "" {
		return true
	}

	return false
}

//TODO:通知chart进行清理
func (rm *RepositoryManager) DeleteRepo(name string) error {
	rm.Locker.Lock()

	_, err := rm.getRepo(name)
	if err != nil {
		rm.Locker.Unlock()
		return fmt.Errorf("repo not found")
	}
	rm.Locker.Unlock()

	//优先删除etcd
	err = db.RDB.DeleteRepo(name)
	if err != nil {
		return err
	}

	rm.Locker.Lock()
	delete(rm.Repos, name)
	rm.Locker.Unlock()

	err2 := cleanRepoLocalDir(name)
	if err2 != nil {
		xlog.Logger.Errorf("clean repo local '%v' dir fail while deleting: %v", setting.LocalRepoPath(name), err2)
	}
	return nil
}

func (rm *RepositoryManager) ListRepos() []Repository {
	rm.Locker.RLock()
	defer rm.Locker.RUnlock()

	repos := make([]Repository, 0)
	for _, v := range rm.Repos {
		repos = append(repos, v)
	}

	return repos
}

/*
func Init() error {
	err := os.MkdirAll( 0755)
	if err != nil {
		return err
	}
	return nil
}
*/
