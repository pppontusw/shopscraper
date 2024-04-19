import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import AuthRoute from '../auth/AuthRoute';
import PrivateRoute from '../auth/PrivateRoute';
import { ApiKeyProvider } from '../auth/ApiKey';
import ProductsView from '../components/ProductsView'; // Import ProductsView

function AppRoutes() {
  return (
    <ApiKeyProvider>
      <Router>
        <Routes>
          <Route path="/auth" element={<AuthRoute />} />
          <Route path="/" element={<PrivateRoute><ProductsView /></PrivateRoute>} />
        </Routes>
      </Router>
    </ApiKeyProvider>
  );
}

export default AppRoutes;
