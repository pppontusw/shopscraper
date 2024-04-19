import React from 'react';
import AppRoutes from './routes/AppRoutes';
import { CustomThemeProvider } from './theme/useTheme';

function App() {
  return (
    <CustomThemeProvider>
      <AppRoutes />
    </CustomThemeProvider>
  )
}

export default App;
