import React from 'react';
import { GridToolbarDensitySelector, GridToolbarQuickFilter, GridToolbarContainer, GridToolbarExport } from '@mui/x-data-grid';
import { Box, Button } from '@mui/material';
import ThemeToggleButton from './ThemeButton';
import RefreshIcon from '@mui/icons-material/Refresh';

const CustomToolbar = ({ onRefresh }) => {
  return (
    <GridToolbarContainer >
      <GridToolbarDensitySelector />
      <GridToolbarExport />
      <Button size="small" onClick={onRefresh} startIcon={<RefreshIcon />} sx={{ p: 0 }}>
        Refresh
      </Button>
      <Box sx={{ flexGrow: 1 }} />
      <GridToolbarQuickFilter sx={{ paddingRight: 2 }} />
      <ThemeToggleButton />
    </GridToolbarContainer>
  );
};

export default CustomToolbar