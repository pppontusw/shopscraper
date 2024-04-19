// frontend/src/components/ThemeToggleButton.test.js
import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import ThemeToggleButton from './ThemeButton';
import { CustomThemeProvider } from '../theme/useTheme';

// Helper function to render the component within the custom theme provider
const renderWithTheme = (component) => {
  return render(
    <CustomThemeProvider>
      {component}
    </CustomThemeProvider>
  );
};

describe('ThemeToggleButton', () => {
  beforeEach(() => {
    localStorage.clear(); // Clear localStorage before each test
  });

  it('toggles theme from dark to light and updates the icon and localStorage accordingly', () => {
    renderWithTheme(<ThemeToggleButton />);

    // Initially, the theme should be 'dark' in localStorage
    expect(localStorage.getItem('themeName')).toBe('dark');

    // Initially, find the SVG element (icon) for the light theme to toggle to dark
    let icon = screen.getByTestId('light-theme-icon');
    expect(icon).toBeInTheDocument();

    // Click the button to toggle the theme to 'light'
    fireEvent.click(screen.getByRole('button'));

    // After clicking, check if the icon has changed to dark theme icon
    icon = screen.getByTestId('dark-theme-icon');
    expect(icon).toBeInTheDocument();

    // Check localStorage to confirm it now holds the 'light' theme
    expect(localStorage.getItem('themeName')).toBe('light');
  });

  it('toggles theme from light to dark and updates the icon and localStorage accordingly', () => {
    // Set initial theme to 'light'
    localStorage.setItem('themeName', 'light');
    renderWithTheme(<ThemeToggleButton />);

    // Initially, the theme should be 'light' in localStorage
    expect(localStorage.getItem('themeName')).toBe('light');

    // Initially, find the SVG element (icon) for the dark theme to toggle to light
    let icon = screen.getByTestId('dark-theme-icon');
    expect(icon).toBeInTheDocument();

    // Click the button to toggle the theme to 'dark'
    fireEvent.click(screen.getByRole('button'));

    // After clicking, check if the icon has changed to light theme icon
    icon = screen.getByTestId('light-theme-icon');
    expect(icon).toBeInTheDocument();

    // Check localStorage to confirm it now holds the 'dark' theme
    expect(localStorage.getItem('themeName')).toBe('dark');
  });
});
