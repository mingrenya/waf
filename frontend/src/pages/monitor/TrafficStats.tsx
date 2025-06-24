import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { useTranslation } from 'react-i18next';
import { TrafficStat } from '@/types/monitor';

interface TrafficStatsProps {
  data: TrafficStat[];
}

const TrafficStats: React.FC<TrafficStatsProps> = ({ data }) => {
  const { t } = useTranslation();

  return (
    <div className="h-80">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart
          data={data}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line 
            type="monotone" 
            dataKey="total" 
            stroke="#8884d8" 
            name={t('totalRequests')} 
          />
          <Line 
            type="monotone" 
            dataKey="blocked" 
            stroke="#ff7300" 
            name={t('blockedRequests')} 
          />
          <Line 
            type="monotone" 
            dataKey="attacks" 
            stroke="#82ca9d" 
            name={t('attackRequests')} 
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default TrafficStats;
