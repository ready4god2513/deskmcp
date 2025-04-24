package ticketstatuses

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

type TicketStatusHandler struct {
	deskClient *desk.Client
}

func NewTicketStatusHandler(deskClient *desk.Client) *TicketStatusHandler {
	return &TicketStatusHandler{
		deskClient: deskClient,
	}
}

func (h *TicketStatusHandler) RegisterTools(s *server.MCPServer) {
	// List ticket statuses
	s.AddTool(mcp.NewTool("list_ticket_statuses",
		mcp.WithDescription("List all ticket statuses"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for ticket statuses. Available fields:
- name: Filter by ticket status name
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
			mcp.Description("Number of ticket statuses per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listTicketStatuses)

	// Get ticket status
	s.AddTool(mcp.NewTool("get_ticket_status",
		mcp.WithDescription("Get a specific ticket status by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Ticket Status ID"),
		),
	), h.getTicketStatus)

	// Create ticket status
	s.AddTool(mcp.NewTool("create_ticket_status",
		mcp.WithDescription("Create a new ticket status"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Ticket status name"),
		),
	), h.createTicketStatus)
}

func (h *TicketStatusHandler) listTicketStatuses(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	utils.AddPaginationToParams(params, request)

	resp, err := h.deskClient.Client.TicketStatuses.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list ticket statuses: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketStatuses)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket statuses: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TicketStatusHandler) getTicketStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid ticket status ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.TicketStatuses.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ticket status: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketStatus)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket status: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TicketStatusHandler) createTicketStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ticketStatus := &models.TicketStatus{
		Name: request.Params.Arguments["name"].(string),
	}
	resp, err := h.deskClient.Client.TicketStatuses.Create(ctx, ticketStatus)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create ticket status: %v", err)), nil
	}
	data, err := json.Marshal(resp.TicketStatus)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket status: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
