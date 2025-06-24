import apiClient from '.';
import { Site } from '@/types/site';

export const getSites = async (): Promise<Site[]> => {
  const response = await apiClient.get('/sites');
  return response.data;
};

export const getSiteById = async (id: string): Promise<Site> => {
  const response = await apiClient.get(`/sites/${id}`);
  return response.data;
};

export const createSite = async (site: Omit<Site, 'id' | 'createdAt'>): Promise<Site> => {
  const response = await apiClient.post('/sites', site);
  return response.data;
};

export const updateSite = async (site: Site): Promise<Site> => {
  const response = await apiClient.put(`/sites/${site.id}`, site);
  return response.data;
};

export const deleteSite = async (id: string): Promise<void> => {
  await apiClient.delete(`/sites/${id}`);
};
