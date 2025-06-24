import React from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { getMonitorData } from '@/api/monitor';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import TrafficStats from './TrafficStats';

const MonitorPage: React.FC = () => {
  const { t } = useTranslation();
  const { data, isLoading, isRefetching } = useQuery({
    queryKey: ['monitorData'],
    queryFn: getMonitorData,
    refetchInterval: 10000, // 10秒刷新一次
  });

  const severityData = [
    { name: t('low'), value: data?.severityCount.low || 0 },
    { name: t('medium'), value: data?.severityCount.medium || 0 },
    { name: t('high'), value: data?.severityCount.high || 0 },
    { name: t('critical'), value: data?.severityCount.critical || 0 },
  ];

  const renderMetricCard = (title: string, value: number | string) => (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading || isRefetching ? (
          <Skeleton className="h-8 w-full" />
        ) : (
          <div className="text-3xl font-bold">{value}</div>
        )}
      </CardContent>
    </Card>
  );

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{t('monitorDashboard')}</h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {renderMetricCard(t('totalRequests'), data?.totalRequests || 0)}
        {renderMetricCard(t('blockedRequests'), data?.blockedRequests || 0)}
        {renderMetricCard(t('attackRequests'), data?.attackRequests || 0)}
        {renderMetricCard(
          t('averageResponseTime'), 
          data?.avgResponseTime ? `${data.avgResponseTime}ms` : '0ms'
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>{t('severityDistribution')}</CardTitle>
          </CardHeader>
          <CardContent className="h-80">
            {isLoading || isRefetching ? (
              <div className="flex items-center justify-center h-full">
                <Skeleton className="h-64 w-full" />
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={severityData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="value" fill="#8884d8" />
                </BarChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{t('trafficStats')}</CardTitle>
          </CardHeader>
          <CardContent className="h-80">
            {isLoading || isRefetching ? (
              <div className="flex items-center justify-center h-full">
                <Skeleton className="h-64 w-full" />
              </div>
            ) : (
              <TrafficStats data={data?.trafficStats || []} />
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default MonitorPage;
