import React from 'react';
import { Navigate } from 'react-router-dom';
import ProductTable from './ProductTable';
import useProducts from '../hooks/useProducts';
import { Typography, Box, Button, useTheme } from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

const ProductsView = () => {
  const {
    loading,
    error,
    redirectToAuth,
    products,
    refreshProducts,
  } = useProducts();

  const theme = useTheme();

  if (redirectToAuth) {
    return <Navigate to={"/auth"} replace />;
  }

  return (
    <div>
      {error && (
        <Box 
          display="flex" 
          flexDirection="column"
          justifyContent="center" 
          alignItems="center" 
          style={{ height: "100vh", padding: 20, backgroundColor: theme.palette.background.default }}
        >
          <ErrorOutlineIcon style={{ fontSize: 60, color: theme.palette.error.main }} />
          <Typography variant="h5" color="error" style={{ margin: '20px 0' }}>
            Oops! Something went wrong.
          </Typography>
          <Typography variant="subtitle1" color="error" style={{ marginBottom: 20 }}>
            {error}
        </Typography>
        <br />
          <Button variant="outlined" color="primary" onClick={refreshProducts}>
            Try Again
          </Button>
        </Box>
      )}
      {!error && (
        <ProductTable
          loading={loading}
          products={products}
          onRefresh={refreshProducts}
        />
      )}
    </div>
  );
};

export default ProductsView;
