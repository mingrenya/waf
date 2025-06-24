import React, { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import AttackLogs from './AttackLogs';
import ProtectLogs from './ProtectLogs';
import { useTranslation } from 'react-i18next';

const LogsPage: React.FC = () => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('attack');

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{t('logs')}</h1>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="attack">{t('attackLogs')}</TabsTrigger>
          <TabsTrigger value="protect">{t('protectLogs')}</TabsTrigger>
        </TabsList>
        
        <TabsContent value="attack" className="mt-4">
          <AttackLogs />
        </TabsContent>
        
        <TabsContent value="protect" className="mt-4">
          <ProtectLogs />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default LogsPage;
