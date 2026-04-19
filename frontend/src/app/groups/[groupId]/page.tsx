'use client';

import React, { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useQuery, useMutation } from '@apollo/client';
import {
  Box, Typography, Tabs, Tab, Card, CardContent, Grid, Avatar,
  Button, Chip, Stack, IconButton, Tooltip, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, MenuItem, List, ListItem,
  ListItemAvatar, ListItemText, ListItemSecondaryAction, Divider,
  CircularProgress, Skeleton, Table, TableBody, TableCell, TableContainer,
  TableHead, TableRow, Paper,
} from '@mui/material';
import {
  Add as AddIcon,
  ContentCopy as CopyIcon,
  Delete as DeleteIcon,
  ArrowForward as ArrowIcon,
  Receipt as ReceiptIcon,
  AccountBalance as SettleIcon,
  People as PeopleIcon,
  Check as CheckIcon,
} from '@mui/icons-material';
import Providers from '../../providers';
import AppLayout from '@/components/layout/AppLayout';
import {
  GET_GROUP, GET_GROUP_EXPENSES, GET_GROUP_BALANCES, GET_SETTLEMENT_PLAN,
  GET_SETTLEMENTS, CREATE_EXPENSE, DELETE_EXPENSE, RECORD_SETTLEMENT,
  MARK_SETTLEMENT_COMPLETED, GENERATE_INVITE_CODE, REMOVE_MEMBER,
} from '@/lib/graphql/operations';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';

const splitTypes = [
  { value: 'EQUAL', label: 'Equal Split' },
  { value: 'CUSTOM', label: 'Custom Amounts' },
  { value: 'PERCENTAGE', label: 'Percentage' },
];

