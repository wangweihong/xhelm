package repo

type CreateRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CreateResponse struct{}

type ListRequest struct{}
type ListResponse struct {
	Total      int          `json:"total"`
	Reposities []Repository `json:"repositories"`
}

type GetRequest struct {
	Name string `json:"name"`
}
type GetResponse struct {
	Reposity Repository `json:"repository"`
}

type DeleteRequest struct {
	Name string `json:"name"`
}
type DeleteResponse struct{}

type Repository struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Remote     bool   `json:"is_remote"`
	CreateTime int64  `json:"create_time"`
}
