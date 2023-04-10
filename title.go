package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TitleService interface {
	List(ctx context.Context, req *TitleQueryRequest) (*Titles, *PaginationResponse, error)

	FindByID(ctx context.Context, id TitleID) (*Title, error)
}

type TitleClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/titles/show.html
func (a *TitleClient) List(ctx context.Context, req *TitleQueryRequest) (*Titles, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "titles", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleTitlesPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/titles/show.html
func (a *TitleClient) FindByID(ctx context.Context, id TitleID) (*Title, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("titles/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleTitleResponse(res)
}

func handleTitleResponse(res *http.Response) (*Title, error) {
	var response Title

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleTitlesResponse(res *http.Response) (*Titles, error) {
	var response Titles

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleTitlesPaginatedResponse(res *http.Response) (*Titles, *PaginationResponse, error) {
	var response Titles

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
