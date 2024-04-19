import React, { useState } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { useApiKey } from '../auth/ApiKey';


export const ApiKeyInput = () => {

  const handleSubmit = (event) => {
    event.preventDefault();
    saveApiKey(tempKey);
  };

  const { saveApiKey, error } = useApiKey();
  const [tempKey, setTempKey] = useState('');

  return (
    <Box
      component="form"
      sx={{
        '& .MuiTextField-root': { m: 1, width: '25ch' },
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
      }}
      noValidate
      autoComplete="off"
      onSubmit={handleSubmit}
    >
      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}
      <TextField
        label="API Key"
        variant="outlined"
        value={tempKey}
        onChange={(e) => setTempKey(e.target.value)}
        required
      />
      <Button type="submit" variant="contained" color="primary" sx={{ mt: 2 }}>
        Submit
      </Button>
    </Box>
  );
};
