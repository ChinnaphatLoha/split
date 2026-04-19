'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation } from '@apollo/client';
import {
  Box, Card, CardContent, Typography, TextField, Button,
  Link, CircularProgress, Alert, InputAdornment, IconButton,
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import Providers from '../providers';
import { LOGIN } from '@/lib/graphql/operations';
import { useAuthStore } from '@/stores/authStore';

function LoginContent() {
  const router = useRouter();
  const login = useAuthStore((s) => s.login);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');

  const [loginMutation, { loading }] = useMutation(LOGIN, {
    onCompleted: (data) => {
      const { user, accessToken, refreshToken } = data.login;
      login(user, accessToken, refreshToken);
      router.push('/dashboard');
    },
    onError: (err) => {
      setError(err.message || 'Invalid email or password');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!email || !password) {
      setError('Please fill in all fields');
      return;
    }
    loginMutation({ variables: { email, password } });
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        p: 2,
        position: 'relative',
        '&::before': {
          content: '""',
          position: 'absolute',
          top: '20%',
          left: '50%',
          transform: 'translateX(-50%)',
          width: 500,
          height: 500,
          borderRadius: '50%',
          background: 'radial-gradient(circle, rgba(124,77,255,0.12) 0%, transparent 70%)',
          filter: 'blur(60px)',
        },
      }}
    >
      <Card sx={{ width: '100%', maxWidth: 440, position: 'relative', zIndex: 1 }}>
        <CardContent sx={{ p: { xs: 3, sm: 4.5 } }}>
          {/* Logo */}
          <Box sx={{ textAlign: 'center', mb: 4 }}>
            <Box
              sx={{
                width: 52, height: 52, borderRadius: 2, mx: 'auto', mb: 2,
                background: 'linear-gradient(135deg, #7C4DFF 0%, #00E5FF 100%)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                fontWeight: 900, fontSize: '1.5rem', color: '#fff',
              }}
            >
              S
            </Box>
            <Typography variant="h4" fontWeight={800}>Welcome back</Typography>
            <Typography color="text.secondary" sx={{ mt: 0.5 }}>Sign in to your Split account</Typography>
          </Box>

          {error && <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>{error}</Alert>}

          <form onSubmit={handleSubmit}>
            <TextField
              id="login-email"
              fullWidth
              label="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              sx={{ mb: 2.5 }}
              autoComplete="email"
            />
            <TextField
              id="login-password"
              fullWidth
              label="Password"
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              sx={{ mb: 3 }}
              autoComplete="current-password"
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={() => setShowPassword(!showPassword)} edge="end" size="small">
                      {showPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            <Button
              id="login-submit"
              type="submit"
              variant="contained"
              fullWidth
              size="large"
              disabled={loading}
              sx={{ py: 1.5, mb: 2.5 }}
            >
              {loading ? <CircularProgress size={24} color="inherit" /> : 'Sign In'}
            </Button>
          </form>

          <Typography variant="body2" color="text.secondary" textAlign="center">
            Don't have an account?{' '}
            <Link
              component="button"
              onClick={() => router.push('/register')}
              sx={{ color: '#7C4DFF', fontWeight: 600, textDecoration: 'none', '&:hover': { textDecoration: 'underline' } }}
            >
              Create one
            </Link>
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}

export default function LoginPage() {
  return (
    <Providers>
      <LoginContent />
    </Providers>
  );
}
