package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ready4god2513/deskmcp/pkg/desk"
	"github.com/ready4god2513/desksdkgo/models"
)

var deskClient *desk.Client

func main() {
	// Get environment variables
	deskURL := os.Getenv("DESK_API_URL")
	if deskURL == "" {
		log.Fatal("DESK_API_URL environment variable is required")
	}

	deskToken := os.Getenv("DESK_API_TOKEN")
	if deskToken == "" {
		log.Fatal("DESK_API_TOKEN environment variable is required")
	}

	// Initialize Desk client
	deskClient = desk.NewClient(deskURL, deskToken)

	// Create MCP server
	s := server.NewMCPServer(
		"Teamwork Desk",
		"1.0.0",
	)

	// Register tools
	registerTicketTools(s)
	registerCustomerTools(s)
	registerCompanyTools(s)
	registerUserTools(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

func registerTicketTools(s *server.MCPServer) {
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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := url.Values{}
		if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
			for field, value := range filterParams {
				params.Add(field, fmt.Sprintf("%v", value))
			}
		}
		resp, err := deskClient.Client.Tickets.List(ctx, params)
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

		// Return the data in the format expected by Claude
		return mcp.NewToolResultText(string(data)), nil
	})

	// Get ticket
	s.AddTool(mcp.NewTool("get_ticket",
		mcp.WithDescription("Get a specific ticket by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Ticket ID"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid ticket ID: %v", err)), nil
		}
		resp, err := deskClient.Client.Tickets.Get(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get ticket: %v", err)), nil
		}
		data, err := json.Marshal(resp.Ticket)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ticket := &models.Ticket{
			Subject:     request.Params.Arguments["subject"].(string),
			PreviewText: request.Params.Arguments["preview_text"].(string),
		}
		resp, err := deskClient.Client.Tickets.Create(ctx, ticket)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create ticket: %v", err)), nil
		}
		data, err := json.Marshal(resp.Ticket)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticket: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerCustomerTools(s *server.MCPServer) {
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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := url.Values{}
		if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
			for field, value := range filterParams {
				params.Add(field, fmt.Sprintf("%v", value))
			}
		}
		resp, err := deskClient.Client.Customers.List(ctx, params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list customers: %v", err)), nil
		}
		data, err := json.Marshal(resp.Customers)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customers: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	// Get customer
	s.AddTool(mcp.NewTool("get_customer",
		mcp.WithDescription("Get a specific customer by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Customer ID"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid customer ID: %v", err)), nil
		}
		resp, err := deskClient.Client.Customers.Get(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get customer: %v", err)), nil
		}
		data, err := json.Marshal(resp.Customer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customer: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		customer := &models.Customer{
			FirstName: request.Params.Arguments["first_name"].(string),
			LastName:  request.Params.Arguments["last_name"].(string),
			Email:     request.Params.Arguments["email"].(string),
		}
		resp, err := deskClient.Client.Customers.Create(ctx, customer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create customer: %v", err)), nil
		}
		data, err := json.Marshal(resp.Customer)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal customer: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerCompanyTools(s *server.MCPServer) {
	// List companies
	s.AddTool(mcp.NewTool("list_companies",
		mcp.WithDescription("List all companies"),
		mcp.WithObject("filter",
			mcp.Description(`Optional filter for companies. Available fields:
- name: Filter by company name
- created_at: Filter by creation date
- updated_at: Filter by last update date`),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := url.Values{}
		if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
			for field, value := range filterParams {
				params.Add(field, fmt.Sprintf("%v", value))
			}
		}
		resp, err := deskClient.Client.Companies.List(ctx, params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list companies: %v", err)), nil
		}
		data, err := json.Marshal(resp.Companies)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal companies: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	// Get company
	s.AddTool(mcp.NewTool("get_company",
		mcp.WithDescription("Get a specific company by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("Company ID"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid company ID: %v", err)), nil
		}
		resp, err := deskClient.Client.Companies.Get(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get company: %v", err)), nil
		}
		data, err := json.Marshal(resp.Company)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal company: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	// Create company
	s.AddTool(mcp.NewTool("create_company",
		mcp.WithDescription("Create a new company"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Company name"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		company := &models.Company{
			Name: request.Params.Arguments["name"].(string),
		}
		resp, err := deskClient.Client.Companies.Create(ctx, company)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create company: %v", err)), nil
		}
		data, err := json.Marshal(resp.Company)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal company: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})
}

func registerUserTools(s *server.MCPServer) {
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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := url.Values{}
		if filterParams, ok := request.Params.Arguments["filter"].(map[string]interface{}); ok {
			for field, value := range filterParams {
				params.Add(field, fmt.Sprintf("%v", value))
			}
		}
		resp, err := deskClient.Client.Users.List(ctx, params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list users: %v", err)), nil
		}
		data, err := json.Marshal(resp.Users)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal users: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

	// Get user
	s.AddTool(mcp.NewTool("get_user",
		mcp.WithDescription("Get a specific user by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("User ID"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := strconv.Atoi(request.Params.Arguments["id"].(string))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid user ID: %v", err)), nil
		}
		resp, err := deskClient.Client.Users.Get(ctx, id)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user: %v", err)), nil
		}
		data, err := json.Marshal(resp.User)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})

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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user := &models.User{
			FirstName: request.Params.Arguments["first_name"].(string),
			LastName:  request.Params.Arguments["last_name"].(string),
			Email:     request.Params.Arguments["email"].(string),
		}
		resp, err := deskClient.Client.Users.Create(ctx, user)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create user: %v", err)), nil
		}
		data, err := json.Marshal(resp.User)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal user: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	})
}
