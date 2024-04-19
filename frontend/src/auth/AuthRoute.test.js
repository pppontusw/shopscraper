import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { ApiKeyProvider } from './ApiKey';
import AuthRoute from './AuthRoute';

describe('AuthRoute', () => {
  it('renders ApiKeyInput if apiKey is not set', () => {
    render(
      <MemoryRouter>
        <ApiKeyProvider>
          <AuthRoute />
        </ApiKeyProvider>
      </MemoryRouter>
    );
    expect(screen.getByLabelText('API Key *')).toBeInTheDocument();
  });

  it('redirects if apiKey is set', () => {
    localStorage.setItem('apiKey', '12345');
    render(
      <MemoryRouter>
        <ApiKeyProvider>
          <AuthRoute />
        </ApiKeyProvider>
      </MemoryRouter>
    );
    expect(screen.queryByLabelText('API Key')).not.toBeInTheDocument();
    localStorage.clear();
  });
});
