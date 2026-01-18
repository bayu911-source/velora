# Database Schema for Workflow Persistence

This document outlines the database schema for storing the state of workflows in Velora.

## Tables

### `workflows`

Stores the definition of a workflow.

| Column      | Type    | Description                                      |
|-------------|---------|--------------------------------------------------|
| `id`          | TEXT    | Unique identifier for the workflow (e.g., UUID).   |
| `name`        | TEXT    | Human-readable name of the workflow.             |
| `description` | TEXT    | A brief description of what the workflow does.   |
| `created_at`  | INTEGER | Timestamp of when the workflow was created.      |

### `steps`

Stores the individual steps of a workflow.

| Column        | Type    | Description                                                |
|---------------|---------|------------------------------------------------------------|
| `id`            | TEXT    | Unique identifier for the step (e.g., UUID).             |
| `workflow_id`   | TEXT    | Foreign key referencing the `workflows` table.             |
| `agent_name`    | TEXT    | The name of the agent to be executed in this step.       |
| `input`         | TEXT    | The input to be passed to the agent.                     |
| `order`         | INTEGER | The order in which this step should be executed.           |

### `executions`

Stores the state of a running workflow execution.

| Column        | Type    | Description                                                |
|---------------|---------|------------------------------------------------------------|
| `id`            | TEXT    | Unique identifier for the execution (e.g., UUID).        |
| `workflow_id`   | TEXT    | Foreign key referencing the `workflows` table.             |
| `status`        | TEXT    | The current status of the execution (e.g., `running`, `completed`, `failed`). |
| `current_step`  | INTEGER | The index of the currently executing step.               |
| `started_at`    | INTEGER | Timestamp of when the execution started.                   |
| `finished_at`   | INTEGER | Timestamp of when the execution finished.                  |

## Relationships

- A `workflow` can have many `steps`.
- A `workflow` can have many `executions`.
- A `step` belongs to a `workflow`.
- An `execution` belongs to a `workflow`.
