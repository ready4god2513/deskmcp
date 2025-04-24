package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/ready4god2513/deskmcp/pkg/companies"
	"github.com/ready4god2513/deskmcp/pkg/customers"
	"github.com/ready4god2513/deskmcp/pkg/desk"
	"github.com/ready4god2513/deskmcp/pkg/tags"
	"github.com/ready4god2513/deskmcp/pkg/tickets"
	"github.com/ready4god2513/deskmcp/pkg/ticketstatuses"
	"github.com/ready4god2513/deskmcp/pkg/tickettypes"
	"github.com/ready4god2513/deskmcp/pkg/users"
)

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
	deskClient := desk.NewClient(deskURL, deskToken)

	// Create MCP server
	s := server.NewMCPServer(
		"Teamwork Desk",
		"1.0.0",
	)

	// Register tools from each package
	ticketHandler := tickets.NewTicketHandler(deskClient)
	ticketHandler.RegisterTools(s)

	customerHandler := customers.NewCustomerHandler(deskClient)
	customerHandler.RegisterTools(s)

	companyHandler := companies.NewCompanyHandler(deskClient)
	companyHandler.RegisterTools(s)

	userHandler := users.NewUserHandler(deskClient)
	userHandler.RegisterTools(s)

	ticketStatusHandler := ticketstatuses.NewTicketStatusHandler(deskClient)
	ticketStatusHandler.RegisterTools(s)

	tagsHandler := tags.NewTagHandler(deskClient)
	tagsHandler.RegisterTools(s)

	ticketTypeHandler := tickettypes.NewTicketTypeHandler(deskClient)
	ticketTypeHandler.RegisterTools(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}
