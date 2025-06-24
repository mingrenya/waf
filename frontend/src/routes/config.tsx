import React from 'react';
import { createBrowserRouter } from 'react-router-dom';
import RootLayout from '@/components/layout/RootLayout';
import ProtectedRoute from '@/feature/auth/components/ProtectedRoute';
import LoginPage from '@/pages/auth/LoginPage';
import DashboardPage from '@/pages/dashboard/DashboardPage';
import SiteSettingsPage from '@/pages/setting/site/SiteSettingsPage';
import CertificateSettingsPage from '@/pages/setting/certificate/CertificateSettingsPage';
import RuleManagementPage from '@/pages/rule/index';
import LogsPage from '@/pages/logs/index';
import MonitorPage from '@/pages/monitor/index';
import { ROUTES } from './constants';

const router = createBrowserRouter([
  {
    path: ROUTES.ROOT.path,
    element: <RootLayout />,
    children: [
      {
        path: ROUTES.LOGIN.path,
        element: (
          <ProtectedRoute requireAuth={false}>
            <LoginPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.DASHBOARD.path,
        element: (
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.SETTINGS_SITE.path,
        element: (
          <ProtectedRoute>
            <SiteSettingsPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.SETTINGS_CERTIFICATE.path,
        element: (
          <ProtectedRoute>
            <CertificateSettingsPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.RULE_MANAGEMENT.path,
        element: (
          <ProtectedRoute>
            <RuleManagementPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.LOGS.path,
        element: (
          <ProtectedRoute>
            <LogsPage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.MONITOR.path,
        element: (
          <ProtectedRoute>
            <MonitorPage />
          </ProtectedRoute>
        ),
      },
    ],
  },
]);

export default router;
