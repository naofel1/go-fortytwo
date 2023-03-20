package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (pID CursusID) String() string {
	return strconv.Itoa(int(pID))
}

type CursusService interface {
	List(context.Context, *CursusQueryRequest) (*CursusSlice, *PaginationResponse, error)

	FindByID(context.Context, CursusID) (*Cursus, error)
	DeleteByID(context.Context, CursusID) (*Cursus, error)
}

type CursusClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusClient) List(ctx context.Context, req *CursusQueryRequest) (*CursusSlice, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "cursus", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleCursusSlicePaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusClient) FindByID(ctx context.Context, id CursusID) (*Cursus, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("cursus/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleCursusResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/cursus/show.html
func (a *CursusClient) DeleteByID(ctx context.Context, id CursusID) (*Cursus, error) {
	res, err := a.apiClient.request(ctx, http.MethodDelete, fmt.Sprintf("cursus/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleCursusResponse(res)
}

func handleCursusResponse(res *http.Response) (*Cursus, error) {
	var response Cursus

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleCursusSliceResponse(res *http.Response) (*CursusSlice, error) {
	var response CursusSlice

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleCursusSlicePaginatedResponse(res *http.Response) (*CursusSlice, *PaginationResponse, error) {
	var response CursusSlice

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
