import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { DataTable } from '@/components/table/DataTable';
import { Button } from '@/components/ui/button';
import { Certificate, getCertificates, createCertificate, updateCertificate, deleteCertificate } from '@/api/certificate';
import { certificateColumns } from './certificateColumns';
import CertificateFormDialog from './CertificateFormDialog';
import { useToast } from '@/components/ui/use-toast';
import { Input } from '@/components/ui/input';

const CertificateSettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedCertificate, setSelectedCertificate] = useState<Certificate | null>(null);
  const [search, setSearch] = useState('');

  const { data: certificates, isLoading, isError } = useQuery({
    queryKey: ['certificates'],
    queryFn: getCertificates,
    staleTime: 5 * 60 * 1000,
  });

  // 过滤证书
  const filteredCertificates = React.useMemo(() => {
    if (!certificates) return [];
    return certificates.filter(cert => 
      cert.domain.toLowerCase().includes(search.toLowerCase()) ||
      cert.issuer.toLowerCase().includes(search.toLowerCase())
    );
  }, [certificates, search]);

  const createMutation = useMutation({
    mutationFn: createCertificate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['certificates'] });
      setIsDialogOpen(false);
      toast({ 
        title: t('certificateCreatedSuccess'),
        description: t('certificateCreatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('certificateCreateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const updateMutation = useMutation({
    mutationFn: updateCertificate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['certificates'] });
      setIsDialogOpen(false);
      setSelectedCertificate(null);
      toast({ 
        title: t('certificateUpdatedSuccess'),
        description: t('certificateUpdatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('certificateUpdateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const deleteMutation = useMutation({
    mutationFn: deleteCertificate,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['certificates'] });
      toast({ 
        title: t('certificateDeletedSuccess'),
        description: t('certificateDeletedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({ 
        title: t('certificateDeleteFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const handleEdit = (certificate: Certificate) => {
    setSelectedCertificate(certificate);
    setIsDialogOpen(true);
  };

  const handleDelete = (id: string) => {
    if (window.confirm(t('confirmDeleteCertificate'))) {
      deleteMutation.mutate(id);
    }
  };

  const handleSubmit = (certificate: Certificate) => {
    if (certificate.id) {
      updateMutation.mutate(certificate);
    } else {
      createMutation.mutate(certificate);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{t('certificateManagement')}</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          {t('addCertificate')}
        </Button>
      </div>

      <div className="flex items-center py-4">
        <Input
          placeholder={t('searchCertificatesPlaceholder')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
      </div>

      {isError ? (
        <div className="text-red-500">{t('loadCertificatesFailed')}</div>
      ) : (
        <DataTable
          data={filteredCertificates}
          columns={certificateColumns(handleEdit, handleDelete, t)}
          isLoading={isLoading}
        />
      )}

      <CertificateFormDialog
        open={isDialogOpen}
        onOpenChange={(open) => {
          if (!open) {
            setIsDialogOpen(false);
            setSelectedCertificate(null);
          }
        }}
        certificate={selectedCertificate}
        onSubmit={handleSubmit}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  );
};

export default CertificateSettingsPage;
