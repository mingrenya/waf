interface RouteConfig {
  path: string;
  breadcrumb: string;
  requireAuth?: boolean;
  roles?: string[];
}

export const ROUTES = {
  ROOT: {
    path: '/',
    breadcrumb: 'home',
  },
  LOGIN: {
    path: '/login',
    breadcrumb: 'login',
    requireAuth: false,
  },
  DASHBOARD: {
    path: '/dashboard',
    breadcrumb: 'dashboard',
    requireAuth: true,
  },
  SETTINGS: {
    path: '/settings',
    breadcrumb: 'settings',
    requireAuth: true,
    roles: ['admin'],
  },
  SETTINGS_SITE: {
    path: '/settings/site',
    breadcrumb: 'siteManagement',
    requireAuth: true,
    roles: ['admin'],
  },
  SETTINGS_CERTIFICATE: {
    path: '/settings/certificate',
    breadcrumb: 'certificateManagement',
    requireAuth: true,
    roles: ['admin'],
  },
  LOGS: {
    path: '/logs',
    breadcrumb: 'logs',
    requireAuth: true,
  },
  FORBIDDEN: {
    path: '/forbidden',
    breadcrumb: 'forbidden',
    requireAuth: true,
  },
};

export type RouteKey = keyof typeof ROUTES;
