'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import {
  Box, Typography, Button, Container, Grid, Card, CardContent,
  Stack, useTheme,
} from '@mui/material';
import {
  GroupWork as GroupIcon,
  ReceiptLong as ReceiptIcon,
  AccountBalance as SettleIcon,
  Speed as SpeedIcon,
  Security as SecurityIcon,
  TrendingUp as TrendingIcon,
} from '@mui/icons-material';
import Providers from './providers';
import { useAuthStore } from '@/stores/authStore';

const features = [
  { icon: <GroupIcon sx={{ fontSize: 40 }} />, title: 'Group Management', desc: 'Create groups, invite members with codes, and manage shared expenses effortlessly.' },
  { icon: <ReceiptIcon sx={{ fontSize: 40 }} />, title: 'Smart Splitting', desc: 'Split expenses equally, by custom amounts, or by percentage. Full flexibility.' },
  { icon: <SettleIcon sx={{ fontSize: 40 }} />, title: 'Optimal Settlement', desc: 'Our algorithm minimizes the number of transactions needed to settle all debts.' },
  { icon: <SpeedIcon sx={{ fontSize: 40 }} />, title: 'Real-time Balances', desc: 'See who owes whom instantly. Balances update as expenses are added.' },
  { icon: <SecurityIcon sx={{ fontSize: 40 }} />, title: 'Secure & Private', desc: 'JWT authentication, encrypted data, and ACID-compliant financial records.' },
  { icon: <TrendingIcon sx={{ fontSize: 40 }} />, title: 'Activity Timeline', desc: 'Full history of expenses and settlements for transparency and auditing.' },
];

