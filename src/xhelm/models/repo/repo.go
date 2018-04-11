package repo

import (
	"xhelm/repository"
)

func transfrom(v *repository.Repository) Repository {
	var r Repository
	r.Name = v.Name
	r.URL = v.URL
	r.CreateTime = v.CreateTime.Unix()
	r.Remote = v.Remote
	return r
}

func List(req *ListRequest) (*ListResponse, error) {
	var resp ListResponse
	repos, err := repository.RM.ListRepos()
	if err != nil {
		return nil, err
	}

	resp.Total = len(repos)
	resp.Reposities = make([]Repository, 0)
	for _, v := range repos {
		r := transfrom(&v)
		resp.Reposities = append(resp.Reposities, r)
	}
	return &resp, nil
}

func Get(req *GetRequest) (*GetResponse, error) {
	var resp GetResponse
	or, err := repository.RM.GetRepo(req.Name)
	if err != nil {
		return nil, err
	}

	r := transfrom(or)
	resp.Reposity = r
	return &resp, nil
}

func New(req *CreateRequest) (*CreateResponse, error) {
	var resp CreateResponse
	if req.URL == "" {
		err := repository.RM.NewRepo(req.Name, nil)
		return &resp, err
	} else {
		var co repository.CreateOption
		co.URL = req.URL

		err := repository.RM.NewRepo(req.Name, &co)
		return &resp, err
	}
}

func Delete(req *DeleteRequest) (*DeleteResponse, error) {
	var resp DeleteResponse
	err := repository.RM.DeleteRepo(req.Name)
	return &resp, err
}
