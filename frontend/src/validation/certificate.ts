import { z } from 'zod';

export const certificateSchema = z.object({
  domain: z.string().min(1, '域名不能为空'),
  certificate: z.string().min(1, '证书内容不能为空'),
  privateKey: z.string().min(1, '私钥不能为空'),
});

export type CertificateFormValues = z.infer<typeof certificateSchema>;
