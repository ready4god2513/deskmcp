package users

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

type UserHandler struct {
	deskClient *desk.Client
}

func NewUserHandler(deskClient *desk.Client) *UserHandler {
	return &UserHandler{
		deskClient: deskClient,
	}
}

func (h *UserHandler) RegisterTools(s *server.MCPServer) {
	// List users
	s.AddTool(mcp.NewTool("list_users",
		mcp.WithDescription("List all users"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for users. Available fields:
- email: Filter by email address
- first_name: Filter by first name
- last_name: Filter by last name
- role: Filter by user role
- created_at: Filter by creation date
- updated_at: Filter by last update date`),
		),
		mcp.WithString("orderBy",
			mcp.Description("Order by field"),
			mcp.Enum("createdAt", "updatedAt", "firstName", "lastName", "email", "role"),
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
			mcp.Description("Number of users per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listUsers)

	// Get user
	s.AddTool(mcp.NewTool("get_user",
		mcp.WithDescription("Get a specific user by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("User ID"),
		),
	), h.getUser)

	// Create user
	s.AddTool(mcp.NewTool("create_user",
		mcp.WithDescription("Create a new user"),
		mcp.WithString("first_name",
			mcp.Required(),
			mcp.Description("User's first name"),
		),
		mcp.WithString("last_name",
			mcp.Required(),
			mcp.Description("User's last name"),
		),
		mcp.WithString("email",
			mcp.Required(),
			mcp.Description("User's email address"),
		),
	), h.createUser)
}

func (h *UserHandler) listUsers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	resp, err := h.deskClient.Client.Users.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list users: %v", err)), nil
	}
	data, err := json.Marshal(resp.Users)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal users: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *UserHandler) getUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid user ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.Users.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get user: %v", err)), nil
	}
	data, err := json.Marshal(resp.User)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *UserHandler) createUser(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	user := &models.User{
		FirstName: request.Params.Arguments["first_name"].(string),
		LastName:  request.Params.Arguments["last_name"].(string),
		Email:     request.Params.Arguments["email"].(string),
	}
	resp, err := h.deskClient.Client.Users.Create(ctx, user)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create user: %v", err)), nil
	}
	data, err := json.Marshal(resp.User)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
