import apiClient from './index';

export const getDashboardData = async () => {
  const response = await apiClient.get('/dashboard');
  return response.data;
};