function GroupDetailContent() {
  const params = useParams();
  const router = useRouter();
  const groupId = params?.groupId as string;
  const user = useAuthStore((s) => s.user);
  const showSnackbar = useUIStore((s) => s.showSnackbar);
  const [tab, setTab] = useState(0);
  const [expenseDialogOpen, setExpenseDialogOpen] = useState(false);

  // Expense form state
  const [expAmount, setExpAmount] = useState('');
  const [expDesc, setExpDesc] = useState('');
  const [expSplitType, setExpSplitType] = useState('EQUAL');
  const [expSplits, setExpSplits] = useState<{ userId: string; amount: number; percentage: number }[]>([]);

  // Queries
  const { data: groupData, loading: groupLoading } = useQuery(GET_GROUP, { variables: { id: groupId }, skip: !groupId });
  const { data: expensesData, loading: expensesLoading, refetch: refetchExpenses } = useQuery(GET_GROUP_EXPENSES, { variables: { groupId, limit: 50, offset: 0 }, skip: !groupId });
  const { data: balancesData, refetch: refetchBalances } = useQuery(GET_GROUP_BALANCES, { variables: { groupId }, skip: !groupId });
  const { data: planData, refetch: refetchPlan } = useQuery(GET_SETTLEMENT_PLAN, { variables: { groupId }, skip: !groupId });
  const { data: settlementsData, refetch: refetchSettlements } = useQuery(GET_SETTLEMENTS, { variables: { groupId }, skip: !groupId });

  const group = groupData?.group;
  const expenses = expensesData?.groupExpenses?.expenses || [];
  const balances = balancesData?.groupBalances?.balances || [];
  const balanceDetails = balancesData?.groupBalances?.details || [];
  const transactions = planData?.settlementPlan?.transactions || [];
  const settlements = settlementsData?.settlements || [];

  const isOwner = group?.ownerId === user?.id;

  // Mutations
  const [createExpense, { loading: creatingExpense }] = useMutation(CREATE_EXPENSE, {
    onCompleted: () => {
      showSnackbar('Expense added!', 'success');
      setExpenseDialogOpen(false);
      resetExpenseForm();
      refetchExpenses();
      refetchBalances();
      refetchPlan();
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [deleteExpense] = useMutation(DELETE_EXPENSE, {
    onCompleted: () => {
      showSnackbar('Expense deleted', 'success');
      refetchExpenses();
      refetchBalances();
      refetchPlan();
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [recordSettlement] = useMutation(RECORD_SETTLEMENT, {
    onCompleted: () => {
      showSnackbar('Settlement recorded!', 'success');
      refetchSettlements();
      refetchBalances();
      refetchPlan();
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [markCompleted] = useMutation(MARK_SETTLEMENT_COMPLETED, {
    onCompleted: () => {
      showSnackbar('Settlement marked as completed!', 'success');
      refetchSettlements();
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [generateInvite] = useMutation(GENERATE_INVITE_CODE, {
    onCompleted: (data) => {
      navigator.clipboard.writeText(data.generateInviteCode);
      showSnackbar('New invite code copied to clipboard!', 'success');
    },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const [removeMember] = useMutation(REMOVE_MEMBER, {
    onCompleted: () => { showSnackbar('Member removed', 'success'); },
    onError: (err) => showSnackbar(err.message, 'error'),
  });

  const resetExpenseForm = () => {
    setExpAmount('');
    setExpDesc('');
    setExpSplitType('EQUAL');
    setExpSplits([]);
  };

  const handleOpenExpenseDialog = () => {
    // Pre-populate splits with all group members
    if (group?.members) {
      setExpSplits(group.members.map((m: any) => ({ userId: m.userId, amount: 0, percentage: 0 })));
    }
    setExpenseDialogOpen(true);
  };

  const handleCreateExpense = () => {
    const amount = parseFloat(expAmount);
    if (!amount || !expDesc) return;

    let splits = expSplits;
    if (expSplitType === 'EQUAL') {
      const perPerson = amount / splits.length;
      splits = splits.map((s) => ({ ...s, amount: perPerson, percentage: 100 / splits.length }));
    }

    createExpense({
      variables: {
        groupId,
        amount,
        description: expDesc,
        splitType: expSplitType,
        splits: splits.map((s) => ({
          userId: s.userId,
          amount: s.amount || 0,
          percentage: s.percentage || 0,
        })),
      },
    });
  };

  const copyInviteCode = () => {
    if (group?.inviteCode) {
      navigator.clipboard.writeText(group.inviteCode);
      showSnackbar('Invite code copied!', 'success');
    }
  };

  if (groupLoading) {
    return (
      <Box>
        <Skeleton variant="text" width={300} height={50} />
        <Skeleton variant="rounded" height={200} sx={{ mt: 2, borderRadius: 4 }} />
      </Box>
    );
  }

  if (!group) {
    return <Typography>Group not found</Typography>;
  }

  return (
    <Box>
      {/* Group Header */}
      <Box sx={{ mb: 4 }}>
        <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1 }}>
          <Typography variant="h3">{group.name}</Typography>
          <Stack direction="row" spacing={1}>
            <Tooltip title="Copy invite code">
              <Button variant="outlined" size="small" startIcon={<CopyIcon />} onClick={copyInviteCode} sx={{ borderColor: 'rgba(255,255,255,0.15)' }}>
                {group.inviteCode}
              </Button>
            </Tooltip>
            <Button variant="contained" startIcon={<AddIcon />} onClick={handleOpenExpenseDialog}>
              Add Expense
            </Button>
          </Stack>
        </Stack>
        <Typography color="text.secondary">{group.description}</Typography>
        <Chip label={group.currency} size="small" sx={{ mt: 1, backgroundColor: 'rgba(124,77,255,0.15)', color: '#B47CFF', fontWeight: 600 }} />
      </Box>

      {/* Tabs */}
      <Tabs
        value={tab}
        onChange={(_, v) => setTab(v)}
        sx={{
          mb: 3,
          '& .MuiTab-root': { textTransform: 'none', fontWeight: 600, fontSize: '0.95rem' },
          '& .MuiTabs-indicator': { backgroundColor: '#7C4DFF', height: 3, borderRadius: 2 },
        }}
      >
        <Tab icon={<ReceiptIcon />} iconPosition="start" label="Expenses" />
        <Tab icon={<SettleIcon />} iconPosition="start" label="Balances & Settle" />
        <Tab icon={<PeopleIcon />} iconPosition="start" label="Members" />
      </Tabs>

      {/* Tab 0: Expenses */}
      {tab === 0 && (
        <Box>
          {expenses.length === 0 ? (
            <Card sx={{ textAlign: 'center', py: 6 }}>
              <ReceiptIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary">No expenses yet</Typography>
              <Typography color="text.secondary" sx={{ mb: 3 }}>Add your first expense to get started.</Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={handleOpenExpenseDialog}>Add Expense</Button>
            </Card>
          ) : (
            <Stack spacing={2}>
              {expenses.map((exp: any) => (
                <Card key={exp.id} sx={{ transition: 'all 0.2s', '&:hover': { borderColor: 'rgba(124,77,255,0.2)' } }}>
                  <CardContent sx={{ display: 'flex', alignItems: 'center', gap: 2, py: 2 }}>
                    <Avatar sx={{ bgcolor: 'rgba(124,77,255,0.15)', color: '#7C4DFF' }}>
                      <ReceiptIcon />
                    </Avatar>
                    <Box sx={{ flex: 1 }}>
                      <Typography fontWeight={600}>{exp.description}</Typography>
                      <Typography variant="body2" color="text.secondary">
                        Paid by {exp.payer?.name || 'Unknown'} • {exp.splitType.toLowerCase()} split
                        {exp.splits && ` • ${exp.splits.length} people`}
                      </Typography>
                    </Box>
                    <Typography variant="h5" fontWeight={700} sx={{ color: '#10B981' }}>
                      {group.currency} {exp.amount?.toFixed(2)}
                    </Typography>
                    <IconButton size="small" onClick={() => deleteExpense({ variables: { id: exp.id } })} sx={{ color: 'text.secondary', '&:hover': { color: 'error.main' } }}>
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </CardContent>
                </Card>
              ))}
            </Stack>
          )}
        </Box>
      )}

      {/* Tab 1: Balances & Settlement */}
      {tab === 1 && (
        <Grid container spacing={3}>
          {/* Balances */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h5" fontWeight={700} sx={{ mb: 2 }}>Balances</Typography>
                {balances.length === 0 ? (
                  <Typography color="text.secondary">No balances yet. Add some expenses first.</Typography>
                ) : (
                  <Stack spacing={1.5}>
                    {balances.map((b: any) => (
                      <Box key={b.userId} sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <Avatar sx={{ width: 36, height: 36, bgcolor: b.balance >= 0 ? '#10B981' : '#EF4444', fontSize: '0.85rem' }}>
                          {b.user?.name?.charAt(0)?.toUpperCase() || 'U'}
                        </Avatar>
                        <Typography sx={{ flex: 1 }} fontWeight={500}>{b.user?.name || b.userId}</Typography>
                        <Typography fontWeight={700} sx={{ color: b.balance >= 0 ? '#10B981' : '#EF4444' }}>
                          {b.balance >= 0 ? '+' : ''}{b.balance?.toFixed(2)} {group.currency}
                        </Typography>
                      </Box>
                    ))}
                  </Stack>
                )}
              </CardContent>
            </Card>
          </Grid>

          {/* Optimal Settlement Plan */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h5" fontWeight={700} sx={{ mb: 2 }}>Optimal Settlement</Typography>
                {transactions.length === 0 ? (
                  <Typography color="text.secondary">Everyone is settled up! 🎉</Typography>
                ) : (
                  <Stack spacing={2}>
                    {transactions.map((t: any, i: number) => (
                      <Box key={i} sx={{ display: 'flex', alignItems: 'center', gap: 1.5, p: 1.5, borderRadius: 2, backgroundColor: 'rgba(255,255,255,0.03)' }}>
                        <Avatar sx={{ width: 32, height: 32, bgcolor: '#EF4444', fontSize: '0.8rem' }}>
                          {t.fromUser?.name?.charAt(0) || '?'}
                        </Avatar>
                        <Typography variant="body2" fontWeight={500}>{t.fromUser?.name || t.fromUserId}</Typography>
                        <ArrowIcon sx={{ color: '#7C4DFF', mx: 0.5 }} />
                        <Avatar sx={{ width: 32, height: 32, bgcolor: '#10B981', fontSize: '0.8rem' }}>
                          {t.toUser?.name?.charAt(0) || '?'}
                        </Avatar>
                        <Typography variant="body2" fontWeight={500}>{t.toUser?.name || t.toUserId}</Typography>
                        <Box sx={{ flex: 1 }} />
                        <Typography fontWeight={700} sx={{ color: '#F59E0B' }}>
                          {group.currency} {t.amount?.toFixed(2)}
                        </Typography>
                        <Button
                          size="small"
                          variant="outlined"
                          onClick={() => recordSettlement({ variables: { groupId, toUserId: t.toUserId, amount: t.amount } })}
                          sx={{ borderColor: 'rgba(16,185,129,0.5)', color: '#10B981', fontSize: '0.75rem', minWidth: 'auto', px: 1.5 }}
                        >
                          Settle
                        </Button>
                      </Box>
                    ))}
                  </Stack>
                )}
              </CardContent>
            </Card>
          </Grid>

          {/* Settlement History */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h5" fontWeight={700} sx={{ mb: 2 }}>Settlement History</Typography>
                {settlements.length === 0 ? (
                  <Typography color="text.secondary">No settlements recorded yet.</Typography>
                ) : (
                  <Stack spacing={1.5}>
                    {settlements.map((s: any) => (
                      <Box key={s.id} sx={{ display: 'flex', alignItems: 'center', gap: 2, p: 1.5, borderRadius: 2, backgroundColor: 'rgba(255,255,255,0.02)' }}>
                        <Typography fontWeight={500} sx={{ flex: 1 }}>
                          {s.fromUser?.name || s.fromUserId} → {s.toUser?.name || s.toUserId}
                        </Typography>
                        <Typography fontWeight={700}>{group.currency} {s.amount?.toFixed(2)}</Typography>
                        <Chip
                          label={s.status}
                          size="small"
                          color={s.status === 'COMPLETED' ? 'success' : 'warning'}
                          sx={{ fontWeight: 600 }}
                        />
                        {s.status === 'PENDING' && (
                          <IconButton size="small" onClick={() => markCompleted({ variables: { id: s.id } })} sx={{ color: '#10B981' }}>
                            <CheckIcon fontSize="small" />
                          </IconButton>
                        )}
                      </Box>
                    ))}
                  </Stack>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Tab 2: Members */}
      {tab === 2 && (
        <Card>
          <CardContent>
            <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
              <Typography variant="h5" fontWeight={700}>Members ({group.members?.length || 0})</Typography>
              {isOwner && (
                <Button size="small" variant="outlined" onClick={() => generateInvite({ variables: { groupId } })} sx={{ borderColor: 'rgba(255,255,255,0.15)' }}>
                  Regenerate Invite Code
                </Button>
              )}
            </Stack>
            <List>
              {(group.members || []).map((member: any, i: number) => (
                <React.Fragment key={member.userId}>
                  {i > 0 && <Divider sx={{ borderColor: 'rgba(255,255,255,0.04)' }} />}
                  <ListItem sx={{ px: 0 }}>
                    <ListItemAvatar>
                      <Avatar sx={{ bgcolor: '#7C4DFF' }}>
                        {member.user?.name?.charAt(0)?.toUpperCase() || 'U'}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={member.user?.name || member.userId}
                      secondary={member.user?.email}
                    />
                    <Chip
                      label={member.role}
                      size="small"
                      sx={{
                        mr: 1,
                        backgroundColor: member.role === 'owner' ? 'rgba(124,77,255,0.15)' : 'rgba(255,255,255,0.06)',
                        color: member.role === 'owner' ? '#B47CFF' : 'text.secondary',
                        fontWeight: 600,
                      }}
                    />
                    {isOwner && member.role !== 'owner' && (
                      <IconButton
                        size="small"
                        onClick={() => removeMember({ variables: { groupId, userId: member.userId } })}
                        sx={{ color: 'text.secondary', '&:hover': { color: 'error.main' } }}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    )}
                  </ListItem>
                </React.Fragment>
              ))}
            </List>
          </CardContent>
        </Card>
      )}

      {/* Create Expense Dialog */}
      <Dialog open={expenseDialogOpen} onClose={() => setExpenseDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ fontWeight: 700 }}>Add New Expense</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth label="Description" value={expDesc}
            onChange={(e) => setExpDesc(e.target.value)}
            sx={{ mt: 1, mb: 2.5 }}
            placeholder="e.g. Dinner, Taxi, Hotel"
          />
          <TextField
            fullWidth label="Amount" type="number" value={expAmount}
            onChange={(e) => setExpAmount(e.target.value)}
            sx={{ mb: 2.5 }}
            InputProps={{ startAdornment: <Typography sx={{ mr: 1, color: 'text.secondary' }}>{group.currency}</Typography> }}
          />
          <TextField
            fullWidth select label="Split Type" value={expSplitType}
            onChange={(e) => setExpSplitType(e.target.value)}
            sx={{ mb: 2.5 }}
          >
            {splitTypes.map((t) => <MenuItem key={t.value} value={t.value}>{t.label}</MenuItem>)}
          </TextField>

          {expSplitType !== 'EQUAL' && (
            <Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1.5 }}>
                {expSplitType === 'CUSTOM' ? 'Enter amount for each person:' : 'Enter percentage for each person:'}
              </Typography>
              {expSplits.map((split, i) => {
                const member = group.members?.find((m: any) => m.userId === split.userId);
                return (
                  <Box key={split.userId} sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 1.5 }}>
                    <Typography sx={{ flex: 1, fontSize: '0.9rem' }}>{member?.user?.name || split.userId}</Typography>
                    <TextField
                      size="small"
                      type="number"
                      label={expSplitType === 'CUSTOM' ? 'Amount' : '%'}
                      value={expSplitType === 'CUSTOM' ? split.amount || '' : split.percentage || ''}
                      onChange={(e) => {
                        const val = parseFloat(e.target.value) || 0;
                        const updated = [...expSplits];
                        if (expSplitType === 'CUSTOM') updated[i].amount = val;
                        else updated[i].percentage = val;
                        setExpSplits(updated);
                      }}
                      sx={{ width: 120 }}
                    />
                  </Box>
                );
              })}
            </Box>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 3, pt: 1 }}>
          <Button onClick={() => { setExpenseDialogOpen(false); resetExpenseForm(); }} sx={{ color: 'text.secondary' }}>Cancel</Button>
          <Button variant="contained" disabled={!expAmount || !expDesc || creatingExpense} onClick={handleCreateExpense}>
            {creatingExpense ? <CircularProgress size={20} /> : 'Add Expense'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default function GroupDetailPage() {
  return (
    <Providers>
      <AppLayout>
        <GroupDetailContent />
      </AppLayout>
    </Providers>
  );
}
