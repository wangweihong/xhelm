package chart

type CreateRequest struct {
	Repo          string `json:"repository"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Description   string `json:"description"`
	Template      string `json:"template"`
	DefaultValues string `json:"default_values"`
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
	Repo    string `json:"repository"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
type GetResponse struct {
	Detail Detail `json:"detail"`
}

type DeleteRequest struct {
	Repo    string  `json:"repository"`
	Name    string  `json:"name"`
	Version *string `json:"version"`
}
type DeleteResponse struct{}
