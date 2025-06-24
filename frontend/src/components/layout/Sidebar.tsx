import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
  LayoutDashboard, 
  Settings, 
  ShieldAlert, 
  FileText, 
  Gauge, 
  List, 
  XCircle,
  ChevronLeft,
  ChevronRight
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { ROUTES } from '@/routes/constants';
import { useAuthStore } from '@/store/auth';

interface SidebarProps {
  isOpen: boolean;
  toggleOpen: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen, toggleOpen }) => {
  const { t } = useTranslation();
  const location = useLocation();
  const user = useAuthStore(state => state.user);
  
  const navigationItems = [
    {
      name: t('dashboard'),
      path: ROUTES.DASHBOARD.path,
      icon: <LayoutDashboard className="h-5 w-5" />,
      roles: ['admin', 'user']
    },
    {
      name: t('monitor'),
      path: ROUTES.MONITOR.path,
      icon: <Gauge className="h-5 w-5" />,
      roles: ['admin', 'user']
    },
    {
      name: t('ruleManagement'),
      path: ROUTES.RULE_MANAGEMENT.path,
      icon: <ShieldAlert className="h-5 w-5" />,
      roles: ['admin']
    },
    {
      name: t('siteManagement'),
      path: ROUTES.SETTINGS_SITE.path,
      icon: <List className="h-5 w-5" />,
      roles: ['admin']
    },
    {
      name: t('certificateManagement'),
      path: ROUTES.SETTINGS_CERTIFICATE.path,
      icon: <FileText className="h-5 w-5" />,
      roles: ['admin']
    },
    {
      name: t('logs'),
      path: ROUTES.LOGS.path,
      icon: <XCircle className="h-5 w-5" />,
      roles: ['admin', 'user']
    },
    {
      name: t('settings'),
      path: '/settings',
      icon: <Settings className="h-5 w-5" />,
      roles: ['admin']
    },
  ];

  // 过滤用户有权限访问的菜单
  const filteredItems = navigationItems.filter(item => 
    item.roles.includes(user?.role || 'user')
  );

  return (
    <aside 
      className={cn(
        "bg-background border-r h-full flex flex-col transition-all duration-300 ease-in-out",
        isOpen ? "w-64" : "w-20"
      )}
    >
      <div className="flex items-center justify-between p-4 border-b">
        <div className={cn(
          "font-bold text-xl transition-opacity",
          isOpen ? "opacity-100" : "opacity-0"
        )}>
          AI WAF
        </div>
        <button 
          onClick={toggleOpen}
          className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
        >
          {isOpen ? (
            <ChevronLeft className="h-5 w-5" />
          ) : (
            <ChevronRight className="h-5 w-5" />
          )}
        </button>
      </div>
      
      <nav className="flex-1 overflow-y-auto py-4">
        <ul className="space-y-1 px-2">
          {filteredItems.map((item) => (
            <li key={item.path}>
              <NavLink
                to={item.path}
                className={({ isActive }) => cn(
                  "flex items-center p-3 rounded-lg transition-colors",
                  isActive 
                    ? "bg-primary text-primary-foreground" 
                    : "text-foreground hover:bg-gray-100 dark:hover:bg-gray-700"
                )}
              >
                <span className="flex-shrink-0">{item.icon}</span>
                <span className={cn(
                  "ml-3 transition-opacity",
                  isOpen ? "opacity-100" : "opacity-0"
                )}>
                  {item.name}
                </span>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
      
      <div className={cn(
        "p-4 border-t flex items-center transition-opacity",
        isOpen ? "opacity-100" : "opacity-0"
      )}>
        <div className="bg-gray-200 border-2 border-dashed rounded-xl w-10 h-10" />
        <div className="ml-3">
          <p className="text-sm font-medium">{user?.name || '用户'}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {user?.role === 'admin' ? t('administrator') : t('user')}
          </p>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;
