package chart

type CreateRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CreateResponse struct{}

type ListRequest struct {
	Repo string `json:"repository"`
}
type ListResponse struct {
	Total  int     `json:"total"`
	Charts []Chart `json:"charts"`
}

type GetRequest struct {
	Repo string `json:"repository"`
	Name string `json:"name"`
}
type GetResponse struct {
	Chart Chart `json:"chart"`
}

type DeleteRequest struct {
	Repo string `json:"repository"`
	Name string `json:"name"`
}
type DeleteResponse struct{}
