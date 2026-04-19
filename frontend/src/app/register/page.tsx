'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation } from '@apollo/client';
import {
  Box, Card, CardContent, Typography, TextField, Button,
  Link, CircularProgress, Alert, InputAdornment, IconButton,
  MenuItem,
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import Providers from '../providers';
import { REGISTER } from '@/lib/graphql/operations';
import { useAuthStore } from '@/stores/authStore';

const currencies = ['USD', 'EUR', 'GBP', 'JPY', 'THB', 'AUD', 'CAD', 'CHF', 'CNY', 'KRW', 'SGD', 'INR'];

function RegisterContent() {
  const router = useRouter();
  const login = useAuthStore((s) => s.login);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');

  const [registerMutation, { loading }] = useMutation(REGISTER, {
    onCompleted: (data) => {
      const { user, accessToken, refreshToken } = data.register;
      login(user, accessToken, refreshToken);
      router.push('/dashboard');
    },
    onError: (err) => {
      setError(err.message || 'Registration failed');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!name || !email || !password) {
      setError('Please fill in all fields');
      return;
    }
    if (password.length < 6) {
      setError('Password must be at least 6 characters');
      return;
    }
    registerMutation({ variables: { email, name, password, currency } });
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
          background: 'radial-gradient(circle, rgba(0,229,255,0.1) 0%, transparent 70%)',
          filter: 'blur(60px)',
        },
      }}
    >
      <Card sx={{ width: '100%', maxWidth: 440, position: 'relative', zIndex: 1 }}>
        <CardContent sx={{ p: { xs: 3, sm: 4.5 } }}>
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
            <Typography variant="h4" fontWeight={800}>Create account</Typography>
            <Typography color="text.secondary" sx={{ mt: 0.5 }}>Start splitting expenses with ease</Typography>
          </Box>

          {error && <Alert severity="error" sx={{ mb: 3, borderRadius: 2 }}>{error}</Alert>}

          <form onSubmit={handleSubmit}>
            <TextField
              id="register-name"
              fullWidth label="Full Name" value={name}
              onChange={(e) => setName(e.target.value)}
              sx={{ mb: 2.5 }}
              autoComplete="name"
            />
            <TextField
              id="register-email"
              fullWidth label="Email" type="email" value={email}
              onChange={(e) => setEmail(e.target.value)}
              sx={{ mb: 2.5 }}
              autoComplete="email"
            />
            <TextField
              id="register-password"
              fullWidth label="Password" type={showPassword ? 'text' : 'password'} value={password}
              onChange={(e) => setPassword(e.target.value)}
              sx={{ mb: 2.5 }}
              autoComplete="new-password"
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
            <TextField
              id="register-currency"
              fullWidth select label="Default Currency" value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              sx={{ mb: 3 }}
            >
              {currencies.map((c) => (
                <MenuItem key={c} value={c}>{c}</MenuItem>
              ))}
            </TextField>
            <Button
              id="register-submit"
              type="submit" variant="contained" fullWidth size="large"
              disabled={loading}
              sx={{ py: 1.5, mb: 2.5 }}
            >
              {loading ? <CircularProgress size={24} color="inherit" /> : 'Create Account'}
            </Button>
          </form>

          <Typography variant="body2" color="text.secondary" textAlign="center">
            Already have an account?{' '}
            <Link
              component="button"
              onClick={() => router.push('/login')}
              sx={{ color: '#7C4DFF', fontWeight: 600, textDecoration: 'none', '&:hover': { textDecoration: 'underline' } }}
            >
              Sign in
            </Link>
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}

export default function RegisterPage() {
  return (
    <Providers>
      <RegisterContent />
    </Providers>
  );
}
