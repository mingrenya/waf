export interface Certificate {
  id: string;
  domain: string;
  issuer: string;
  validFrom: string;
  validTo: string;
  status: 'valid' | 'expired' | 'invalid';
  certificate: string;
  privateKey: string;
}
