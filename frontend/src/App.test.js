import React from 'react';
import App from './App';
import { CustomThemeProvider } from './theme/useTheme';
import { BrowserRouter as Router } from 'react-router-dom';
import { ApiKeyProvider } from './auth/ApiKey';
import { render, screen } from '@testing-library/react';

jest.mock('./routes/AppRoutes', () => () => <div>Mocked App Routes</div>);

describe('App component', () => {
  it('renders without crashing', () => {
    render(
      <CustomThemeProvider>
        <ApiKeyProvider>
          <Router>
            <App />
          </Router>
        </ApiKeyProvider>
      </CustomThemeProvider>
    );
    expect(screen.getByText('Mocked App Routes')).toBeInTheDocument();
  });
});
