package tickettypes

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

type TicketTypeHandler struct {
	deskClient *desk.Client
}

func NewTicketTypeHandler(deskClient *desk.Client) *TicketTypeHandler {
	return &TicketTypeHandler{
		deskClient: deskClient,
	}
}

func (h *TicketTypeHandler) RegisterTools(s *server.MCPServer) {
	// List ticket types
	s.AddTool(mcp.NewTool("list_ticket_types",
		mcp.WithDescription("List all ticket types"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for ticket types. Available fields:
- name: Filter by ticket type name
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
			mcp.Description("Number of ticket types per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listTicketTypes)

	// Get ticket type
	s.AddTool(mcp.NewTool("get_ticket_type",
		mcp.WithDescription("Get a specific ticket type by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Ticket Type ID"),
		),
	), h.getTicketType)

	// Create ticket type
	s.AddTool(mcp.NewTool("create_ticket_type",
		mcp.WithDescription("Create a new ticket type"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Ticket type name"),
		),
	), h.createTicketType)
}

func (h *TicketTypeHandler) listTicketTypes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	utils.AddPaginationToParams(params, request)

	resp, err := h.deskClient.Client.TicketTypes.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list ticket types: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketTypes)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket types: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TicketTypeHandler) getTicketType(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid ticket type ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.TicketTypes.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ticket type: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket type: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TicketTypeHandler) createTicketType(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ticketType := &models.TicketType{
		Name: request.Params.Arguments["name"].(string),
	}
	resp, err := h.deskClient.Client.TicketTypes.Create(ctx, ticketType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create ticket type: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket type: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
