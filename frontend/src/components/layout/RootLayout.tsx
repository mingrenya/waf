import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Layout } from '@/components/ui/layout';
import Sidebar from './Sidebar';
import Breadcrumb from './Breadcrumb';
import { Button } from '@/components/ui/button';
import { Languages } from 'lucide-react';

const RootLayout: React.FC = () => {
  const { t, i18n } = useTranslation();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const toggleLanguage = () => {
    const newLang = i18n.language === 'zh' ? 'en' : 'zh';
    i18n.changeLanguage(newLang);
  };

  return (
    <Layout className="flex h-screen">
      <Sidebar 
        isOpen={sidebarOpen} 
        toggleOpen={() => setSidebarOpen(!sidebarOpen)} 
      />
      
      <div className="flex flex-col flex-1 overflow-hidden">
        <header className="flex items-center justify-between p-4 border-b">
          <Breadcrumb />
          <div className="flex items-center gap-4">
            <Button 
              variant="outline" 
              size="icon"
              onClick={toggleLanguage}
              title={t('toggleLanguage')}
            >
              <Languages className="h-4 w-4" />
            </Button>
            {/* 用户头像/信息 */}
          </div>
        </header>
        
        <main className="flex-1 overflow-auto p-4">
          <Outlet />
        </main>
      </div>
    </Layout>
  );
};

export default RootLayout;
