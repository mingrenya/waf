import { ColumnDef } from '@tanstack/react-table';
import { AttackLog } from '@/types/logs';
import { format } from 'date-fns';
import { Badge } from '@/components/ui/badge';

export const attackLogColumns = (
  t: (key: string) => string
): ColumnDef<AttackLog>[] => [
  {
    accessorKey: 'timestamp',
    header: t('timestamp'),
    cell: ({ row }) => format(new Date(row.original.timestamp), 'yyyy-MM-dd HH:mm:ss'),
  },
  {
    accessorKey: 'clientIP',
    header: t('clientIP'),
  },
  {
    accessorKey: 'method',
    header: t('method'),
  },
  {
    accessorKey: 'path',
    header: t('path'),
  },
  {
    accessorKey: 'ruleId',
    header: t('ruleId'),
  },
  {
    accessorKey: 'severity',
    header: t('severity'),
    cell: ({ row }) => {
      const severity = row.original.severity;
      let variant: 'destructive' | 'warning' | 'default' = 'default';
      if (severity >= 8) variant = 'destructive';
      else if (severity >= 5) variant = 'warning';
      
      return <Badge variant={variant}>{severity}</Badge>;
    },
  },
  {
    accessorKey: 'message',
    header: t('message'),
    cell: ({ row }) => (
      <div className="max-w-xs truncate" title={row.original.message}>
        {row.original.message}
      </div>
    ),
  },
];
