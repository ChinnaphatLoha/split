'use client';

import React, { Suspense, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useQuery, useMutation } from '@apollo/client';
import {
  Box, Typography, Grid, Card, CardContent, Button, TextField,
  Dialog, DialogTitle, DialogContent, DialogActions, CardActionArea,
  Chip, AvatarGroup, Avatar, Stack, MenuItem, CircularProgress,
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import Providers from '../providers';
import AppLayout from '@/components/layout/AppLayout';
import { GET_MY_GROUPS, CREATE_GROUP, JOIN_BY_INVITE_CODE } from '@/lib/graphql/operations';
import { useUIStore } from '@/stores/uiStore';

const currencies = ['USD', 'EUR', 'GBP', 'JPY', 'THB', 'AUD', 'CAD', 'CHF', 'CNY', 'KRW', 'SGD', 'INR'];

function GroupsContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const showSnackbar = useUIStore((s) => s.showSnackbar);
  const [createOpen, setCreateOpen] = useState(searchParams?.get('create') === 'true');
  const [joinOpen, setJoinOpen] = useState(false);

  // Create group state
  const [newName, setNewName] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [newCurrency, setNewCurrency] = useState('USD');

  // Join group state
  const [inviteCode, setInviteCode] = useState('');

  const { data, loading, refetch } = useQuery(GET_MY_GROUPS);
  const groups = data?.myGroups || [];

  const [createGroup, { loading: creating }] = useMutation(CREATE_GROUP, {
    onCompleted: (data) => {
      showSnackbar('Group created successfully!', 'success');
      setCreateOpen(false);
      setNewName('');
      setNewDesc('');
      refetch();
      router.push(`/groups/${data.createGroup.id}`);
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [joinGroup, { loading: joining }] = useMutation(JOIN_BY_INVITE_CODE, {
    onCompleted: (data) => {
      showSnackbar('Joined group successfully!', 'success');
      setJoinOpen(false);
      setInviteCode('');
      refetch();
      router.push(`/groups/${data.joinByInviteCode.id}`);
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 4 }}>
        <Typography variant="h3">Groups</Typography>
        <Stack direction="row" spacing={1.5}>
          <Button variant="outlined" onClick={() => setJoinOpen(true)} sx={{ borderColor: 'rgba(255,255,255,0.15)' }}>
            Join Group
          </Button>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setCreateOpen(true)}>
            Create Group
          </Button>
        </Stack>
      </Box>

      <Grid container spacing={3}>
        {groups.map((group: any) => (
          <Grid item xs={12} sm={6} md={4} key={group.id}>
            <Card
              sx={{
                height: 200, cursor: 'pointer', transition: 'all 0.3s',
                '&:hover': { transform: 'translateY(-4px)', boxShadow: '0 12px 40px rgba(124,77,255,0.12)' },
              }}
              onClick={() => router.push(`/groups/${group.id}`)}
            >
              <CardContent sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
                <Box sx={{ flex: 1 }}>
                  <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1 }}>
                    <Typography variant="h5" fontWeight={700} noWrap sx={{ flex: 1 }}>{group.name}</Typography>
                    <Chip label={group.currency} size="small" sx={{ backgroundColor: 'rgba(124,77,255,0.15)', color: '#B47CFF', fontWeight: 600 }} />
                  </Stack>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2, display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                    {group.description || 'No description'}
                  </Typography>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <AvatarGroup max={4} sx={{ '& .MuiAvatar-root': { width: 28, height: 28, fontSize: '0.75rem', borderColor: '#111827' } }}>
                    {(group.members || []).map((m: any) => (
                      <Avatar key={m.userId} sx={{ bgcolor: '#7C4DFF' }}>
                        {m.user?.name?.charAt(0)?.toUpperCase() || 'U'}
                      </Avatar>
                    ))}
                  </AvatarGroup>
                  <Typography variant="caption" color="text.secondary">{group.members?.length || 0} members</Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}

        {!loading && groups.length === 0 && (
          <Grid item xs={12}>
            <Box sx={{ textAlign: 'center', py: 8, color: 'text.secondary' }}>
              <Typography variant="h5" sx={{ mb: 1 }}>No groups yet</Typography>
              <Typography>Create a group or join one with an invite code.</Typography>
            </Box>
          </Grid>
        )}
      </Grid>

      {/* Create Group Dialog */}
      <Dialog open={createOpen} onClose={() => setCreateOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ fontWeight: 700 }}>Create New Group</DialogTitle>
        <DialogContent>
          <TextField fullWidth label="Group Name" value={newName} onChange={(e) => setNewName(e.target.value)} sx={{ mt: 1, mb: 2.5 }} />
          <TextField fullWidth label="Description" multiline rows={2} value={newDesc} onChange={(e) => setNewDesc(e.target.value)} sx={{ mb: 2.5 }} />
          <TextField fullWidth select label="Currency" value={newCurrency} onChange={(e) => setNewCurrency(e.target.value)}>
            {currencies.map((c) => <MenuItem key={c} value={c}>{c}</MenuItem>)}
          </TextField>
        </DialogContent>
        <DialogActions sx={{ p: 3, pt: 1 }}>
          <Button onClick={() => setCreateOpen(false)} sx={{ color: 'text.secondary' }}>Cancel</Button>
          <Button
            variant="contained"
            disabled={!newName || creating}
            onClick={() => createGroup({ variables: { name: newName, description: newDesc, currency: newCurrency } })}
          >
            {creating ? <CircularProgress size={20} /> : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Join Group Dialog */}
      <Dialog open={joinOpen} onClose={() => setJoinOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ fontWeight: 700 }}>Join Group</DialogTitle>
        <DialogContent>
          <Typography color="text.secondary" sx={{ mb: 2 }}>Enter the invite code to join an existing group.</Typography>
          <TextField fullWidth label="Invite Code" value={inviteCode} onChange={(e) => setInviteCode(e.target.value)} placeholder="e.g. a1b2c3d4" sx={{ mt: 1 }} />
        </DialogContent>
        <DialogActions sx={{ p: 3, pt: 1 }}>
          <Button onClick={() => setJoinOpen(false)} sx={{ color: 'text.secondary' }}>Cancel</Button>
          <Button
            variant="contained"
            disabled={!inviteCode || joining}
            onClick={() => joinGroup({ variables: { inviteCode } })}
          >
            {joining ? <CircularProgress size={20} /> : 'Join'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default function GroupsPage() {
  return (
    <Providers>
      <AppLayout>
        <Suspense fallback={null}>
          <GroupsContent />
        </Suspense>
      </AppLayout>
    </Providers>
  );
}
