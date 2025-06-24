export interface Site {
  id: string;
  name: string;
  domain: string;
  description?: string;
  createdAt: string;
  updatedAt?: string;
}

// API 请求类型
export type CreateSiteRequest = Omit<Site, 'id' | 'createdAt'>;
export type UpdateSiteRequest = Site;
