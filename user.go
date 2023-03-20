package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (pID UserID) String() string {
	return strconv.Itoa(int(pID))
}

type UserService interface {
	Me(ctx context.Context, token string) (*User, error)

	List(ctx context.Context, req *CursusQueryRequest) (*Users, *PaginationResponse, error)

	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByCampus(ctx context.Context, id CursusID) (*Users, error)

	LocationStats(ctx context.Context, id UserID) (*LocationsStat, error)
}

type UserClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/Users/show.html
func (a *UserClient) List(ctx context.Context, req *CursusQueryRequest) (*Users, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "users", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleUsersPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/Users/show.html
func (a *UserClient) FindByID(ctx context.Context, id UserID) (*User, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("users/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleUserResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/Users/show.html
func (a *UserClient) FindByCampus(ctx context.Context, id CursusID) (*Users, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("cursus/%s/users", id), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleUsersResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/Users/show.html
func (a *UserClient) Me(ctx context.Context, tok string) (*User, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "me", tok, nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleUserResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/Users/show.html
func (a *UserClient) LocationStats(ctx context.Context, id UserID) (*LocationsStat, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("users/%s/locations_stats", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	var response LocationsStat

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleUserResponse(res *http.Response) (*User, error) {
	var response User

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleUsersResponse(res *http.Response) (*Users, error) {
	var response Users

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleUsersPaginatedResponse(res *http.Response) (*Users, *PaginationResponse, error) {
	var response Users

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
