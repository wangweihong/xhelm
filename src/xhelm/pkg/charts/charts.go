package charts

import (
	"xhelm/repository"
)

const (
	localChartRelativePath = "charts"
	localCacheRelativePath = "cache"
)

func LocalChartPath(repo, chart string) string {
	return repository.LocalRepoPath(repo) + "/" + chart

}
