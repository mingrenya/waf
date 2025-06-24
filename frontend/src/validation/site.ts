import { z } from 'zod';

export const siteSchema = z.object({
  name: z.string().min(1, '站点名称不能为空'),
  domain: z.string().url('请输入有效的URL地址'),
  description: z.string().optional(),
});

export type SiteFormValues = z.infer<typeof siteSchema>;
