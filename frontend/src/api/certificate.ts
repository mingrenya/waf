import apiClient from '.';
import { Certificate } from '@/types/certificate';

export const getCertificates = async (): Promise<Certificate[]> => {
  const response = await apiClient.get('/certificates');
  return response.data;
};

export const getCertificateById = async (id: string): Promise<Certificate> => {
  const response = await apiClient.get(`/certificates/${id}`);
  return response.data;
};

export const createCertificate = async (certificate: Omit<Certificate, 'id'>): Promise<Certificate> => {
  const response = await apiClient.post('/certificates', certificate);
  return response.data;
};

export const updateCertificate = async (certificate: Certificate): Promise<Certificate> => {
  const response = await apiClient.put(`/certificates/${certificate.id}`, certificate);
  return response.data;
};

export const deleteCertificate = async (id: string): Promise<void> => {
  await apiClient.delete(`/certificates/${id}`);
};
