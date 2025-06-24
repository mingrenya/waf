import { ColumnDef } from '@tanstack/react-table';
import { Rule } from '@/api/rule';
import { Button } from '@/components/ui/button';
import { Pen, Trash2 } from 'lucide-react';

export const ruleColumns = (
  onEdit: (rule: Rule) => void,
  onDelete: (id: string) => void,
  t: (key: string) => string
): ColumnDef<Rule>[] => [
  {
    accessorKey: 'id',
    header: t('id'),
  },
  {
    accessorKey: 'name',
    header: t('name'),
  },
  {
    accessorKey: 'description',
    header: t('description'),
  },
  {
    accessorKey: 'action',
    header: t('action'),
  },
  {
    accessorKey: 'status',
    header: t('status'),
    cell: ({ row }) => (
      <span className={row.original.status === 'enabled' ? 'text-green-600' : 'text-gray-500'}>
        {row.original.status === 'enabled' ? t('enabled') : t('disabled')}
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
