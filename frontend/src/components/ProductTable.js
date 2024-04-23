import React, { useEffect, useState } from 'react';
import { Box, Link, useTheme, CircularProgress } from '@mui/material';
import { DataGrid, useGridApiRef } from '@mui/x-data-grid';
import timeAgo from '../utils/timeAgo';
import CustomToolbar from './CustomToolbar';

const getApplyFilterFn = (value) => {
  if (!value) {
    return null;
  }

  const words = value.split(' ') || [];
  // Filter out empty quotes
  const filteredWords = words.filter(word => word !== '""');

  return (cellValue) => {
    // Split the cell value into individual words
    const cellWords = cellValue.split(' ');

    return filteredWords.every(word => {
      if (word.startsWith('"') && word.endsWith('"')) {
        // Exact word matching for quoted phrases
        const exactWord = word.slice(1, -1);
        return cellWords.some(cellWord => cellWord.toLowerCase() === exactWord.toLowerCase());
      } else {

        // Normalize commas to periods for matching interchangeably
        const normalizedWord = word.replace(/,/g, '.');
        return cellWords.some(cellWord => {
          const normalizedCellWord = cellWord.replace(/,/g, '.');
          return normalizedCellWord.toLowerCase().includes(normalizedWord.toLowerCase());
        });
      }
    });
  };
};


const ProductTable = ({
  loading,
  products,
  onRefresh
}) => {
  const theme = useTheme();
  const apiRef = useGridApiRef();
  const [initialState, setInitialState] = useState(null);

  const columns = [
    {
      field: 'name',
      headerName: 'Name',
      minWidth: 250,
      flex: 1,
      renderCell: (cellValues) => (
        <Link href={cellValues.row.link} target="_blank" rel="noopener noreferrer">
          {cellValues.value}
        </Link>
      ),
      getApplyQuickFilterFn: getApplyFilterFn,
    },
    {
      field: 'previousPrice',
      headerName: 'Previous Price',
      type: 'number',
      minWidth: 120,
      valueGetter: (params) => params.Valid ? params.Int64 : null
    },
    { field: 'price', headerName: 'Price', type: 'number', minWidth: 120 },
    { field: 'shop', headerName: 'Shop', minWidth: 150 },
    {
      field: 'lastSeen',
      headerName: 'Last Seen',
      minWidth: 150,
      valueGetter: (params) => new Date(params),
      renderCell: (cellValues) => (timeAgo(cellValues.value)),
    },
    {
      field: 'firstSeen',
      headerName: 'First Seen',
      minWidth: 150,
      valueGetter: (params) => new Date(params),
      renderCell: (cellValues) => (timeAgo(cellValues.value)),
    },
    { field: 'notified', headerName: 'Notified', minWidth: 50 },
  ];

  const rows = (products || []).map((product, index) => ({
    id: product.price + product.link + product.name,
    name: product.name,
    previousPrice: product.previousPrice,
    price: product.price,
    shop: product.shop,
    lastSeen: product.lastSeen,
    firstSeen: product.firstSeen,
    link: product.link,
    notified: product.notified,
  }));

  useEffect(() => {
    const stateFromLocalStorage = localStorage.getItem('productGridState');
    if (stateFromLocalStorage) {
      setInitialState(JSON.parse(stateFromLocalStorage));
    } else {
      setInitialState({
        columns: {
          columnVisibilityModel: {
            previousPrice: false,
            notified: false
          }
        }
      });
    }

    const saveSnapshot = () => {
      if (apiRef?.current?.exportState) {
        const currentState = apiRef.current.exportState();
        localStorage.setItem('productGridState', JSON.stringify(currentState));
      }
    };

    window.addEventListener('beforeunload', saveSnapshot);

    return () => {
      window.removeEventListener('beforeunload', saveSnapshot);
      saveSnapshot();
    };
  }, [apiRef]);

  const getRowClassName = (params) => {
    const isDarkMode = theme.palette.mode === 'dark';
    return params.indexRelativeToCurrentPage % 2 === 0 ? (isDarkMode ? 'darkStripe' : 'lightStripe') : '';
  };

  if (!initialState) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" style={{ height: "100vh" }}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <div style={{ width: '100%' }}>
        <DataGrid
          apiRef={apiRef}
          initialState={initialState}
          autosizeOptions={{
            includeOutliers: true,
            includeHeaders: false,
          }}
          autoHeight
          rows={rows}
          columns={columns}
          pageSize={5}
          rowsPerPageOptions={[5]}
          disableSelectionOnClick
          isRowSelectable={() => false}
          slotProps={{
            toolbar: {
              csvOptions: {
                allColumns: true
              },
              printOptions: {
                disableToolbarButton: true
              }
            }
          }
          }
          slots={{ toolbar: (props) => <CustomToolbar {...props} onRefresh={onRefresh} /> }}
          loading={loading}
          ignoreDiacritics
          getRowClassName={getRowClassName}
          sx={{
            '& .MuiDataGrid-row': {
              // Default background color for non-striped rows
              backgroundColor: 'inherit',
            },
            '& .MuiDataGrid-row:nth-of-type(odd)': {
              // Striped row background color based on theme
              backgroundColor: theme.palette.mode === 'dark' ? theme.palette.grey[900] : theme.palette.action.hover,
            },
            '& .MuiDataGrid-row:hover': {
              // Match hover state to the normal state to "disable" hover effect
              backgroundColor: 'inherit !important',
            },
            '& .MuiDataGrid-row:nth-of-type(odd):hover': {
              // Ensure odd rows maintain striped color on hover
              backgroundColor: `${theme.palette.mode === 'dark' ? theme.palette.grey[900] : theme.palette.action.hover} !important`,
            },
          }}
        />
      </div>
    </Box>
  );
};

export default ProductTable;
