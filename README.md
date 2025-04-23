# DeskMCP

An MCP (Model Context Protocol) server for the Teamwork Desk API. This server allows AI agents to interact with Teamwork Desk through a standardized protocol.

## Features

- Full support for Teamwork Desk API resources:
  - Tickets
  - Customers
  - Companies
  - Users
- Advanced filtering capabilities for all list operations
- JSON response formatting for easy parsing
- Docker support for easy deployment

## Installation

### Option 1: Direct Installation
```bash
go install github.com/ready4god2513/deskmcp/cmd/mcp
```

### Option 2: Docker Installation
1. Build the Docker image:
```bash
docker build -t deskmcp .
```

2. Run the container (recommended to use a secrets manager or environment file for sensitive data):
```bash
# Using environment variables (not recommended for production)
docker run -d \
  -e DESK_API_URL=https://yourcompany.teamwork.com \
  -e DESK_API_TOKEN=your_api_token \
  --name deskmcp \
  deskmcp

# Using environment file (more secure)
docker run -d \
  --env-file .env \
  --name deskmcp \
  deskmcp
```

3. To view logs:
```bash
docker logs -f deskmcp
```

4. To stop the container:
```bash
docker stop deskmcp
```

## Configuration

The server requires the following environment variables:

- `DESK_API_URL`: Your Teamwork Desk API URL
- `DESK_API_TOKEN`: Your Teamwork Desk API token

## Getting Started

### What is this tool?

DeskMCP is a bridge between your Teamwork Desk account and AI assistants like Claude. It allows AI assistants to perform actions in your Teamwork Desk account, such as:
- Creating and managing support tickets
- Looking up customer information
- Managing companies and users
- Generating reports and insights

### How to Use with AI Assistants

1. **Install the Tool**
   - Make sure you have Go installed on your computer
   - Run the installation command shown above
   - Set up your Teamwork Desk API credentials

2. **Start the Server**
   - Open a terminal window
   - Run the `mcp` command
   - Keep this terminal window open while using the tool

3. **Using with Claude**
   - Open Claude in your web browser
   - In your conversation, you can ask Claude to:
     - "Show me all open tickets"
     - "Create a new ticket for customer John Smith"
     - "Find all tickets from company Acme Inc"
     - "List all customers who haven't been contacted in 30 days"

4. **Tips for Best Results**
   - Be specific in your requests
   - Use natural language (e.g., "show me" instead of technical commands)
   - Ask for clarification if you're not sure what information is available
   - Request data in specific formats (e.g., "as a table" or "with dates sorted")
   - Ask for explanations of the data when needed

5. **Safety and Security**
   - The tool only has access to what your API token allows
   - All actions are logged in your Teamwork Desk account
   - You can revoke access at any time by changing your API token
   - The tool runs locally on your machine, keeping your data secure

## Available Tools

### Tickets
- `list_tickets`: List all tickets with optional filters
- `get_ticket`: Get a specific ticket by ID
- `create_ticket`: Create a new ticket

### Customers
- `list_customers`: List all customers with optional filters
- `get_customer`: Get a specific customer by ID
- `create_customer`: Create a new customer

### Companies
- `list_companies`: List all companies with optional filters
- `get_company`: Get a specific company by ID
- `create_company`: Create a new company

### Users
- `list_users`: List all users with optional filters
- `get_user`: Get a specific user by ID
- `create_user`: Create a new user

## Filter Usage

All list operations support filtering through the `filter` parameter. Here are some examples:

### Simple Equality Filters
```json
{
  "filter": {
    "status": "open",
    "priority": "high"
  }
}
```

### Multiple Conditions
```json
{
  "filter": {
    "status": "open",
    "priority": "high",
    "customer_id": "123"
  }
}
```

### Common Filter Fields

#### Tickets
- `status`: "open", "pending", "closed", etc.
- `priority`: "low", "medium", "high", "urgent"
- `customer_id`: Customer ID
- `agent_id`: Agent ID
- `company_id`: Company ID
- `created_at`: Date range
- `updated_at`: Date range

