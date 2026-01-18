import React, { useState, useEffect } from 'react';

const WorkflowList = () => {
  const [workflows, setWorkflows] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('http://localhost:8080/workflows')
      .then((response) => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then((data) => setWorkflows(data))
      .catch((error) => setError(error.message));
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h2>Available Workflows</h2>
      <ul>
        {workflows.map((workflow) => (
          <li key={workflow.name}>
            <strong>{workflow.name}</strong>: {workflow.description}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default WorkflowList;