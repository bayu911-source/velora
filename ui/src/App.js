import React from 'react';
import './App.css';
import AgentList from './components/AgentList';
import WorkflowList from './components/WorkflowList';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Velora UI</h1>
      </header>
      <main>
        <AgentList />
        <WorkflowList />
      </main>
    </div>
  );
}

export default App;