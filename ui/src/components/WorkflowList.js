import React, { useState, useEffect } from 'react';

const WorkflowList = () => {
  const [workflows, setWorkflows] = useState([]);
  const [selectedWorkflow, setSelectedWorkflow] = useState(null);
  const [agents, setAgents] = useState([]);
  const [error, setError] = useState(null);
  const [statusMessage, setStatusMessage] = useState('');
  const [form, setForm] = useState({
    name: '',
    description: '',
    agents: '',
  });
  const [runInput, setRunInput] = useState('');
  const [runOutput, setRunOutput] = useState('');

  useEffect(() => {
    fetchAgents();
    fetchWorkflows();
  }, []);

  const fetchAgents = () => {
    fetch('/api/agents')
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
    fetch('/api/workflows')
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
    fetch(`/api/workflows/${id}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error('Failed to load workflow details');
        }
        return response.json();
      })
      .then((data) => {
        setSelectedWorkflow(data);
        setRunOutput('');
        setError(null);
      })
      .catch((error) => setError(error.message));
  };

  const handleCreate = async (event) => {
    event.preventDefault();
    setStatusMessage('');
    setError(null);

    const agentsArray = form.agents
      .split(',')
      .map((agent) => agent.trim())
      .filter(Boolean);

    const response = await fetch('/api/workflows', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: form.name,
        description: form.description,
        agents: agentsArray,
      }),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Failed to create workflow');
      return;
    }

    const workflow = await response.json();
    setWorkflows((prev) => [...prev, workflow]);
    setForm({ name: '', description: '', agents: '' });
    setStatusMessage('Workflow created successfully');
  };

  const handleRun = async (event) => {
    event.preventDefault();
    if (!selectedWorkflow) {
      return;
    }
    setStatusMessage('Running workflow...');
    setError(null);

    const response = await fetch(`/api/workflows/${selectedWorkflow.id}/run`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ input: runInput }),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Workflow execution failed');
      setStatusMessage('');
      return;
    }

    const result = await response.json();
    setRunOutput(result.output || '');
    setStatusMessage(`Workflow completed: ${result.state}`);
    fetchWorkflows();
    fetchWorkflowDetails(selectedWorkflow.id);
  };

  const handleDelete = async (id) => {
    setStatusMessage('Deleting workflow...');
    setError(null);

    const response = await fetch(`/api/workflows/${id}`, {
      method: 'DELETE',
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
      setRunOutput('');
    }
    setStatusMessage('Workflow deleted successfully');
  };

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <h2>Workflows</h2>
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
              Agents (comma separated)
              <input name="agents" value={form.agents} onChange={handleChange} required />
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
                  <span>({workflow.state})</span>
                </div>
                <div className="workflow-actions">
                  <button type="button" className="danger" onClick={() => handleDelete(workflow.id)}>
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        </section>

        <section className="workflow-details">
          <h3>Workflow Details</h3>
          {selectedWorkflow ? (
            <div>
              <p><strong>Name:</strong> {selectedWorkflow.name}</p>
              <p><strong>Description:</strong> {selectedWorkflow.description || '—'}</p>
              <p><strong>State:</strong> {selectedWorkflow.state}</p>
              <div>
                <h4>Steps</h4>
                <ol>
                  {selectedWorkflow.steps.map((step, index) => (
                    <li key={`${step.AgentName}-${index}`}>
                      <strong>{step.AgentName}</strong>
                      <div>Input: {step.Input || '—'}</div>
                      <div>Output: {step.Output || '—'}</div>
                    </li>
                  ))}
                </ol>
              </div>

              <form onSubmit={handleRun} className="workflow-run-form">
                <label>
                  Run Input
                  <input value={runInput} onChange={(event) => setRunInput(event.target.value)} />
                </label>
                <button type="submit">Run Workflow</button>
              </form>

              {runOutput && (
                <div className="workflow-output">
                  <h4>Run Output</h4>
                  <pre>{runOutput}</pre>
                </div>
              )}
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
