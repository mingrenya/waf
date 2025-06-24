import apiClient from '.';
import { User } from '@/types/auth';

interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  user: User;
}

export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await apiClient.post('/auth/login', data);
  return response.data;
};

export const getCurrentUser = async (): Promise<User> => {
  const response = await apiClient.get('/auth/me');
  return response.data;
};

export const resetPassword = async (data: { 
  currentPassword: string; 
  newPassword: string 
}): Promise<void> => {
  await apiClient.post('/auth/reset-password', data);
};
