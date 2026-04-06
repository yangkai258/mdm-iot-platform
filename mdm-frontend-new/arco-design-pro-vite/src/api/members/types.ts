import axios from '@/api/interceptor';

export interface Member {
  id: number;
  name: string;
  phone: string;
  email: string;
  level: string;
  points: number;
  created_at: string;
  status: string;
}

export interface PageResult<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}
