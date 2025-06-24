import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { DataTable } from '@/components/table/DataTable';
import { Button } from '@/components/ui/button';
import { Site, getSites, createSite, updateSite, deleteSite } from '@/api/site';
import { siteColumns } from './siteColumns';
import SiteFormDialog from './SiteFormDialog';
import { useToast } from '@/components/ui/use-toast';
import { Input } from '@/components/ui/input';

const SiteSettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedSite, setSelectedSite] = useState<Site | null>(null);
  const [search, setSearch] = useState('');

  const { data: sites, isLoading, isError } = useQuery({
    queryKey: ['sites'],
    queryFn: getSites,
    staleTime: 5 * 60 * 1000,
  });

  // 过滤站点
  const filteredSites = React.useMemo(() => {
    if (!sites) return [];
    return sites.filter(site => 
      site.name.toLowerCase().includes(search.toLowerCase()) ||
      site.domain.toLowerCase().includes(search.toLowerCase())
    );
  }, [sites, search]);

  const createMutation = useMutation({
    mutationFn: createSite,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sites'] });
      setIsDialogOpen(false);
      toast({ 
        title: t('siteCreatedSuccess'),
        description: t('siteCreatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('siteCreateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const updateMutation = useMutation({
    mutationFn: updateSite,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sites'] });
      setIsDialogOpen(false);
      setSelectedSite(null);
      toast({ 
        title: t('siteUpdatedSuccess'),
        description: t('siteUpdatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('siteUpdateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteSite,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sites'] });
      toast({ 
        title: t('siteDeletedSuccess'),
        description: t('siteDeletedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('siteDeleteFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const handleEdit = (site: Site) => {
    setSelectedSite(site);
    setIsDialogOpen(true);
  };

  const handleDelete = (id: string) => {
    if (window.confirm(t('confirmDeleteSite'))) {
      deleteMutation.mutate(id);
    }
  };

  const handleSubmit = (site: Site) => {
    if (site.id) {
      updateMutation.mutate(site);
    } else {
      createMutation.mutate(site);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{t('siteManagement')}</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          {t('addSite')}
        </Button>
      </div>

      <div className="flex items-center py-4">
        <Input
          placeholder={t('searchSitesPlaceholder')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
      </div>

      {isError ? (
        <div className="text-red-500">{t('loadSitesFailed')}</div>
      ) : (
        <DataTable
          data={filteredSites}
          columns={siteColumns(handleEdit, handleDelete, t)}
          isLoading={isLoading}
        />
      )}

      <SiteFormDialog
        open={isDialogOpen}
        onOpenChange={(open) => {
          if (!open) {
            setIsDialogOpen(false);
            setSelectedSite(null);
          }
        }}
        site={selectedSite}
        onSubmit={handleSubmit}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  );
};

export default SiteSettingsPage;
