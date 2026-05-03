import React, { useEffect, useState } from 'react';
import './App.css';
import AgentList from './components/AgentList';
import WorkflowList from './components/WorkflowList';

const API_BASE = '/api/v1';

function App() {
  const [token, setToken] = useState(localStorage.getItem('velora_access_token') || '');
  const [refreshToken, setRefreshToken] = useState(localStorage.getItem('velora_refresh_token') || '');
  const [tenant, setTenant] = useState(JSON.parse(localStorage.getItem('velora_tenant') || 'null'));
  const [user, setUser] = useState(JSON.parse(localStorage.getItem('velora_user') || 'null'));
  const [page, setPage] = useState(token ? 'dashboard' : 'login');
  const [error, setError] = useState('');
  const [statusMessage, setStatusMessage] = useState('');
  const [loginForm, setLoginForm] = useState({ tenant_id: '', email: '', password: '' });
  const [registerForm, setRegisterForm] = useState({ tenant_name: '', name: '', email: '', password: '' });

  useEffect(() => {
    if (token) {
      setPage('dashboard');
    }
  }, [token]);

  const saveAuth = (accessToken, refreshTokenValue, userObj, tenantObj) => {
    setToken(accessToken);
    setRefreshToken(refreshTokenValue);
    setUser(userObj);
    setTenant(tenantObj);
    localStorage.setItem('velora_access_token', accessToken);
    localStorage.setItem('velora_refresh_token', refreshTokenValue);
    localStorage.setItem('velora_user', JSON.stringify(userObj));
    localStorage.setItem('velora_tenant', JSON.stringify(tenantObj));
  };

  const clearAuth = () => {
    setToken('');
    setRefreshToken('');
    setUser(null);
    setTenant(null);
    localStorage.removeItem('velora_access_token');
    localStorage.removeItem('velora_refresh_token');
    localStorage.removeItem('velora_user');
    localStorage.removeItem('velora_tenant');
    setPage('login');
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    setError('');
    setStatusMessage('Logging in...');

    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginForm),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Login failed');
      setStatusMessage('');
      return;
    }

    const data = await response.json();
    saveAuth(data.access_token, data.refresh_token, data.user, { id: loginForm.tenant_id });
    setStatusMessage('Login successful');
  };

  const handleRegister = async (event) => {
    event.preventDefault();
    setError('');
    setStatusMessage('Creating tenant...');

    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(registerForm),
    });

    if (!response.ok) {
      const message = await response.text();
      setError(message || 'Registration failed');
      setStatusMessage('');
      return;
    }

    const data = await response.json();
    saveAuth(data.access_token, data.refresh_token, data.user, data.tenant);
    setStatusMessage('Registration successful');
  };

  const handleChange = (event, formType) => {
    const { name, value } = event.target;
    if (formType === 'login') {
      setLoginForm((prev) => ({ ...prev, [name]: value }));
    } else {
      setRegisterForm((prev) => ({ ...prev, [name]: value }));
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Velora SaaS Dashboard</h1>
        {token && (
          <div className="workspace-bar">
            <div>
              <strong>Tenant:</strong> {tenant?.name || tenant?.id || 'Unknown'}
            </div>
            <div>
              <strong>User:</strong> {user?.name || user?.email}
            </div>
            <button type="button" onClick={clearAuth} className="logout-button">
              Logout
            </button>
          </div>
        )}
      </header>

      <main>
        {!token ? (
          <div className="auth-page">
            <div className="auth-tabs">
              <button className={page === 'login' ? 'auth-tab active' : 'auth-tab'} onClick={() => setPage('login')}>
                Login
              </button>
              <button className={page === 'register' ? 'auth-tab active' : 'auth-tab'} onClick={() => setPage('register')}>
                Register
              </button>
            </div>
            <div className="auth-card">
              {error && <div className="error">{error}</div>}
              {statusMessage && <div className="status-message">{statusMessage}</div>}
              {page === 'login' ? (
                <form onSubmit={handleLogin} className="auth-form">
                  <label>
                    Tenant ID
                    <input name="tenant_id" value={loginForm.tenant_id} onChange={(e) => handleChange(e, 'login')} required />
                  </label>
                  <label>
                    Email
                    <input name="email" type="email" value={loginForm.email} onChange={(e) => handleChange(e, 'login')} required />
                  </label>
                  <label>
                    Password
                    <input name="password" type="password" value={loginForm.password} onChange={(e) => handleChange(e, 'login')} required />
                  </label>
                  <button type="submit">Login</button>
                </form>
              ) : (
                <form onSubmit={handleRegister} className="auth-form">
                  <label>
                    Tenant Name
                    <input name="tenant_name" value={registerForm.tenant_name} onChange={(e) => handleChange(e, 'register')} required />
                  </label>
                  <label>
                    Full Name
                    <input name="name" value={registerForm.name} onChange={(e) => handleChange(e, 'register')} required />
                  </label>
                  <label>
                    Email
                    <input name="email" type="email" value={registerForm.email} onChange={(e) => handleChange(e, 'register')} required />
                  </label>
                  <label>
                    Password
                    <input name="password" type="password" value={registerForm.password} onChange={(e) => handleChange(e, 'register')} required />
                  </label>
                  <button type="submit">Register</button>
                </form>
              )}
            </div>
          </div>
        ) : (
          <div className="dashboard-grid">
            <AgentList token={token} />
            <WorkflowList token={token} />
          </div>
        )}
      </main>
    </div>
  );
}

export default App;