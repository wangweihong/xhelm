package repository

const (
	localRepoRootPath = "/var/local/xhelm"

	etcdRepoRootPath = "/xhelm/repo"
)

func LocalRepoPath(repo string) string {
	return localRepoRootPath + "/" + repo
}
