package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (pID ProjectID) String() string {
	return strconv.Itoa(int(pID))
}

type ProjectService interface {
	List(context.Context, *ProjectQueryRequest) (*Projects, *PaginationResponse, error)

	GetProjectsByCursus(context.Context, CursusID) (*Projects, *PaginationResponse, error)

	FindByID(context.Context, ProjectID) (*Project, error)
	DeleteByID(context.Context, ProjectID) (*Project, error)
}

type ProjectClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/projects/show.html
func (a *ProjectClient) List(ctx context.Context, req *ProjectQueryRequest) (*Projects, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "projects", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleProjectsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/projects/show.html
func (a *ProjectClient) GetProjectsByCursus(ctx context.Context, id CursusID) (*Projects, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("cursus/%s/projects", id.String()), "", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleProjectsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/projects/show.html
func (a *ProjectClient) FindByID(ctx context.Context, id ProjectID) (*Project, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("projects/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleProjectResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/projects/show.html
func (a *ProjectClient) DeleteByID(ctx context.Context, id ProjectID) (*Project, error) {
	res, err := a.apiClient.request(ctx, http.MethodDelete, fmt.Sprintf("projects/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleProjectResponse(res)
}

func handleProjectResponse(res *http.Response) (*Project, error) {
	var response Project

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleProjectsResponse(res *http.Response) (*Projects, error) {
	var response Projects

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleProjectsPaginatedResponse(res *http.Response) (*Projects, *PaginationResponse, error) {
	var response Projects

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
