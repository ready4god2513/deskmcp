package companies

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ready4god2513/deskmcp/pkg/desk"
	"github.com/ready4god2513/deskmcp/pkg/utils"
	"github.com/ready4god2513/desksdkgo/models"
)

type CompanyHandler struct {
	deskClient *desk.Client
}

func NewCompanyHandler(deskClient *desk.Client) *CompanyHandler {
	return &CompanyHandler{
		deskClient: deskClient,
	}
}

func (h *CompanyHandler) RegisterTools(s *server.MCPServer) {
	// List companies
	s.AddTool(mcp.NewTool("list_companies",
		mcp.WithDescription("List all companies"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for companies. Available fields:
- name: Filter by company name
- created_at: Filter by creation date
- updated_at: Filter by last update date`),
		),
		mcp.WithString("orderBy",
			mcp.Description("Order by field"),
			mcp.Enum("createdAt", "updatedAt", "name"),
		),
		mcp.WithString("orderMode",
			mcp.Description("Order mode"),
			mcp.Enum("asc", "desc"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number"),
			mcp.Min(1),
		),
		mcp.WithNumber("pageSize",
			mcp.Description("Number of companies per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listCompanies)

	// Get company
	s.AddTool(mcp.NewTool("get_company",
		mcp.WithDescription("Get a specific company by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Company ID"),
		),
	), h.getCompany)

	// Create company
	s.AddTool(mcp.NewTool("create_company",
		mcp.WithDescription("Create a new company"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Company name"),
		),
	), h.createCompany)
}

func (h *CompanyHandler) listCompanies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	utils.AddPaginationToParams(params, request)

	resp, err := h.deskClient.Client.Companies.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list companies: %v", err)), nil
	}
	data, err := json.Marshal(resp.Companies)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal companies: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *CompanyHandler) getCompany(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid company ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.Companies.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get company: %v", err)), nil
	}
	data, err := json.Marshal(resp.Company)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal company: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *CompanyHandler) createCompany(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	company := &models.Company{
		Name: request.Params.Arguments["name"].(string),
	}
	resp, err := h.deskClient.Client.Companies.Create(ctx, company)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create company: %v", err)), nil
	}
	data, err := json.Marshal(resp.Company)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal company: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
