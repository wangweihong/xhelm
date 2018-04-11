package setting

const (
	//TODO: test
	localRootPath = "/home/wwh/xhelm"

	localRepoRootPath   = localRootPath + "/" + "repository"
	localPluginRootPath = localRootPath + "/" + "plugins"
)

//获取repo本地路径
func LocalRepoRootPath() string {
	return localRepoRootPath
}

func LocalRepoPath(repo string) string {
	return localRepoRootPath + "/" + repo
}

func LocalPluginPath() string {
	return localPluginRootPath
}

func LocalRepoChartsRootPath(repo string) string {
	return LocalRepoPath(repo) + "/charts"
}

func LocalRepoCacheRootPath(repo string) string {
	return LocalRepoPath(repo) + "/cache"
}

func LocalRepoIndexFile(repo string) string {
	return LocalRepoPath(repo) + "/index.yaml"
}
