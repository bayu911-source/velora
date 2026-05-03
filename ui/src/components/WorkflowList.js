import React, { useEffect, useState } from 'react';

const API_BASE = '/api/v1';

const WorkflowList = ({ token }) => {
  const [workflows, setWorkflows] = useState([]);
  const [selectedWorkflow, setSelectedWorkflow] = useState(null);
  const [agents, setAgents] = useState([]);
  const [error, setError] = useState(null);
  const [statusMessage, setStatusMessage] = useState('');
  const [form, setForm] = useState({
    name: '',
    description: '',
    agentIdentifiers: '',
    actionPrompt: '',
  });
  const [runInput, setRunInput] = useState('');

  useEffect(() => {
    if (token) {
      fetchAgents();
      fetchWorkflows();
    }
  }, [token]);

  const authHeaders = () => ({
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json',
  });

  const fetchAgents = () => {
    fetch(`${API_BASE}/agents`, { headers: authHeaders() })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => setAgents(data))
      .catch((error) => setError(error.message));
  };

  const fetchWorkflows = () => {
    fetch(`${API_BASE}/workflows`, { headers: authHeaders() })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => setWorkflows(data))
      .catch((error) => setError(error.message));
  };

  const fetchWorkflowDetails = (id) => {
    fetch(`${API_BASE}/workflows/${id}`, { headers: authHeaders() })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Failed to load workflow details');
        }
        return response.json();
      })
      .then((data) => {
        setSelectedWorkflow(data);
        setRunInput('');
        setError(null);
      })
      .catch((error) => setError(error.message));
  };

  const buildActions = () => {
    const identifiers = form.agentIdentifiers
      .split(',')
      .map((agent) => agent.trim())
      .filter(Boolean);

    if (identifiers.length === 0) {
      throw new Error('Please enter at least one agent identifier');
    }

    return identifiers.map((identifier) => {
      const agent = agents.find(
        (item) => item.id === identifier || item.name.toLowerCase() === identifier.toLowerCase(),
      );
      if (!agent) {
        throw new Error(`Agent not found: ${identifier}`);
      }
      return {
        name: agent.name,
        agent_id: agent.id,
        prompt: form.actionPrompt || agent.prompt || '',
      };
    });
  };

  const handleCreate = async (event) => {
    event.preventDefault();
    setStatusMessage('');
    setError(null);

    let actions = [];
    try {
      actions = buildActions();
    } catch (err) {
      setError(err.message);
      return;
    }

    const response = await fetch(`${API_BASE}/workflows`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({
        name: form.name,
        description: form.description,
        trigger: 'manual',
        actions,
      }),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Failed to create workflow');
      return;
    }

    const workflow = await response.json();
    setWorkflows((prev) => [...prev, workflow]);
    setForm({ name: '', description: '', agentIdentifiers: '', actionPrompt: '' });
    setStatusMessage('Workflow created successfully');
  };

  const handleRun = async (event) => {
    event.preventDefault();
    if (!selectedWorkflow) {
      return;
    }
    setStatusMessage('Enqueuing workflow run...');
    setError(null);

    const response = await fetch(`${API_BASE}/workflows/${selectedWorkflow.id}/run`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ input: runInput }),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Workflow execution failed');
      setStatusMessage('');
      return;
    }

    const result = await response.json();
    setStatusMessage(result.status ? `Workflow ${result.status}` : 'Workflow run queued');
    fetchWorkflows();
    fetchWorkflowDetails(selectedWorkflow.id);
  };

  const handleDelete = async (id) => {
    setStatusMessage('Deleting workflow...');
    setError(null);

    const response = await fetch(`${API_BASE}/workflows/${id}`, {
      method: 'DELETE',
      headers: authHeaders(),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Failed to delete workflow');
      setStatusMessage('');
      return;
    }

    setWorkflows((prev) => prev.filter((workflow) => workflow.id !== id));
    if (selectedWorkflow?.id === id) {
      setSelectedWorkflow(null);
      setRunInput('');
    }
    setStatusMessage('Workflow deleted successfully');
  };

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const renderWorkflowActions = () => {
    if (!selectedWorkflow) {
      return null;
    }

    let actions = [];
    if (Array.isArray(selectedWorkflow.actions)) {
      actions = selectedWorkflow.actions;
    } else if (selectedWorkflow.actions) {
      try {
        actions = JSON.parse(selectedWorkflow.actions);
      } catch (err) {
        actions = [];
      }
    }

    return (
      <div>
        <h4>Actions</h4>
        <ol>
          {actions.map((action, index) => (
            <li key={`${action.agent_id || action.name}-${index}`}>
              <strong>{action.name || action.agent_id}</strong>
              <div>Agent: {action.agent_id}</div>
              <div>Prompt: {action.prompt || '—'}</div>
            </li>
          ))}
        </ol>
      </div>
    );
  };

  if (!token) {
    return <div>Please log in to manage workflows.</div>;
  }

  return (
    <div className="dashboard-panel">
      <h2>Workflows</h2>
      {error && <div className="error">Error: {error}</div>}
      <div className="workflow-grid">
        <section className="workflow-list">
          <h3>Create Workflow</h3>
          <form onSubmit={handleCreate} className="workflow-form">
            <label>
              Name
              <input name="name" value={form.name} onChange={handleChange} required />
            </label>
            <label>
              Description
              <input name="description" value={form.description} onChange={handleChange} />
            </label>
            <label>
              Agents (name or ID, comma separated)
              <input
                name="agentIdentifiers"
                value={form.agentIdentifiers}
                onChange={handleChange}
                required
                placeholder="e.g. Lead Generator, agent-uuid"
              />
            </label>
            <label>
              Action Prompt
              <input
                name="actionPrompt"
                value={form.actionPrompt}
                onChange={handleChange}
                placeholder="Optional override prompt for all actions"
              />
            </label>
            <button type="submit">Create</button>
          </form>

          <h3>Saved Workflows</h3>
          <ul>
            {workflows.map((workflow) => (
              <li key={workflow.id}>
                <div className="workflow-row">
                  <button type="button" onClick={() => fetchWorkflowDetails(workflow.id)}>
                    {workflow.name}
                  </button>
                  <span>({workflow.status || 'unknown'})</span>
                </div>
                <div className="workflow-actions">
                  <button type="button" className="danger" onClick={() => handleDelete(workflow.id)}>
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>

          <div className="workflow-helper">
            <h4>Available agents</h4>
            <p>Use agent names or IDs above when creating a workflow.</p>
            <ul>
              {agents.slice(0, 10).map((agent) => (
                <li key={agent.id}>{agent.name} ({agent.id})</li>
              ))}
            </ul>
          </div>
        </section>

        <section className="workflow-details">
          <h3>Workflow Details</h3>
          {selectedWorkflow ? (
            <div>
              <p><strong>Name:</strong> {selectedWorkflow.name}</p>
              <p><strong>Description:</strong> {selectedWorkflow.description || '—'}</p>
              <p><strong>Status:</strong> {selectedWorkflow.status || 'unknown'}</p>
              <p><strong>Trigger:</strong> {selectedWorkflow.trigger || 'manual'}</p>
              {renderWorkflowActions()}

              <form onSubmit={handleRun} className="workflow-run-form">
                <label>
                  Run Input
                  <input value={runInput} onChange={(event) => setRunInput(event.target.value)} />
                </label>
                <button type="submit">Run Workflow</button>
              </form>
            </div>
          ) : (
            <p>Select a workflow to see details and run it.</p>
          )}
        </section>
      </div>
      {statusMessage && <div className="status-message">{statusMessage}</div>}
    </div>
  );
};

export default WorkflowList;
