# Velora

Velora is a flexible and extensible framework for building and managing AI agents and workflows. It allows developers to create complex, multi-step processes by chaining together different AI agents, each specialized for a specific task.

## Features

*   **Agent-Based Architecture:** Build modular agents that can be reused across different workflows.
*   **Workflow Management:** Define, run, and monitor complex workflows with ease.
*   **Plugin System:** Extend Velora's functionality with custom plugins.
*   **Go & Gemini:** Built with the power of Go and Google's Gemini generative AI models.

## Getting Started

### Prerequisites

*   Go 1.18 or later
*   A valid `GEMINI_API_KEY` environment variable

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

To run the Velora agent, use the following command:

```bash
go run main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.
