package controllers

import (
	"fmt"
	"strings"
	"xhelm/models/repo"
)

type RepoController struct {
	baseController
}

// Listrepos
// @Title repo
// @Description   获取仓库列表
// @Success 201 {object} repo.ListResponse
// @Failure 500
// @router /repositories [Get]
func (this *RepoController) ListRepos() {
	req := &repo.ListRequest{}
	this.setOutput(repo.List(req))
}

// Getrepo
// @Title repo
// @Description   获取仓库列表
// @Param repo path string true "vsphere集群 uuid"
// @Success 201 {object} repo.GetResponse
// @Failure 500
// @router /repositories/:repository [Get]
func (this *RepoController) GetRepo() {
	req := &repo.GetRequest{}
	repoName := this.Ctx.Input.Param(":repository")
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repository name"))
		return
	}
	req.Name = repoName
	this.setOutput(repo.Get(req))
}

// NewRepo
// @Title repo
// @Description   添加新仓库
// @Param body body repo.CreateRequest true "仓库数据"
// @Success 201 {string} create success!
// @Failure 500
// @router /repositories [Post]
func (this *RepoController) NewRepo() {
	req := &repo.CreateRequest{}
	if err := this.unmarshalRequestBodyJSON(req); err != nil {
		this.setOutput(nil, err)
		return
	}

	this.setOutput(repo.New(req))
}

// DeleteRepo
// @Title repo
// @Description   删除新仓库
// @Param body body repo.CreateRequest true "仓库数据"
// @Success 201 {string} create success!
// @Failure 500
// @router /repository/:repository [Delete]
func (this *RepoController) DeleteRepo() {
	req := &repo.DeleteRequest{}
	repoName := this.Ctx.Input.Param(":repository")
	if strings.TrimSpace(repoName) == "" {
		this.setOutput(nil, fmt.Errorf("invaild repository name"))
		return
	}
	req.Name = repoName

	this.setOutput(repo.Delete(req))
}
