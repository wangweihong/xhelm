package repository

import (
	"fmt"
	"os"
	"sync"
	"time"
	"xhelm/db"
	"xhelm/setting"
	"xhelm/xlog"

	goflock "github.com/theckman/go-flock"
	helmchartutil "k8s.io/helm/pkg/chartutil"
	helmdowloader "k8s.io/helm/pkg/downloader"
	helmgetter "k8s.io/helm/pkg/getter"
	helmrepo "k8s.io/helm/pkg/repo"
)

var (
	RM                         = &RepositoryManager{}
	repoflocker                = goflock.NewFlock(setting.LocalRepoRootPath() + "/" + "repo.lock")
	errRemoteRepoNotSupport    = fmt.Errorf("remote repo don't support this action")
	errRemoteRepoNotSupportYet = fmt.Errorf("remote repo don't support now")
)

//TODO: 不需要在内存中维护缓存, 这个访问量实在太少
type RepositoryManager struct {
}

type Repository struct {
	locker sync.RWMutex
	helmrepo.Entry

	Remote     bool
	CreateTime time.Time
}

type CreateOption struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func (rm *RepositoryManager) GetRepo(name string) (*Repository, error) {
	var r Repository
	err := db.RDB.GetRepo(name, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

//TODO: file lock
func cleanRepoLocalDir(repo string) error {
	repoflocker.Lock()
	defer repoflocker.Unlock()

	rp := setting.LocalRepoPath(repo)
	err := os.RemoveAll(rp)
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func createRepoLocalDir(repo string) error {
	repoflocker.Lock()
	defer repoflocker.Unlock()
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

func downloadRemoteChart(repo *Repository, chartName string, version *string) error {
	dest := setting.LocalRepoChartsRootPath(repo.Name)
	getters := helmgetter.Providers{
		{
			Schemes: []string{"http", "https"},
			New:     newHTTPGetter,
		},
	}

	c := helmdowloader.ChartDownloader{
		Getters:  getters,
		Username: repo.Username,
		Password: repo.Password,
	}

	chartVersion := ""
	if version != nil {
		chartVersion = *version
	}

	chartURL, err := helmrepo.FindChartInAuthRepoURL(repo.URL, repo.Username, repo.Password, chartName, chartVersion, repo.CertFile, repo.KeyFile, repo.CAFile, getters)
	if err != nil {
		return nil
	}

	saved, v, err := c.DownloadTo(chartURL, chartVersion, dest)
	if err != nil {
		return err
	}

	xlog.Logger.Info("Verification: %v\n", v)
	xlog.Logger.Info("Chart Download to :%v", saved)

	ud := setting.LocalRepoCacheRootPath(repo.Name)
	err = helmchartutil.ExpandFile(ud, saved)
	if err != nil {
		return err
	}

	return nil
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

func (rm *RepositoryManager) NewRepo(name string, opt *CreateOption) error {
	if err := isNameValid(name); err != nil {
		return err
	}

	var repo Repository
	repo.CreateTime = time.Now()
	repo.Name = name
	if opt != nil {
		return errRemoteRepoNotSupportYet
		repo.Remote = true
		repo.CAFile = opt.CertFile
		repo.KeyFile = opt.KeyFile
		repo.URL = opt.URL
		repo.Username = opt.Username
		repo.Password = opt.Password
	}

	//如果因为已存在, 则直接创建失败
	if err := db.RDB.CreateRepo(name, repo); err != nil {
		return fmt.Errorf(" create repo in db failed: %v", err)
	}

	var e error
	var err2 error

	//创建本地目录
	//TODO: 处理可能会残留的旧的repo的目录
	err := createRepoLocalDir(name)
	if err != nil {
		e = err
		goto clean_etcd
	}

	//TODO:清理
	if repo.Remote {
		err = downloadRemoteRepoIndex(&repo)
		if err != nil {
			e = err
			goto clean_local
		}
	}

	return nil

clean_local:
	err2 = cleanRepoLocalDir(name)
	if err2 != nil {
		xlog.Logger.Errorf("clean repo local '%v' dir fail while creating: %v", setting.LocalRepoPath(name), err2)
	}

clean_etcd:
	err2 = db.RDB.DeleteRepo(name)
	if err2 != nil {
		xlog.Logger.Error(err2)
	}

	return e
}

func isRepoLocal(repo *Repository) bool {
	if repo.URL == "" {
		return true
	}

	return false
}

//TODO:通知chart进行清理
func (rm *RepositoryManager) DeleteRepo(name string) error {
	//优先删除etcd
	err := db.RDB.DeleteRepo(name)
	if err != nil {
		//TODO: 检测etcd不存在的报错
		return err
	}

	err2 := cleanRepoLocalDir(name)
	if err2 != nil {
		xlog.Logger.Errorf("clean repo local '%v' dir fail while deleting: %v", setting.LocalRepoPath(name), err2)
	}
	return nil
}

func (rm *RepositoryManager) ListRepos() ([]Repository, error) {
	repos := make([]Repository, 0)
	err := db.RDB.ListRepos(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

//TODO:
//1. 在切换节点时, 从数据库中获取所有local repo charts的信息,
//2. 这些信息组成构建index.yaml文件.
//3. 在添加/删除chart,以及chart版本时, 对数据库中的文件进行清理,以及更新index.yaml文件
//4. 不要长期缓存index.yaml文件到内存中, 太占内存, 用完立即释放
func LoadLocalRepo() error {
	//添加生成本地的目录
	//加载所有repo信息
	//db.RDB.ListRepos()
	//

	return nil
}
