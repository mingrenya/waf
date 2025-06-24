import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { DataTable } from '@/components/table/DataTable';
import { Button } from '@/components/ui/button';
import { Rule, getRules, createRule, updateRule, deleteRule } from '@/api/rule';
import { ruleColumns } from './ruleColumns';
import RuleFormDialog from './RuleFormDialog';
import { useToast } from '@/components/ui/use-toast';
import { Input } from '@/components/ui/input';

const RuleManagementPage: React.FC = () => {
  const { t } = useTranslation();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
  const [search, setSearch] = useState('');

  // 获取规则列表
  const { data: rules, isLoading, isError } = useQuery({
    queryKey: ['rules'],
    queryFn: getRules,
    staleTime: 5 * 60 * 1000,
  });

  // 过滤规则
  const filteredRules = React.useMemo(() => {
    if (!rules) return [];
    return rules.filter(rule => 
      rule.name.toLowerCase().includes(search.toLowerCase()) ||
      rule.description?.toLowerCase().includes(search.toLowerCase()) ||
      rule.id.toLowerCase().includes(search.toLowerCase())
    );
  }, [rules, search]);

  // 创建规则
  const createMutation = useMutation({
    mutationFn: createRule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rules'] });
      setIsDialogOpen(false);
      toast({ 
        title: t('ruleCreatedSuccess'),
        description: t('ruleCreatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({
        title: t('ruleCreateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  // 更新规则
  const updateMutation = useMutation({
    mutationFn: updateRule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rules'] });
      setIsDialogOpen(false);
      setSelectedRule(null);
      toast({ 
        title: t('ruleUpdatedSuccess'),
        description: t('ruleUpdatedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({
        title: t('ruleUpdateFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  // 删除规则
  const deleteMutation = useMutation({
    mutationFn: deleteRule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rules'] });
      toast({ 
        title: t('ruleDeletedSuccess'),
        description: t('ruleDeletedSuccessDesc')
      });
    },
    onError: (error) => {
      toast({
        title: t('ruleDeleteFailed'),
        description: error.message,
        variant: 'destructive'
      });
    }
  });

  const handleEdit = (rule: Rule) => {
    setSelectedRule(rule);
    setIsDialogOpen(true);
  };

  const handleDelete = (id: string) => {
    if (window.confirm(t('confirmDeleteRule'))) {
      deleteMutation.mutate(id);
    }
  };

  const handleSubmit = (rule: Rule) => {
    if (rule.id) {
      updateMutation.mutate(rule);
    } else {
      createMutation.mutate(rule);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">{t('ruleManagement')}</h1>
        <Button onClick={() => setIsDialogOpen(true)}>
          {t('addRule')}
        </Button>
      </div>

      <div className="flex items-center py-4">
        <Input
          placeholder={t('searchRulesPlaceholder')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-sm"
        />
      </div>

      {isError ? (
        <div className="text-red-500">{t('loadRulesFailed')}</div>
      ) : (
        <DataTable
          data={filteredRules}
          columns={ruleColumns(handleEdit, handleDelete, t)}
          isLoading={isLoading}
        />
      )}

      <RuleFormDialog
        open={isDialogOpen}
        onOpenChange={(open) => {
          if (!open) {
            setIsDialogOpen(false);
            setSelectedRule(null);
          }
        }}
        rule={selectedRule}
        onSubmit={handleSubmit}
        isSubmitting={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  );
};

export default RuleManagementPage;
