package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CursusUserService interface {
	List(context.Context, *CursusUserQueryRequest) (*CursusUsers, *PaginationResponse, error)

	FindByID(context.Context, UserID) (*CursusUsers, error)
	FindByCursus(context.Context, CursusID) (*CursusUsers, error)
}

type CursusUserClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusUserClient) List(ctx context.Context, req *CursusUserQueryRequest) (*CursusUsers, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "cursus_users", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleCursusUsersPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusUserClient) FindByID(ctx context.Context, id UserID) (*CursusUsers, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("users/%s/cursus_users", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleCursusUsersResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusUserClient) FindByCursus(ctx context.Context, id CursusID) (*CursusUsers, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("cursus/%s/cursus_users", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleCursusUsersResponse(res)
}

func handleCursusUserResponse(res *http.Response) (*CursusUser, error) {
	var response CursusUser

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleCursusUsersResponse(res *http.Response) (*CursusUsers, error) {
	var response CursusUsers

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleCursusUsersPaginatedResponse(res *http.Response) (*CursusUsers, *PaginationResponse, error) {
	var response CursusUsers

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
