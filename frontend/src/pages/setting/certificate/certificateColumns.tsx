import { ColumnDef } from '@tanstack/react-table';
import { Certificate } from '@/api/certificate';
import { Button } from '@/components/ui/button';
import { Pen, Trash2 } from 'lucide-react';
import { format } from 'date-fns';

export const certificateColumns = (
  onEdit: (certificate: Certificate) => void,
  onDelete: (id: string) => void,
  t: (key: string) => string
): ColumnDef<Certificate>[] => [
  {
    accessorKey: 'domain',
    header: t('domain'),
  },
  {
    accessorKey: 'issuer',
    header: t('issuer'),
  },
  {
    accessorKey: 'validFrom',
    header: t('validFrom'),
    cell: ({ row }) => format(new Date(row.original.validFrom), 'yyyy-MM-dd'),
  },
  {
    accessorKey: 'validTo',
    header: t('validTo'),
    cell: ({ row }) => format(new Date(row.original.validTo), 'yyyy-MM-dd'),
  },
  {
    accessorKey: 'status',
    header: t('status'),
    cell: ({ row }) => (
      <span className={row.original.status === 'valid' ? 'text-green-600' : 'text-red-600'}>
        {row.original.status === 'valid' ? t('valid') : t('expired')}
      </span>
    ),
  },
  {
    id: 'actions',
    header: t('actions'),
    cell: ({ row }) => (
      <div className="flex space-x-2">
        <Button
          variant="outline"
          size="icon"
          onClick={() => onEdit(row.original)}
          aria-label={t('edit')}
        >
          <Pen className="h-4 w-4" />
        </Button>
        <Button
          variant="destructive"
          size="icon"
          onClick={() => onDelete(row.original.id)}
          aria-label={t('delete')}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    ),
  },
];
