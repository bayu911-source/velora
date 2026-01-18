import React, { useState, useEffect } from 'react';

const AgentList = () => {
  const [agents, setAgents] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('http://localhost:8080/agents')
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => setAgents(data))
      .catch((error) => setError(error.message));
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h2>Available Agents</h2>
      <ul>
        {agents.map((agent) => (
          <li key={agent.name}>
            <strong>{agent.name}</strong>: {agent.description}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default AgentList;