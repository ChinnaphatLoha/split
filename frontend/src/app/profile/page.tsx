'use client';

import React, { useState, useEffect } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import {
  Box, Typography, Card, CardContent, TextField, Button, Grid,
  Avatar, MenuItem, CircularProgress, Divider,
} from '@mui/material';
import { Save as SaveIcon, Person as PersonIcon } from '@mui/icons-material';
import Providers from '../providers';
import AppLayout from '@/components/layout/AppLayout';
import { GET_ME, UPDATE_PROFILE } from '@/lib/graphql/operations';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';

const currencies = ['USD', 'EUR', 'GBP', 'JPY', 'THB', 'AUD', 'CAD', 'CHF', 'CNY', 'KRW', 'SGD', 'INR'];

function ProfileContent() {
  const authUser = useAuthStore((s) => s.user);
  const updateAuthUser = useAuthStore((s) => s.updateUser);
  const showSnackbar = useUIStore((s) => s.showSnackbar);

  const [name, setName] = useState(authUser?.name || '');
  const [currency, setCurrency] = useState(authUser?.currency || 'USD');
  const [avatarUrl, setAvatarUrl] = useState(authUser?.avatarUrl || '');

  const { data } = useQuery(GET_ME);

  useEffect(() => {
    if (data?.me) {
      setName(data.me.name);
      setCurrency(data.me.currency);
      setAvatarUrl(data.me.avatarUrl || '');
    }
  }, [data]);

  const [updateProfile, { loading }] = useMutation(UPDATE_PROFILE, {
    onCompleted: (data) => {
      const u = data.updateProfile;
      updateAuthUser({ name: u.name, currency: u.currency, avatarUrl: u.avatarUrl });
      showSnackbar('Profile updated successfully!', 'success');
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const handleSave = () => {
    updateProfile({ variables: { name, currency, avatarUrl } });
  };

  const userInfo = data?.me || authUser;

  return (
    <Box>
      <Typography variant="h3" sx={{ mb: 4 }}>Profile</Typography>

      <Grid container spacing={3}>
        {/* Avatar Card */}
        <Grid item xs={12} md={4}>
          <Card sx={{ textAlign: 'center', p: 4 }}>
            <Avatar
              sx={{
                width: 100,
                height: 100,
                mx: 'auto',
                mb: 2,
                bgcolor: '#7C4DFF',
                fontSize: '2.5rem',
                fontWeight: 800,
              }}
            >
              {userInfo?.name?.charAt(0)?.toUpperCase() || 'U'}
            </Avatar>
            <Typography variant="h5" fontWeight={700}>{userInfo?.name}</Typography>
            <Typography color="text.secondary" sx={{ mt: 0.5 }}>{userInfo?.email}</Typography>
            <Divider sx={{ my: 2, borderColor: 'rgba(255,255,255,0.06)' }} />
            <Box sx={{ display: 'flex', justifyContent: 'center', gap: 3 }}>
              <Box>
                <Typography variant="h6" fontWeight={700}>{userInfo?.currency}</Typography>
                <Typography variant="caption" color="text.secondary">Currency</Typography>
              </Box>
              <Box>
                <Typography variant="h6" fontWeight={700}>
                  {userInfo?.createdAt ? new Date(userInfo.createdAt).getFullYear() : '-'}
                </Typography>
                <Typography variant="caption" color="text.secondary">Joined</Typography>
              </Box>
            </Box>
          </Card>
        </Grid>

        {/* Edit Form */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent sx={{ p: 3.5 }}>
              <Typography variant="h5" fontWeight={700} sx={{ mb: 3 }}>Edit Profile</Typography>
              <TextField
                id="profile-name"
                fullWidth label="Full Name" value={name}
                onChange={(e) => setName(e.target.value)}
                sx={{ mb: 2.5 }}
              />
              <TextField
                id="profile-email"
                fullWidth label="Email" value={userInfo?.email || ''}
                disabled
                sx={{ mb: 2.5 }}
                helperText="Email cannot be changed"
              />
              <TextField
                id="profile-avatar"
                fullWidth label="Avatar URL" value={avatarUrl}
                onChange={(e) => setAvatarUrl(e.target.value)}
                sx={{ mb: 2.5 }}
                placeholder="https://example.com/avatar.jpg"
              />
              <TextField
                id="profile-currency"
                fullWidth select label="Default Currency" value={currency}
                onChange={(e) => setCurrency(e.target.value)}
                sx={{ mb: 3 }}
              >
                {currencies.map((c) => <MenuItem key={c} value={c}>{c}</MenuItem>)}
              </TextField>
              <Button
                id="profile-save"
                variant="contained"
                startIcon={loading ? <CircularProgress size={18} /> : <SaveIcon />}
                onClick={handleSave}
                disabled={loading}
                sx={{ px: 4 }}
              >
                Save Changes
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}

export default function ProfilePage() {
  return (
    <Providers>
      <AppLayout>
        <ProfileContent />
      </AppLayout>
    </Providers>
  );
}
