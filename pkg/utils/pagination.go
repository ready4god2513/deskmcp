package utils

import (
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// PaginationParams represents the pagination and sorting parameters
type PaginationParams struct {
	OrderBy   string
	OrderMode string
	Page      int
	PageSize  int
}

// DefaultPaginationParams returns the default pagination parameters
func DefaultPaginationParams() PaginationParams {
	return PaginationParams{
		OrderBy:   "createdAt",
		OrderMode: "desc",
		Page:      1,
		PageSize:  10,
	}
}

// AddPaginationToParams adds pagination and sorting parameters to the URL values
func AddPaginationToParams(params url.Values, request mcp.CallToolRequest) {
	// Get pagination parameters from request
	orderBy, ok := request.Params.Arguments["orderBy"].(string)
	if !ok {
		orderBy = "createdAt"
	}
	orderMode, ok := request.Params.Arguments["orderMode"].(string)
	if !ok {
		orderMode = "desc"
	}
	page, ok := request.Params.Arguments["page"].(float64)
	if !ok {
		page = 1
	}
	pageSize, ok := request.Params.Arguments["pageSize"].(float64)
	if !ok {
		pageSize = 10
	}

	// Add the parameters to the URL values
	params.Add("orderBy", orderBy)
	params.Add("orderMode", orderMode)
	params.Add("page", strconv.Itoa(int(page)))
	params.Add("pageSize", strconv.Itoa(int(pageSize)))
}
