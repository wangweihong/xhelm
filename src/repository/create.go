package repository

import (
	"charts"
	"db"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"setting"

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
	repo, err := rm.getRepo(repoName)
	if err != nil {
		return fmt.Errorf("repo not found")
	}

	if repo.Remote {
		return errRemoteRepoNotSupport
	}

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

	var cc charts.Chart
	cc.CompressedData = data
	cc.Name = chartName
	cc.Version = opt.Version
	//
	err = db.CDB.CreateChartVersion(repoName, chartName, opt.Version, &cc)
	if err != nil {
		return fmt.Errorf("create chart fail when save chart to db: %v", err)
	}

	return nil
}
