import { useState, useEffect, useCallback } from 'react';
import { fetchProducts } from '../services/ProductService';
import { useApiKey } from '../auth/ApiKey';

function useProducts() {
  const { apiKey, clearApiKey } = useApiKey();
  const [loading, setLoading] = useState(true);
  const [redirectToAuth, setRedirectToAuth] = useState(false);
  const [products, setProducts] = useState([]);
  const [error, setError] = useState(null);

  const refreshProducts = useCallback(() => {
    if (!apiKey) return;

    setLoading(true);
    setError(null);
    fetchProducts(apiKey)
      .then(products => {
        setProducts(products);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
        if (err.message === 'Invalid API key') {
          clearApiKey()
          setRedirectToAuth(true);
        }
        setError(err.message);
      });
  }, [apiKey, clearApiKey]);

  useEffect(() => {
    refreshProducts();
  }, [refreshProducts])

  useEffect(() => {
    if (error) {
      console.log(error);
    }
  }, [error]);

  return {
    loading,
    error,
    redirectToAuth,
    products,
    refreshProducts,
  };
}

export default useProducts