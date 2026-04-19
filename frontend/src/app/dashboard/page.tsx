'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@apollo/client';
import {
  Box, Typography, Grid, Card, CardContent, CardActionArea,
  Skeleton, Avatar, AvatarGroup, Chip, Stack,
} from '@mui/material';
import {
  Add as AddIcon,
  Group as GroupIcon,
  TrendingUp as TrendingIcon,
  AccountBalanceWallet as WalletIcon,
} from '@mui/icons-material';
import Providers from '../providers';
import AppLayout from '@/components/layout/AppLayout';
import { GET_MY_GROUPS } from '@/lib/graphql/operations';
import { useAuthStore } from '@/stores/authStore';

function DashboardContent() {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const { data, loading } = useQuery(GET_MY_GROUPS);

  const groups = data?.myGroups || [];

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return 'Good morning';
    if (hour < 18) return 'Good afternoon';
    return 'Good evening';
  };

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h3" sx={{ mb: 0.5 }}>
          {getGreeting()}, {user?.name?.split(' ')[0] || 'there'} 👋
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Here's an overview of your groups and expenses.
        </Typography>
      </Box>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {[
          { label: 'Total Groups', value: groups.length, icon: <GroupIcon />, color: '#7C4DFF' },
          { label: 'Active Members', value: groups.reduce((acc: number, g: any) => acc + (g.members?.length || 0), 0), icon: <TrendingIcon />, color: '#00E5FF' },
          { label: 'Currency', value: user?.currency || 'USD', icon: <WalletIcon />, color: '#10B981' },
        ].map((stat, i) => (
          <Grid item xs={12} sm={4} key={i}>
            <Card sx={{ transition: 'transform 0.2s', '&:hover': { transform: 'translateY(-2px)' } }}>
              <CardContent sx={{ display: 'flex', alignItems: 'center', gap: 2, py: 3 }}>
                <Box sx={{ width: 48, height: 48, borderRadius: 2, backgroundColor: `${stat.color}20`, display: 'flex', alignItems: 'center', justifyContent: 'center', color: stat.color }}>
                  {stat.icon}
                </Box>
                <Box>
                  <Typography variant="caption" color="text.secondary" textTransform="uppercase" letterSpacing={1}>
                    {stat.label}
                  </Typography>
                  <Typography variant="h4" fontWeight={700}>{stat.value}</Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Groups */}
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Your Groups</Typography>
      </Box>

      {loading ? (
        <Grid container spacing={3}>
          {[1, 2, 3].map((i) => (
            <Grid item xs={12} sm={6} md={4} key={i}>
              <Skeleton variant="rounded" height={200} sx={{ borderRadius: 4 }} />
            </Grid>
          ))}
        </Grid>
      ) : (
        <Grid container spacing={3}>
          {/* Create Group Card */}
          <Grid item xs={12} sm={6} md={4}>
            <Card
              sx={{
                height: 200,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                border: '2px dashed rgba(124,77,255,0.3)',
                backgroundColor: 'transparent',
                cursor: 'pointer',
                transition: 'all 0.3s ease',
                '&:hover': {
                  borderColor: '#7C4DFF',
                  backgroundColor: 'rgba(124,77,255,0.05)',
                },
              }}
              onClick={() => router.push('/groups?create=true')}
            >
              <CardActionArea sx={{ height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                <AddIcon sx={{ fontSize: 48, color: '#7C4DFF', mb: 1 }} />
                <Typography variant="h6" color="primary">Create New Group</Typography>
              </CardActionArea>
            </Card>
          </Grid>

          {/* Group Cards */}
          {groups.map((group: any) => (
            <Grid item xs={12} sm={6} md={4} key={group.id}>
              <Card
                sx={{
                  height: 200,
                  cursor: 'pointer',
                  transition: 'all 0.3s ease',
                  '&:hover': { transform: 'translateY(-4px)', boxShadow: '0 12px 40px rgba(124,77,255,0.12)' },
                }}
                onClick={() => router.push(`/groups/${group.id}`)}
              >
                <CardContent sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
                  <Box sx={{ flex: 1 }}>
                    <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1 }}>
                      <Typography variant="h5" fontWeight={700} noWrap sx={{ flex: 1 }}>
                        {group.name}
                      </Typography>
                      <Chip label={group.currency} size="small" sx={{ backgroundColor: 'rgba(124,77,255,0.15)', color: '#B47CFF', fontWeight: 600 }} />
                    </Stack>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2, display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                      {group.description || 'No description'}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <AvatarGroup max={4} sx={{ '& .MuiAvatar-root': { width: 30, height: 30, fontSize: '0.8rem', borderColor: '#111827' } }}>
                      {(group.members || []).map((m: any) => (
                        <Avatar key={m.userId} sx={{ bgcolor: '#7C4DFF' }}>
                          {m.user?.name?.charAt(0)?.toUpperCase() || 'U'}
                        </Avatar>
                      ))}
                    </AvatarGroup>
                    <Typography variant="caption" color="text.secondary">
                      {group.members?.length || 0} members
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}

          {groups.length === 0 && (
            <Grid item xs={12} sm={6} md={8}>
              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: 200, color: 'text.secondary' }}>
                <Typography>No groups yet. Create one to get started!</Typography>
              </Box>
            </Grid>
          )}
        </Grid>
      )}
    </Box>
  );
}

export default function DashboardPage() {
  return (
    <Providers>
      <AppLayout>
        <DashboardContent />
      </AppLayout>
    </Providers>
  );
}
