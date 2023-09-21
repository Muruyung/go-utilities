package utils

import (
	"errors"
	"math"
)

// MetaResponse meta pagination response
type MetaResponse struct {
	Data Pagination `json:"pagination"`
}

// Pagination pagination attributes
type Pagination struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

// MapMetaResponse map meta pagination response
func MapMetaResponse(totalCount int, currentPageCount int, currentPage int, limitPerPage int) MetaResponse {
	totalPagesCount := math.Ceil(float64(totalCount) / float64(limitPerPage))
	return MetaResponse{
		Data: Pagination{
			Total:       totalCount,
			Count:       currentPageCount,
			PerPage:     limitPerPage,
			CurrentPage: currentPage,
			TotalPages:  int(totalPagesCount),
		},
	}
}

type paginationOption struct {
	limit  int
	page   int
	offset int
}

// NewPagination build new pagination
func NewPagination(page, limit int) *paginationOption {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return &paginationOption{
		limit:  limit,
		page:   page,
		offset: offset,
	}
}

func (pagination *paginationOption) GetLimit() int {
	return pagination.limit
}

func (pagination *paginationOption) GetOffset() int {
	return pagination.offset
}

func (pagination *paginationOption) GetPage() int {
	return pagination.page
}

// RequestOption pagination request option
type RequestOption struct {
	pagination *paginationOption
	sortBy     *map[string]direction
}

// NewRequestOption build new request option
func NewRequestOption() *RequestOption {
	return &RequestOption{
		pagination: nil,
		sortBy:     nil,
	}
}

// GetPagination get pagination request
func (request *RequestOption) GetPagination() *paginationOption {
	return request.pagination
}

// GetSortBy get sort by request
func (request *RequestOption) GetSortBy() *map[string]direction {
	return request.sortBy
}

// SetPagination set pagination request
func (request *RequestOption) SetPagination(pagination *paginationOption) *RequestOption {
	request.pagination = pagination
	return request
}

// SetSortBy set sort by request
func (request *RequestOption) SetSortBy(sortDir direction, sortBy ...string) (*RequestOption, error) {
	if sortDir.dir == "" {
		err := errors.New("invalid sort type value")
		return nil, err
	}

	sort := make(map[string]direction)
	if request.sortBy != nil {
		sort = *request.sortBy
	}

	for _, val := range sortBy {
		sort[val] = sortDir
	}
	request.sortBy = &sort

	return request, nil
}

// SetPaginationWithSort Set pagination with sort
func (request *RequestOption) SetPaginationWithSort(query QueryBuilderInteractor) (q QueryBuilderInteractor, page int, limit int) {
	var (
		sortBy     = request.GetSortBy()
		pagination = request.GetPagination()
	)
	query.AddPagination(pagination)

	if sortBy != nil {
		for sort, dir := range *sortBy {
			query.AddSort(dir, sort)
		}
	}

	if pagination != nil {
		page = pagination.GetPage()
		limit = pagination.GetLimit()
	}

	return query, page, limit
}
