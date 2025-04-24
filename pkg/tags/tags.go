package tags

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

type TagHandler struct {
	deskClient *desk.Client
}

func NewTagHandler(deskClient *desk.Client) *TagHandler {
	return &TagHandler{
		deskClient: deskClient,
	}
}

func (h *TagHandler) RegisterTools(s *server.MCPServer) {
	// List tags
	s.AddTool(mcp.NewTool("list_tags",
		mcp.WithDescription("List all tags"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for tags. Available fields:
- name: Filter by tag name
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
			mcp.Description("Number of tags per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listTags)

	// Get tag
	s.AddTool(mcp.NewTool("get_tag",
		mcp.WithDescription("Get a specific tag by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Tag ID"),
		),
	), h.getTag)

	// Create tag
	s.AddTool(mcp.NewTool("create_tag",
		mcp.WithDescription("Create a new tag"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Tag name"),
		),
	), h.createTag)
}

func (h *TagHandler) listTags(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	utils.AddPaginationToParams(params, request)

	resp, err := h.deskClient.Client.Tags.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list tags: %v", err)), nil
	}
	data, err := json.Marshal(resp.Tags)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal tags: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TagHandler) getTag(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid tag ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.Tags.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tag: %v", err)), nil
	}
	data, err := json.Marshal(resp.Tag)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal tag: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TagHandler) createTag(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tag := &models.Tag{
		Name: request.Params.Arguments["name"].(string),
	}
	resp, err := h.deskClient.Client.Tags.Create(ctx, tag)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create tag: %v", err)), nil
	}
	data, err := json.Marshal(resp.Tag)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal tag: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
