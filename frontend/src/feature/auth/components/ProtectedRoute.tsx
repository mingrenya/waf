import React, { useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import { ROUTES } from '@/routes/constants';
import { useQuery } from '@tanstack/react-query';
import { getCurrentUser } from '@/api/auth';
import { Loader2 } from 'lucide-react';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireAuth?: boolean;
  roles?: string[];
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ 
  children, 
  requireAuth = true,
  roles = []
}) => {
  const location = useLocation();
  const { token, user, isAuthenticated, setUser } = useAuthStore();
  
  // 获取当前用户信息
  const { isLoading, isError } = useQuery({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
    enabled: !!token && !user,
    onSuccess: (data) => {
      setUser(data);
    }
  });

  // 检查角色权限
  const hasRequiredRole = () => {
    if (roles.length === 0) return true;
    return roles.some(role => user?.role === role);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Loader2 className="h-12 w-12 animate-spin" />
      </div>
    );
  }

  if (requireAuth && !isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />;
  }

  if (requireAuth && isAuthenticated && !hasRequiredRole()) {
    return <Navigate to={ROUTES.FORBIDDEN} replace />;
  }

  if (!requireAuth && isAuthenticated) {
    return <Navigate to={ROUTES.DASHBOARD} replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
