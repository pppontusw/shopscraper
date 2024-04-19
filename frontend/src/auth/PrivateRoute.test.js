import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { ApiKeyProvider } from './ApiKey';
import PrivateRoute from './PrivateRoute';

describe('PrivateRoute', () => {
  it('renders children when apiKey is set', () => {
    localStorage.setItem('apiKey', '12345');
    render(
      <MemoryRouter>
        <ApiKeyProvider>
          <PrivateRoute><div>Protected Content</div></PrivateRoute>
        </ApiKeyProvider>
      </MemoryRouter>
    );
    expect(screen.getByText('Protected Content')).toBeInTheDocument();
    localStorage.clear();
  });

  it('redirects when apiKey is not set', () => {
    render(
      <MemoryRouter>
        <ApiKeyProvider>
          <PrivateRoute><div>Protected Content</div></PrivateRoute>
        </ApiKeyProvider>
      </MemoryRouter>
    );
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });
});
