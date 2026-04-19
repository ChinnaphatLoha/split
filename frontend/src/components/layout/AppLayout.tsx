'use client';

import React, { useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import {
  AppBar, Toolbar, Typography, Box, IconButton, Drawer, List, ListItem,
  ListItemButton, ListItemIcon, ListItemText, Avatar, Divider, useMediaQuery,
  useTheme, Tooltip,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  Group as GroupIcon,
  Person as PersonIcon,
  Logout as LogoutIcon,
  Receipt as ReceiptIcon,
} from '@mui/icons-material';
import { useAuthStore } from '@/stores/authStore';

const DRAWER_WIDTH = 260;

const navItems = [
  { label: 'Dashboard', path: '/dashboard', icon: <DashboardIcon /> },
  { label: 'Groups', path: '/groups', icon: <GroupIcon /> },
  { label: 'Profile', path: '/profile', icon: <PersonIcon /> },
];

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const router = useRouter();
  const pathname = usePathname();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  const drawerContent = (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      {/* Logo */}
      <Box sx={{ p: 3, display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box
          sx={{
            width: 40, height: 40, borderRadius: 2,
            background: 'linear-gradient(135deg, #7C4DFF 0%, #00E5FF 100%)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontWeight: 900, fontSize: '1.2rem', color: '#fff',
          }}
        >
          S
        </Box>
        <Typography variant="h5" sx={{ fontWeight: 800, background: 'linear-gradient(135deg, #7C4DFF, #00E5FF)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
          Split
        </Typography>
      </Box>

      <Divider sx={{ borderColor: 'rgba(255,255,255,0.06)' }} />

      {/* Nav Items */}
      <List sx={{ px: 1.5, py: 2, flex: 1 }}>
        {navItems.map((item) => {
          const isActive = pathname === item.path || pathname?.startsWith(item.path + '/');
          return (
            <ListItem key={item.path} disablePadding sx={{ mb: 0.5 }}>
              <ListItemButton
                onClick={() => { router.push(item.path); setMobileOpen(false); }}
                sx={{
                  borderRadius: 2,
                  py: 1.2,
                  backgroundColor: isActive ? 'rgba(124, 77, 255, 0.15)' : 'transparent',
                  '&:hover': { backgroundColor: isActive ? 'rgba(124, 77, 255, 0.2)' : 'rgba(255,255,255,0.04)' },
                }}
              >
                <ListItemIcon sx={{ color: isActive ? '#7C4DFF' : 'text.secondary', minWidth: 40 }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.label}
                  primaryTypographyProps={{
                    fontWeight: isActive ? 600 : 400,
                    color: isActive ? '#F9FAFB' : 'text.secondary',
                    fontSize: '0.95rem',
                  }}
                />
              </ListItemButton>
            </ListItem>
          );
        })}
      </List>

      <Divider sx={{ borderColor: 'rgba(255,255,255,0.06)' }} />

      {/* User Profile & Logout */}
      <Box sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1.5, px: 1 }}>
          <Avatar sx={{ width: 36, height: 36, bgcolor: '#7C4DFF', fontSize: '0.9rem', fontWeight: 700 }}>
            {user?.name?.charAt(0)?.toUpperCase() || 'U'}
          </Avatar>
          <Box sx={{ flex: 1, minWidth: 0 }}>
            <Typography variant="body2" fontWeight={600} noWrap>{user?.name || 'User'}</Typography>
            <Typography variant="caption" color="text.secondary" noWrap>{user?.email}</Typography>
          </Box>
        </Box>
        <ListItemButton onClick={handleLogout} sx={{ borderRadius: 2, py: 1 }}>
          <ListItemIcon sx={{ minWidth: 36, color: 'error.main' }}>
            <LogoutIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText primary="Logout" primaryTypographyProps={{ fontSize: '0.9rem', color: 'error.main' }} />
        </ListItemButton>
      </Box>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      {/* Desktop Drawer */}
      {!isMobile && (
        <Drawer
          variant="permanent"
          sx={{
            width: DRAWER_WIDTH,
            flexShrink: 0,
            '& .MuiDrawer-paper': {
              width: DRAWER_WIDTH,
              boxSizing: 'border-box',
              backgroundColor: '#0F1320',
              borderRight: '1px solid rgba(255,255,255,0.06)',
            },
          }}
        >
          {drawerContent}
        </Drawer>
      )}

      {/* Mobile Drawer */}
      {isMobile && (
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={() => setMobileOpen(false)}
          sx={{
            '& .MuiDrawer-paper': {
              width: DRAWER_WIDTH,
              backgroundColor: '#0F1320',
            },
          }}
        >
          {drawerContent}
        </Drawer>
      )}

      {/* Main Content */}
      <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        {isMobile && (
          <AppBar position="sticky" elevation={0} sx={{ backgroundColor: '#0F1320', borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
            <Toolbar>
              <IconButton edge="start" color="inherit" onClick={() => setMobileOpen(true)}>
                <MenuIcon />
              </IconButton>
              <Typography variant="h6" sx={{ fontWeight: 700, ml: 1, background: 'linear-gradient(135deg, #7C4DFF, #00E5FF)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                Split
              </Typography>
            </Toolbar>
          </AppBar>
        )}

        <Box
          component="main"
          sx={{
            flex: 1,
            p: { xs: 2, sm: 3, md: 4 },
            maxWidth: 1200,
            width: '100%',
            mx: 'auto',
          }}
        >
          {children}
        </Box>
      </Box>
    </Box>
  );
}
