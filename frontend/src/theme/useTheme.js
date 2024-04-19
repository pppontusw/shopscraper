import React, { createContext, useContext, useState, useMemo, useEffect } from 'react';
import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { lightTheme, darkTheme } from './theme';

const ThemeToggleContext = createContext();

export const useTheme = () => useContext(ThemeToggleContext);

export const CustomThemeProvider = ({ children }) => {
  // Initialize themeName state from localStorage or default to 'dark'
  const [themeName, setThemeName] = useState(() => {
    const storedThemeName = localStorage.getItem('themeName');
    return storedThemeName || 'dark';
  });

  const theme = useMemo(() => {
    return themeName === 'light' ? lightTheme : darkTheme;
  }, [themeName]);

  useEffect(() => {
    // Update localStorage whenever themeName changes
    localStorage.setItem('themeName', themeName);
  }, [themeName]);

  useEffect(() => {
    document.body.style.backgroundColor = theme.palette.background.default;
  }, [theme]);

  const toggleTheme = () => {
    setThemeName((prevThemeName) => {
      const newThemeName = prevThemeName === 'light' ? 'dark' : 'light';
      return newThemeName;
    });
  };

  return (
    <ThemeToggleContext.Provider value={{ themeName, toggleTheme }}>
      <MuiThemeProvider theme={theme}>{children}</MuiThemeProvider>
    </ThemeToggleContext.Provider>
  );
};
