package repository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"xhelm/charts"
	"xhelm/db"
	"xhelm/setting"
	"xhelm/xlog"

	helmchartutil "k8s.io/helm/pkg/chartutil"
	helmchart "k8s.io/helm/pkg/proto/hapi/chart"
)

const defaultIgnore = `# Patterns to ignore when building packages.
# This supports shell glob matching, relative path matching, and
# negation (prefixed with !). Only one pattern per line.
.DS_Store
# Common VCS dirs
.git/
.gitignore
.bzr/
.bzrignore
.hg/
.hgignore
.svn/
# Common backup files
*.swp
*.bak
*.tmp
*~
# Various IDEs
.project
.idea/
*.tmproj
`

const defaultHelpers = `{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "<CHARTNAME>.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "<CHARTNAME>.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "<CHARTNAME>.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
`
const (
	templateName    = "template.yaml"
	notesFileSuffix = "NOTES.txt"
	GoTplEngine     = "gotpl"
)

type ChartCreateOption struct {
	Version       string
	Description   string
	Template      []byte
	DefaultValues []byte //默认配置
}

func (rm *RepositoryManager) CreateChart(repoName string, chartName string, opt ChartCreateOption) error {
	repo, err := rm.GetRepo(repoName)
	if err != nil {
		return fmt.Errorf("repo not found")
	}

	if repo.Remote {
		return errRemoteRepoNotSupport
	} else {

		chartfile := &helmchart.Metadata{
			Name:        chartName,
			Version:     opt.Version,
			Description: opt.Description,
			Engine:      GoTplEngine,
			//		KubeVersion: opt.KubeVersion,
		}

		tmpDir, err := ioutil.TempDir("/tmp", "chart")
		if err != nil {
			return fmt.Errorf("create chart fail: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		cdir := filepath.Join(tmpDir, chartName)
		if fi, err := os.Stat(cdir); err == nil && !fi.IsDir() {
			return fmt.Errorf("create chart fail: file %s already exists and is not a directory", cdir)
		}

		if err := os.MkdirAll(cdir, 0755); err != nil {
			return fmt.Errorf("create chart fail: %v", err)
		}
		cf := filepath.Join(cdir, helmchartutil.ChartfileName)
		if _, err := os.Stat(cf); err != nil {
			if err := helmchartutil.SaveChartfile(cf, chartfile); err != nil {
				return fmt.Errorf("create chart fail: %v", err)
			}
		}

		for _, d := range []string{helmchartutil.TemplatesDir, helmchartutil.ChartsDir} {
			if err := os.MkdirAll(filepath.Join(cdir, d), 0755); err != nil {
				return err
			}
		}

		files := []struct {
			path    string
			content []byte
		}{
			{
				path:    filepath.Join(cdir, helmchartutil.ValuesfileName),
				content: opt.DefaultValues,
			}, {
				path:    filepath.Join(cdir, helmchartutil.IgnorefileName),
				content: []byte(defaultIgnore),
			}, {
				path:    filepath.Join(cdir, helmchartutil.TemplatesDir, templateName),
				content: opt.Template,
			}, {
				path:    filepath.Join(cdir, helmchartutil.TemplatesDir, helmchartutil.HelpersName),
				content: helmchartutil.Transform(defaultHelpers, "<CHARTNAME>", chartfile.Name),
			},
		}

		for _, file := range files {
			//TODO:修正该逻辑
			if _, err := os.Stat(file.path); err == nil {
				// File exists and is okay. Skip it.
				continue
			}
			if err := ioutil.WriteFile(file.path, file.content, 0644); err != nil {
				return fmt.Errorf("create chart fail when writing to file content:%v", err)
			}
		}
		//压缩, 然后上传到etcd中

		ch, err := helmchartutil.LoadDir(cdir)
		if err != nil {
			return fmt.Errorf("create chart fail when try to load dir:%v", err)
		}

		if filepath.Base(cdir) != ch.Metadata.Name {
			return fmt.Errorf("directory name (%s) and Chart.yaml name (%s) must match", filepath.Base(cdir), ch.Metadata.Name)
		}

		dest := setting.LocalRepoChartsRootPath(repoName)
		name, err := helmchartutil.Save(ch, dest)
		if err != nil {
			return fmt.Errorf("create chart fail when try to save chart: %v", err)
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			return fmt.Errorf("create chart fail when ty to load chart: %v", err)
		}

		var metadata charts.Metadata
		metadata.Name = chartName
		metadata.Version = opt.Version
		metadata.CreateTime = time.Now()
		CompressedData := data
		//
		err = db.CDB.CreateChart(repoName, chartName, opt.Version, &metadata, &CompressedData)
		if err != nil {
			return fmt.Errorf("create chart fail when save chart to db: %v", err)
		}

	}
	return nil
}

func (rm *RepositoryManager) DeleteChart(repoName, chartName string, version *string) error {
	repo, err := rm.GetRepo(repoName)
	if err != nil {
		return fmt.Errorf("repo not found")
	}

	if repo.Remote {
		return errRemoteRepoNotSupport
	} else {
		if version != nil {
			return db.CDB.RemoveChartVersion(repoName, chartName, *version)
		} else {
			return db.CDB.RemoveChart(repoName, chartName)
		}
	}
	return nil
}

func (rm *RepositoryManager) GetChartVersionDetail(repoName, chartName string, version string) (*int, *string, error) {
	repo, err := rm.GetRepo(repoName)
	if err != nil {
		return nil, nil, fmt.Errorf("repo not found")
	}

	if repo.Remote {
		return nil, nil, errRemoteRepoNotSupport
	} else {
		metadata := &charts.Metadata{}
		compressedData := []byte{}

		if err := db.CDB.GetChartVersionMetadata(repoName, chartName, version, metadata); err != nil {
			return nil, nil, err
		}
		if err := db.CDB.GetChartVersionTemplate(repoName, chartName, version, &compressedData); err != nil {
			return nil, nil, err
		}

	}
	return nil, nil, fmt.Errorf("wrong logic")
}

func (rm *RepositoryManager) GetChart(repoName string, chartName string) {
}

//TODO: file lock
func (rm *RepositoryManager) ListCharts(repoName string) ([]charts.Metadata, error) {

	repo, err := rm.GetRepo(repoName)
	if err != nil {
		return nil, fmt.Errorf("repo not found")
	}

	mds := make([]charts.Metadata, 0)
	if repo.Remote {
		return nil, errRemoteRepoNotSupportYet
		/*
			indexFile := setting.LocalRepoIndexFile(repoName)
			indf, err := helmrepo.LoadIndexFile(indexFile)
			if err != nil {
				return nil, err
			}
			for k, v := range indf.Entries {
				var c charts.Chart
				c.Name = k
				c.Versions = append(c.Versions, v...)

				cs = append(cs, c)
			}
		*/
	} else {

		err := db.CDB.ListAllChartsMetadata(repoName, &mds)
		if err != nil {
			return nil, err
		}
	}
	return mds, nil
}

//不指定version,则拉取最新的版本
func (rm *RepositoryManager) GetChartVersion(repoName string, chartName string, version string) (*charts.Chart, error) {
	repo, err := rm.GetRepo(repoName)
	if err != nil {
		return nil, fmt.Errorf("repo not found")
	}

	if repo.Remote {
		//指定文件
		return nil, errRemoteRepoNotSupportYet
		err := downloadRemoteChart(repo, chartName, &version)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		/*
			var chart charts.Chart
			err := db.CDB.GetChartVersionMetadata(repoName, chartName, version, &chart.Metadata)
			if err != nil {
				return nil, err
			}
		*/
		compressedData := []byte{}
		err = db.CDB.GetChartVersionTemplate(repoName, chartName, version, &compressedData)
		if err != nil {
			return nil, err
		}

		var chart *charts.Chart
		err := rm.uncompressData(repoName, chartName, compressedData, &chart)
		if err != nil {
			return nil, err
		}
		return chart, nil
	}
	return nil, fmt.Errorf("wrong logic")
}

func (rm *RepositoryManager) uncompressData(repoName string, chartName string, compressedData []byte, chart **charts.Chart) error {
	ud := setting.LocalRepoCacheRootPath(repoName)
	if fi, err := os.Stat(ud); err != nil {
		if err := os.MkdirAll(ud, 0755); err != nil {
			return fmt.Errorf("Failed to untar (mkdir): %s", err)
		}

	} else if !fi.IsDir() {
		return fmt.Errorf("Failed to untar: %s is not a directory", ud)
	}

	r := bytes.NewReader(compressedData)
	err := helmchartutil.Expand(ud, r)
	if err != nil {
		return err
	}

	xlog.Logger.Info("expand path:", ud)
	cpbc, err := helmchartutil.Load(ud + "/" + chartName)
	if err != nil {
		return err
	}
	//	chart = (*charts.Chart)cpbc
	*chart = (*charts.Chart)(cpbc)
	return nil

} //最新版本
