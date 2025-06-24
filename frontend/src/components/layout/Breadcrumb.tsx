import React from 'react';
import { useLocation, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { ChevronRight } from 'lucide-react';
import { ROUTES } from '@/routes/constants';

const Breadcrumb: React.FC = () => {
  const { t } = useTranslation();
  const location = useLocation();
  
  // 获取当前路由路径
  const pathnames = location.pathname.split('/').filter((x) => x);
  
  // 排除登录页面
  if (pathnames[0] === 'login') {
    return null;
  }
  
  // 获取面包屑名称
  const getBreadcrumbName = (path: string) => {
    const route = Object.values(ROUTES).find(r => r.path === `/${path}`);
    return route ? t(route.breadcrumb) : path;
  };

  return (
    <nav className="flex" aria-label="Breadcrumb">
      <ol className="flex items-center space-x-2">
        <li>
          <Link to={ROUTES.DASHBOARD.path} className="text-sm font-medium hover:text-primary">
            {t('dashboard')}
          </Link>
        </li>
        {pathnames.map((value, index) => {
          const to = `/${pathnames.slice(0, index + 1).join('/')}`;
          const isLast = index === pathnames.length - 1;
          
          return (
            <li key={to} className="flex items-center">
              <ChevronRight className="h-4 w-4 text-muted-foreground" />
              {isLast ? (
                <span className="ml-2 text-sm font-medium">
                  {getBreadcrumbName(value)}
                </span>
              ) : (
                <Link to={to} className="ml-2 text-sm font-medium hover:text-primary">
                  {getBreadcrumbName(value)}
                </Link>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
};

export default Breadcrumb;
