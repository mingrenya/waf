import { ColumnDef } from '@tanstack/react-table';
import { Site } from '@/api/site';
import { Button } from '@/components/ui/button';
import { Pen, Trash2 } from 'lucide-react';

export const siteColumns = (
  onEdit: (site: Site) => void,
  onDelete: (id: string) => void,
  t: (key: string) => string
): ColumnDef<Site>[] => [
  {
    accessorKey: 'name',
    header: t('name'),
  },
  {
    accessorKey: 'domain',
    header: t('domain'),
  },
  {
    accessorKey: 'description',
    header: t('description'),
  },
  {
    accessorKey: 'createdAt',
    header: t('createdAt'),
    cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
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
