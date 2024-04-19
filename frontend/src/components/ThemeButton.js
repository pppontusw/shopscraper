import React from 'react';
import { useTheme } from '../theme/useTheme';
import Button from '@mui/material/Button';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';

const ThemeToggleButton = () => {
  const { themeName, toggleTheme } = useTheme();

  const icon = themeName === "light" ? <DarkModeIcon data-testid="dark-theme-icon" /> : <LightModeIcon data-testid="light-theme-icon" />;
  const text = themeName === "light" ? "Dark" : "Light"

  return (
    <Button startIcon={icon} size="small" onClick={toggleTheme}>
      { text }
    </Button>
  );
};

export default ThemeToggleButton;
