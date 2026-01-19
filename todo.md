
# Velora Development TODO

This is a to-do list for tracking the development of the Velora project. This project aims to be a flexible and extensible framework for building and managing AI agents and workflows.

## Milestone: Initial Release (v0.1.0)

### Completed

- [x] Set up the basic project structure for a Go application.
- [x] Implement a `CodeGenerator` agent that can generate Go code.
- [x] Integrate the Gemini API for text generation.
- [x] Implement a command-line interface (CLI) using Cobra.
- [x] Implement an `Agent` interface to standardize agent behavior.
- [x] Create a `Registry` to manage and discover agents.
- [x] Implement a multimodal agent that can process text and images.
- [x] Create a conversational chat agent using the chat session feature of the Gemini SDK.
- [x] Create a `workflow` to allow agents to work together on complex tasks.
- [x] Create a dynamic agent registration mechanism to avoid hardcoding in `cmd/agent.go`.
- [x] Improve error handling throughout the application.
- [x] Implement a retry mechanism with exponential backoff for API requests.
- [x] Enhance the AppBuilderAgent to create application files and directories.
- [x] Add support for other large language models.
- [x] Develop a `velora-ui` to provide a graphical interface for managing agents and workflows.

## Milestone: Persistence and State Management (v0.2.0)

### Completed

- [x] Design a database schema for storing workflow state.
- [x] Choose a database technology (e.g., SQLite, PostgreSQL).
- [x] Implement a data access layer for interacting with the database.
- [x] Integrate the persistence layer into the workflow engine.

**Decision:** We will use **SQLite** for its simplicity and ease of embedding within the application. This avoids the need for a separate database server.

## Future Release (v0.3.0)

- [ ] Create a plugin system to allow the community to extend Velora's functionality.

## Future Release (v0.4.0)

- [ ] **Create a Workflow Engine:**
  - [x] Design a system for defining workflows that orchestrate multiple agents.
  - [x] Implement a `Workflow` struct that can be executed by the engine.
  - [x] Create a `WorkflowManager` to manage the lifecycle of workflows.
