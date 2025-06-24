export interface Rule {
  id: string;
  name: string;
  description?: string;
  action: 'block' | 'allow' | 'challenge';
  status: 'enabled' | 'disabled';
  conditions: Condition[];
}

export interface Condition {
  field: string;
  operator: string;
  value: string;
}
