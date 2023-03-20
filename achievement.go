package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AchievementID int

func (pID AchievementID) String() string {
	return strconv.Itoa(int(pID))
}

type AchievementService interface {
	List(context.Context, *AchievementQueryRequest) (*Achievements, *PaginationResponse, error)
	FindByCursus(context.Context, CursusID) (*Achievements, *PaginationResponse, error)
	FindByCampus(context.Context, CampusID) (*Achievements, *PaginationResponse, error)

	FindByTitle(context.Context, TitleID) (*Achievements, *PaginationResponse, error)

	FindByID(context.Context, AchievementID) (*Achievement, error)

	DeleteByID(context.Context, AchievementID) (*Achievement, error)
}

type AchievementClient struct {
	apiClient *Client
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) List(ctx context.Context, req *AchievementQueryRequest) (*Achievements, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, "achievements", "", req.Pagination.ToQuery(), nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) FindByCursus(ctx context.Context, id CursusID) (*Achievements, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("cursus/%s/achievements", id.String()), "", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) FindByCampus(ctx context.Context, id CampusID) (*Achievements, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("campus/%s/achievements", id.String()), "", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) FindByTitle(ctx context.Context, id TitleID) (*Achievements, *PaginationResponse, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("titles/%s/achievements", id.String()), "", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementsPaginatedResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) FindByID(ctx context.Context, id AchievementID) (*Achievement, error) {
	res, err := a.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("achievements/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementResponse(res)
}

// Get https://api.intra.42.fr/apidoc/2.0/achievements/show.html
func (a *AchievementClient) DeleteByID(ctx context.Context, id AchievementID) (*Achievement, error) {
	res, err := a.apiClient.request(ctx, http.MethodDelete, fmt.Sprintf("achievements/%s", id.String()), "", nil, nil)
	if err != nil {
		return nil, err
	}

	defer closeBody(res.Body)

	return handleAchievementResponse(res)
}

func handleAchievementResponse(res *http.Response) (*Achievement, error) {
	var response Achievement

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleAchievementsResponse(res *http.Response) (*Achievements, error) {
	var response Achievements

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func handleAchievementsPaginatedResponse(res *http.Response) (*Achievements, *PaginationResponse, error) {
	var response Achievements

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, nil, err
	}

	return &response, GetPaginationInfo(res.Header), nil
}
