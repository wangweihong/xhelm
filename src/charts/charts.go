package charts

import (
	helmrepo "k8s.io/helm/pkg/repo"
)

type Index struct {
}
type Chart struct {
	Name     string `json:"name"`
	Versions helmrepo.ChartVersions
}
