import React, { useState, useContext } from 'react';

export const ApiKeyContext = React.createContext();

export const useApiKey = () => useContext(ApiKeyContext);

export const ApiKeyProvider = ({ children }) => {
  const [apiKey, setApiKey] = useState(localStorage.getItem('apiKey') || '');
  const [error, setError] = useState('');

  const saveApiKey = (key) => {
    localStorage.setItem('apiKey', key);
    setApiKey(key);
    setError(''); // Clear error when key is updated
  };

  const clearApiKey = () => {
    localStorage.removeItem('apiKey');
    setApiKey('');
    setError('Invalid API key'); // Set error when API key is cleared due to being invalid
  };

  return (
    <ApiKeyContext.Provider value={{ apiKey, saveApiKey, error, clearApiKey }}>
      {children}
    </ApiKeyContext.Provider>
  );
};