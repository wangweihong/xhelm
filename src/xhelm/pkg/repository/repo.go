package repository

const (
	localRepoRootPath = "/var/local/xhelm"

	etcdRepoRootPath = "/xhelm/repo"
)

//获取repo本地路径
func LocalRepoPath(repo string) string {
	return localRepoRootPath + "/" + repo
}
