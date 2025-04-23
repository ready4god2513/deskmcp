package customers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ready4god2513/deskmcp/pkg/desk"
	"github.com/ready4god2513/desksdkgo/models"
)

type CustomerHandler struct {
	deskClient *desk.Client
}

func NewCustomerHandler(deskClient *desk.Client) *CustomerHandler {
	return &CustomerHandler{
		deskClient: deskClient,
	}
}

func (h *CustomerHandler) RegisterTools(s *server.MCPServer) {
	// List customers
	s.AddTool(mcp.NewTool("list_customers",
		mcp.WithDescription("List all customers"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for customers. Available fields:
- email: Filter by email address
- first_name: Filter by first name
- last_name: Filter by last name
- company_id: Filter by company ID
- created_at: Filter by creation date
- updated_at: Filter by last update date`),
		),
		mcp.WithString("orderBy",
			mcp.Description("Order by field"),
			mcp.Enum("createdAt", "updatedAt", "firstName", "lastName", "email"),
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
			mcp.Description("Number of customers per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listCustomers)

	// Get customer
	s.AddTool(mcp.NewTool("get_customer",
		mcp.WithDescription("Get a specific customer by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Customer ID"),
		),
	), h.getCustomer)

	// Create customer
	s.AddTool(mcp.NewTool("create_customer",
		mcp.WithDescription("Create a new customer"),
		mcp.WithString("first_name",
			mcp.Required(),
			mcp.Description("Customer's first name"),
		),
		mcp.WithString("last_name",
			mcp.Required(),
			mcp.Description("Customer's last name"),
		),
		mcp.WithString("email",
			mcp.Required(),
			mcp.Description("Customer's email address"),
		),
	), h.createCustomer)
}

func (h *CustomerHandler) listCustomers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	orderBy := request.Params.Arguments["orderBy"].(string)
	orderMode := request.Params.Arguments["orderMode"].(string)
	page := request.Params.Arguments["page"].(float64)
	pageSize := request.Params.Arguments["pageSize"].(float64)

	if orderBy == "" {
		orderBy = "createdAt"
	}
	if orderMode == "" {
		orderMode = "desc"
	}

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	params.Add("orderBy", orderBy)
	params.Add("orderMode", orderMode)
	params.Add("page", strconv.Itoa(int(page)))
	params.Add("pageSize", strconv.Itoa(int(pageSize)))

	resp, err := h.deskClient.Client.Customers.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list customers: %v", err)), nil
	}
	data, err := json.Marshal(resp.Customers)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customers: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *CustomerHandler) getCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid customer ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.Customers.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get customer: %v", err)), nil
	}
	data, err := json.Marshal(resp.Customer)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customer: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *CustomerHandler) createCustomer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	customer := &models.Customer{
		FirstName: request.Params.Arguments["first_name"].(string),
		LastName:  request.Params.Arguments["last_name"].(string),
		Email:     request.Params.Arguments["email"].(string),
	}
	resp, err := h.deskClient.Client.Customers.Create(ctx, customer)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create customer: %v", err)), nil
	}
	data, err := json.Marshal(resp.Customer)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customer: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
