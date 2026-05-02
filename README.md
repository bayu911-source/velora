# Velora

Velora is a flexible and extensible framework for building and managing AI agents and workflows. It allows developers to create complex, multi-step processes by chaining together different AI agents, each specialized for a specific task.

## Features

*   **Agent-Based Architecture:** Build modular agents that can be reused across different workflows.
*   **Workflow Management:** Define, run, and monitor complex workflows with ease.
*   **Plugin System:** Extend Velora's functionality with custom plugins.
*   **Multi-LLM Support:** Supports Google's Gemini and OpenAI models.
*   **Go & AI:** Built with the power of Go and modern AI APIs.

## Getting Started

### Prerequisites

*   Go 1.18 or later
*   A valid API key for your preferred LLM provider:
  - `GEMINI_API_KEY` for Google Gemini
  - `OPENAI_API_KEY` for OpenAI (optional)

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/velora-id/velora.git
    ```

2.  Navigate to the project directory:

    ```bash
    cd velora
    ```

3.  Install dependencies:

    ```bash
    go mod tidy
    ```

## Usage

To run the Velora CLI, use the following command:

```bash
go run main.go
```

### CLI Commands

- `velora agent list`: List all available agents
- `velora workflow create <name> <agents>`: Create a new workflow
- `velora workflow run <id> <input>`: Run a workflow
- `velora workflow load <file>`: Load a workflow from YAML file
- `velora workflow list`: List saved workflows
- `velora workflow show <id>`: Show workflow details and step history
- `velora workflow delete <id>`: Delete a workflow and its history
- `velora server`: Start the HTTP server for UI access

### HTTP Server Endpoints

- `GET /api/agents`: List available agents
- `GET /api/workflows`: List saved workflows
- `POST /api/workflows`: Create a new workflow with JSON body `{ "name": "...", "description": "...", "agents": ["AgentA", "AgentB"] }`
- `GET /api/workflows/{id}`: Show workflow details and steps
- `POST /api/workflows/{id}/run`: Run a saved workflow with JSON body `{ "input": "..." }`
- `DELETE /api/workflows/{id}`: Delete a saved workflow

### UI

The Velora UI is built with React and served from `ui/build` when the backend server is running.

To develop the UI locally:

```bash
cd ui
npm install
npm start
```

To build the production UI assets:

```bash
cd ui
npm install
npm run build
```

Then start the backend with `go run main.go` and open the UI from the built files or use the React dev server with proxy.

### Environment Variables

- `GEMINI_API_KEY`: Your Google Gemini API key
- `OPENAI_API_KEY`: Your OpenAI API key (optional)
- `DATABASE_PATH`: Path to SQLite database (default: velora.db)

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.
