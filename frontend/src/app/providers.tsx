'use client';

import React, { useEffect } from 'react';
import { ThemeProvider, CssBaseline, Snackbar, Alert } from '@mui/material';
import { ApolloProvider } from '@apollo/client';
import theme from '@/theme/theme';
import { apolloClient } from '@/lib/graphql/client';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';

export default function Providers({ children }: { children: React.ReactNode }) {
  const loadFromStorage = useAuthStore((s) => s.loadFromStorage);
  const snackbar = useUIStore((s) => s.snackbar);
  const hideSnackbar = useUIStore((s) => s.hideSnackbar);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  return (
    <ApolloProvider client={apolloClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {children}
        <Snackbar
          open={snackbar.open}
          autoHideDuration={4000}
          onClose={hideSnackbar}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        >
          <Alert onClose={hideSnackbar} severity={snackbar.severity} variant="filled" sx={{ borderRadius: 2 }}>
            {snackbar.message}
          </Alert>
        </Snackbar>
      </ThemeProvider>
    </ApolloProvider>
  );
}