#### Customers
- `email`: Customer email
- `company_id`: Company ID
- `created_at`: Date range
- `updated_at`: Date range

#### Companies
- `name`: Company name
- `created_at`: Date range
- `updated_at`: Date range

#### Users
- `email`: User email
- `role`: User role
- `created_at`: Date range
- `updated_at`: Date range

## Use Cases

### 1. Customer Support Automation
- Automatically create tickets from customer emails
- Route tickets to appropriate agents based on content
- Generate responses using AI
- Track customer satisfaction metrics

### 2. Customer Relationship Management
- Track customer interactions across multiple channels
- Identify high-value customers
- Monitor customer satisfaction trends
- Generate customer reports

### 3. Team Management
- Monitor agent performance
- Track ticket resolution times
- Identify training opportunities
- Balance workload across team members

### 4. Business Intelligence
- Generate custom reports
- Track key performance indicators
- Analyze customer support trends
- Identify areas for improvement

### 5. Integration with Other Systems
- Connect with CRM systems
- Integrate with project management tools
- Sync with marketing automation platforms
- Link with accounting software

### 6. AI-Powered Features
- Automated ticket categorization
- Sentiment analysis of customer messages
- Smart ticket routing
- Predictive response suggestions
- Automated follow-up scheduling

## Claude Desktop Configuration

To use DeskMCP with Claude Desktop:

1. Make sure DeskMCP is running locally (either through direct installation or Docker)
2. Open Claude Desktop
3. Go to Settings (gear icon)
4. Navigate to the "Tools" section
5. Add a new MCP server with the following configuration:

```json
{
  "mcpServers": {
    "Teamwork Desk": {
      "command": "mcp",
      "args": [],
      "env": {
        "DESK_API_URL": "https://yourcompany.teamwork.com",
        "DESK_API_TOKEN": "your_api_token"
      },
      "autostart": true,
      "autorestart": true
    }
  }
}
```

Note: The `mcp` command must be in your system's PATH. If you installed it using `go install`, it should be available. If you're using Docker, you'll need to use the appropriate Docker command instead.

6. Save the configuration
7. Restart Claude Desktop

Now you can use natural language to interact with your Teamwork Desk data through Claude. For example:
- "Show me all open tickets"
- "Create a new ticket for customer John Smith"
- "Find all tickets from company Acme Inc"

## Troubleshooting

### Server Disconnected Error

If you see "MCP Teamwork Desk: Server disconnected" or "spawn mcp ENOENT" error:

1. Find the location of your mcp executable:
   ```bash
   which mcp
   ```

2. Verify the full path to mcp in your configuration:
   ```json
   {
     "mcpServers": {
       "Teamwork Desk": {
         "command": "/Users/YOUR_USERNAME/go/bin/mcp",  # Use the full path from `which mcp`
         "args": [],
         "env": {
           "DESK_API_URL": "https://yourcompany.teamwork.com",
           "DESK_API_TOKEN": "your_api_token"
         },
         "autostart": true,
         "autorestart": true
       }
     }
   }
   ```

3. If `which mcp` returns nothing, the executable isn't in your PATH. Add your Go bin directory to PATH:
   ```bash
   # Add to your ~/.zshrc or ~/.bash_profile
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
   Then restart your terminal and Claude Desktop.

4. Alternatively, reinstall the tool:
   ```bash
   go install github.com/ready4god2513/deskmcp/cmd/mcp@latest
   ```

5. For Docker users:
   ```bash
   # Use the full Docker command in the configuration
   {
     "mcpServers": {
       "Teamwork Desk": {
         "command": "docker",
         "args": ["run", "--rm", "-e", "DESK_API_URL=https://yourcompany.teamwork.com", "-e", "DESK_API_TOKEN=your_api_token", "deskmcp"],
         "autostart": true,
         "autorestart": true
       }
     }
   }
   ```

## License

MIT 