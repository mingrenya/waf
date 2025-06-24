import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { DataTable } from '@/components/table/DataTable';
import { attackLogColumns } from './attackLogColumns';
import { getAttackLogs } from '@/api/logs';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { RefreshCw } from 'lucide-react';
import { DateRangePicker } from '@/components/common/DateRangePicker';
import { useDateFilter } from '@/hooks/useDateFilter';

const AttackLogs: React.FC = () => {
  const { t } = useTranslation();
  const [search, setSearch] = useState('');
  const { dateRange, setDateRange } = useDateFilter();
  
  // 获取攻击日志
  const { data, isLoading, isError, refetch, isRefetching } = useQuery({
    queryKey: ['attackLogs', dateRange, search],
    queryFn: () => getAttackLogs({
      startDate: dateRange?.from?.toISOString(),
      endDate: dateRange?.to?.toISOString(),
      search,
    }),
    keepPreviousData: true,
  });

  const handleRefresh = () => {
    refetch();
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">{t('attackLogs')}</h2>
        <Button 
          variant="outline" 
          onClick={handleRefresh}
          disabled={isRefetching}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${isRefetching ? 'animate-spin' : ''}`} />
          {t('refresh')}
        </Button>
      </div>

      <div className="flex flex-col md:flex-row gap-4">
        <Input
          placeholder={t('searchLogsPlaceholder')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="md:max-w-sm"
        />
        <DateRangePicker 
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
        />
      </div>

      {isError ? (
        <div className="text-red-500">{t('loadLogsFailed')}</div>
      ) : (
        <DataTable
          data={data || []}
          columns={attackLogColumns(t)}
          isLoading={isLoading || isRefetching}
        />
      )}
    </div>
  );
};

export default AttackLogs;