function LandingContent() {
  const router = useRouter();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return (
    <Box sx={{ minHeight: '100vh', overflow: 'hidden' }}>
      {/* Hero Section */}
      <Box
        sx={{
          position: 'relative',
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: '-20%',
            right: '-10%',
            width: 600,
            height: 600,
            borderRadius: '50%',
            background: 'radial-gradient(circle, rgba(124,77,255,0.15) 0%, transparent 70%)',
            filter: 'blur(60px)',
          },
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: '-10%',
            left: '-10%',
            width: 500,
            height: 500,
            borderRadius: '50%',
            background: 'radial-gradient(circle, rgba(0,229,255,0.1) 0%, transparent 70%)',
            filter: 'blur(60px)',
          },
        }}
      >
        <Container maxWidth="lg" sx={{ position: 'relative', zIndex: 1 }}>
          {/* Navbar */}
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', py: 3 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
              <Box
                sx={{
                  width: 44, height: 44, borderRadius: 2,
                  background: 'linear-gradient(135deg, #7C4DFF 0%, #00E5FF 100%)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  fontWeight: 900, fontSize: '1.3rem', color: '#fff',
                }}
              >
                S
              </Box>
              <Typography variant="h5" fontWeight={800} sx={{ background: 'linear-gradient(135deg, #7C4DFF, #00E5FF)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                Split
              </Typography>
            </Box>
            <Stack direction="row" spacing={1.5}>
              {isAuthenticated ? (
                <Button variant="contained" onClick={() => router.push('/dashboard')}>Dashboard</Button>
              ) : (
                <>
                  <Button variant="outlined" onClick={() => router.push('/login')} sx={{ borderColor: 'rgba(255,255,255,0.2)', color: '#fff', '&:hover': { borderColor: '#7C4DFF' } }}>
                    Sign In
                  </Button>
                  <Button variant="contained" onClick={() => router.push('/register')}>
                    Get Started
                  </Button>
                </>
              )}
            </Stack>
          </Box>

          {/* Hero Content */}
          <Box sx={{ textAlign: 'center', mt: { xs: 8, md: 12 }, mb: 8 }}>
            <Typography
              variant="h1"
              sx={{
                fontSize: { xs: '2.5rem', md: '4rem', lg: '4.5rem' },
                fontWeight: 900,
                lineHeight: 1.1,
                mb: 3,
                background: 'linear-gradient(135deg, #F9FAFB 0%, #9CA3AF 100%)',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
              }}
            >
              Split expenses.
              <br />
              <Box component="span" sx={{ background: 'linear-gradient(135deg, #7C4DFF 0%, #00E5FF 100%)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                Settle smart.
              </Box>
            </Typography>
            <Typography
              variant="h5"
              sx={{ color: 'text.secondary', maxWidth: 640, mx: 'auto', mb: 5, fontWeight: 400, lineHeight: 1.6-0, fontSize: { xs: '1.1rem', md: '1.3rem' } }}
            >
              Track group expenses, split bills any way you want, and let our algorithm find the fewest payments to settle up.
            </Typography>
            <Stack direction="row" spacing={2} justifyContent="center">
              <Button variant="contained" size="large" onClick={() => router.push('/register')} sx={{ px: 5, py: 1.5, fontSize: '1.1rem' }}>
                Start for Free
              </Button>
              <Button variant="outlined" size="large" sx={{ px: 5, py: 1.5, fontSize: '1.1rem', borderColor: 'rgba(255,255,255,0.15)', color: '#fff', '&:hover': { borderColor: '#7C4DFF', backgroundColor: 'rgba(124,77,255,0.08)' } }}>
                Learn More
              </Button>
            </Stack>
          </Box>
        </Container>
      </Box>

      {/* Features Section */}
      <Box sx={{ py: { xs: 8, md: 12 }, position: 'relative' }}>
        <Container maxWidth="lg">
          <Typography variant="h2" textAlign="center" sx={{ mb: 2 }}>
            Everything you need
          </Typography>
          <Typography variant="h6" textAlign="center" color="text.secondary" sx={{ mb: 8, maxWidth: 500, mx: 'auto', fontWeight: 400 }}>
            Powerful features designed to make group finances simple and stress-free.
          </Typography>
          <Grid container spacing={3}>
            {features.map((f, i) => (
              <Grid item xs={12} sm={6} md={4} key={i}>
                <Card
                  sx={{
                    height: '100%',
                    transition: 'all 0.3s ease',
                    '&:hover': {
                      transform: 'translateY(-4px)',
                      boxShadow: '0 12px 40px rgba(124,77,255,0.15)',
                      borderColor: 'rgba(124,77,255,0.3)',
                    },
                  }}
                >
                  <CardContent sx={{ p: 3.5 }}>
                    <Box sx={{ color: '#7C4DFF', mb: 2 }}>{f.icon}</Box>
                    <Typography variant="h5" sx={{ mb: 1.5 }}>{f.title}</Typography>
                    <Typography variant="body2" color="text.secondary" lineHeight={1.7}>
                      {f.desc}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* CTA Section */}
      <Box sx={{ py: { xs: 8, md: 12 } }}>
        <Container maxWidth="sm">
          <Card sx={{ textAlign: 'center', p: { xs: 4, md: 6 }, background: 'linear-gradient(135deg, rgba(124,77,255,0.1) 0%, rgba(0,229,255,0.05) 100%)', border: '1px solid rgba(124,77,255,0.2)' }}>
            <Typography variant="h3" sx={{ mb: 2 }}>Ready to simplify?</Typography>
            <Typography color="text.secondary" sx={{ mb: 4, fontSize: '1.1rem' }}>
              Join Split and never argue about money again.
            </Typography>
            <Button variant="contained" size="large" onClick={() => router.push('/register')} sx={{ px: 6 }}>
              Create Free Account
            </Button>
          </Card>
        </Container>
      </Box>

      {/* Footer */}
      <Box sx={{ py: 4, borderTop: '1px solid rgba(255,255,255,0.06)' }}>
        <Container maxWidth="lg">
          <Typography variant="body2" color="text.secondary" textAlign="center">
            © 2026 Split. Built with ❤️ for hassle-free group expenses.
          </Typography>
        </Container>
      </Box>
    </Box>
  );
}

export default function Home() {
  return (
    <Providers>
      <LandingContent />
    </Providers>
  );
}
