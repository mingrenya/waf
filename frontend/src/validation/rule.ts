import { z } from 'zod';

export const ruleSchema = z.object({
  name: z.string().min(1, '规则名称不能为空'),
  description: z.string().optional(),
  action: z.enum(['block', 'allow', 'challenge']),
  status: z.enum(['enabled', 'disabled']),
  conditions: z.array(z.object({
    field: z.string(),
    operator: z.string(),
    value: z.string(),
  })),
});

export type RuleFormValues = z.infer<typeof ruleSchema>;
