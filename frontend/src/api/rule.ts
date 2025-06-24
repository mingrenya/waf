import apiClient from './index';
import { Rule } from '@/types/rule';

export const getRules = async (): Promise<Rule[]> => {
  const response = await apiClient.get('/rules');
  return response.data;
};

export const createRule = async (rule: Omit<Rule, 'id'>): Promise<Rule> => {
  const response = await apiClient.post('/rules', rule);
  return response.data;
};

export const updateRule = async (rule: Rule): Promise<Rule> => {
  const response = await apiClient.put(`/rules/${rule.id}`, rule);
  return response.data;
};

export const deleteRule = async (id: string): Promise<void> => {
  await apiClient.delete(`/rules/${id}`);
};
