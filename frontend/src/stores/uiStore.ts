'use client';

import { create } from 'zustand';

interface UIState {
  snackbar: { open: boolean; message: string; severity: 'success' | 'error' | 'info' | 'warning' };
  showSnackbar: (message: string, severity?: 'success' | 'error' | 'info' | 'warning') => void;
  hideSnackbar: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  snackbar: { open: false, message: '', severity: 'info' },

  showSnackbar: (message, severity = 'info') => {
    set({ snackbar: { open: true, message, severity } });
  },

  hideSnackbar: () => {
    set((state) => ({ snackbar: { ...state.snackbar, open: false } }));
  },
}));
