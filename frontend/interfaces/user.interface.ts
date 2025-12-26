export interface User {
  id: number;
  name: string;
  email: string;
  role?: string;
  phone: string;
  delivery_address: string;
  created_at?: string;
  updated_at?: string;
}
