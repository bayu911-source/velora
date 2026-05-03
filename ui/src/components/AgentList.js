import React, { useEffect, useState } from 'react';

const API_BASE = '/api/v1';

const AgentList = ({ token }) => {
  const [agents, setAgents] = useState([]);
  const [error, setError] = useState(null);
  const [statusMessage, setStatusMessage] = useState('');
  const [form, setForm] = useState({
    name: '',
    type: '',
    description: '',
    prompt: '',
  });

  useEffect(() => {
    if (token) {
      fetchAgents();
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

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);
    setStatusMessage('');

    const response = await fetch(`${API_BASE}/agents`, {
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify(form),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Failed to create agent');
      return;
    }

    setForm({ name: '', type: '', description: '', prompt: '' });
    setStatusMessage('Agent created successfully');
    fetchAgents();
  };

  if (!token) {
    return <div>Please log in to manage agents.</div>;
  }

  return (
    <div className="dashboard-panel">
      <h2>Available Agents</h2>
      {error && <div className="error">Error: {error}</div>}
      {statusMessage && <div className="status-message">{statusMessage}</div>}
      <ul>
        {agents.map((agent) => (
          <li key={agent.id || agent.name}>
            <strong>{agent.name}</strong> <em>({agent.type})</em>
            <div>{agent.description || 'No description provided.'}</div>
          </li>
        ))}
      </ul>

      <h3>Create Agent</h3>
      <form onSubmit={handleSubmit} className="workflow-form">
        <label>
          Name
          <input name="name" value={form.name} onChange={handleChange} required />
        </label>
        <label>
          Type
          <input name="type" value={form.type} onChange={handleChange} required />
        </label>
        <label>
          Description
          <input name="description" value={form.description} onChange={handleChange} />
        </label>
        <label>
          Prompt
          <input name="prompt" value={form.prompt} onChange={handleChange} />
        </label>
        <button type="submit">Create Agent</button>
      </form>
    </div>
  );
};

export default AgentList;
