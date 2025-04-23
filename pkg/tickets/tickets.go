package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ready4god2513/deskmcp/pkg/desk"
	"github.com/ready4god2513/deskmcp/pkg/utils"
	"github.com/ready4god2513/desksdkgo/models"
)

type TicketHandler struct {
	deskClient *desk.Client
}

func NewTicketHandler(deskClient *desk.Client) *TicketHandler {
	return &TicketHandler{
		deskClient: deskClient,
	}
}

func (h *TicketHandler) RegisterTools(s *server.MCPServer) {
	// List tickets
	s.AddTool(mcp.NewTool("list_tickets",
		mcp.WithDescription("List all tickets"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for tickets. Available fields:
- status: Filter by ticket status (e.g. "open", "closed", "pending")
- priority: Filter by priority level
- created_at: Filter by creation date
- updated_at: Filter by last update date
- customer_id: Filter by customer ID
- company_id: Filter by company ID
- assigned_user_id: Filter by assigned user ID`),
		),
		mcp.WithString("orderBy",
			mcp.Description("Order by field"),
			mcp.Enum("createdAt", "updatedAt"),
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
			mcp.Description("Number of tickets per page"),
			mcp.Min(1),
			mcp.Max(100),
		),
	), h.listTickets)

	// Count tickets
	s.AddTool(mcp.NewTool("count_tickets",
		mcp.WithDescription("Count all filtered tickets"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for tickets. Available fields:
- status: Filter by ticket status (e.g. "open", "closed", "pending")
- priority: Filter by priority level
- created_at: Filter by creation date
- updated_at: Filter by last update date
- customer_id: Filter by customer ID
- company_id: Filter by company ID
- assigned_user_id: Filter by assigned user ID`),
		),
	), h.countTickets)

	// Get ticket
	s.AddTool(mcp.NewTool("get_ticket",
		mcp.WithDescription("Get a specific ticket by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Ticket ID"),
		),
	), h.getTicket)

	// Create ticket
	s.AddTool(mcp.NewTool("create_ticket",
		mcp.WithDescription("Create a new ticket"),
		mcp.WithString("subject",
			mcp.Required(),
			mcp.Description("Ticket subject"),
		),
		mcp.WithString("preview_text",
			mcp.Required(),
			mcp.Description("Ticket preview text"),
		),
	), h.createTicket)
}

func (h *TicketHandler) listTickets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
		for field, value := range filterParams {
			params.Add(field, fmt.Sprintf("%v", value))
		}
	}

	utils.AddPaginationToParams(params, request)

	resp, err := h.deskClient.Client.Tickets.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list tickets: %v", err)), nil
	}

	// Format the tickets into a more readable structure
	type formattedTicket struct {
		ID          int    `json:"id"`
		Subject     string `json:"subject"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		PreviewText string `json:"preview_text"`
	}

	tickets := make([]formattedTicket, 0, len(resp.Tickets))
	for _, t := range resp.Tickets {
		var status string

		for _, s := range resp.Included.Ticketstatuses {
			if s.ID == t.Status.ID {
				status = s.Name
				break
			}
		}
		tickets = append(tickets, formattedTicket{
			ID:          t.ID,
			Subject:     t.Subject,
			Status:      status,
			CreatedAt:   t.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
			PreviewText: t.PreviewText,
		})
	}

	data, err := json.Marshal(tickets)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal tickets: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// Count tickets
func (h *TicketHandler) countTickets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := url.Values{}
	utils.AddPaginationToParams(params, request)

	tickets, err := h.deskClient.Client.Tickets.List(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to count tickets: %v", err)), nil
	}

	return mcp.NewToolResultText(strconv.Itoa(tickets.Pagination.Records)), nil
}

func (h *TicketHandler) getTicket(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid ticket ID: %v", err)), nil
	}
	resp, err := h.deskClient.Client.Tickets.Get(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get ticket: %v", err)), nil
	}
	data, err := json.Marshal(resp.Ticket)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func (h *TicketHandler) createTicket(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ticket := &models.Ticket{
		Subject:     request.Params.Arguments["subject"].(string),
		PreviewText: request.Params.Arguments["preview_text"].(string),
	}
	resp, err := h.deskClient.Client.Tickets.Create(ctx, ticket)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create ticket: %v", err)), nil
	}
	data, err := json.Marshal(resp.Ticket)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}
