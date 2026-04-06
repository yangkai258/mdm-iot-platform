import axios from '@/api/interceptor';
import type { Member, PageResult } from './types';

export const getMemberList = (params: { page?: number; pageSize?: number; keyword?: string }) => {
  return axios.get<any, any>('/members', { params }) as Promise<PageResult<Member>>;
};

export const getMember = (id: number) => {
  return axios.get<any, any>(`/members/${id}`) as Promise<Member>;
};

export const createMember = (data: Partial<Member>) => {
  return axios.post<any, any>('/members', data);
};

export const updateMember = (id: number, data: Partial<Member>) => {
  return axios.put<any, any>(`/members/${id}`, data);
};

export const deleteMember = (id: number) => {
  return axios.delete<any, any>(`/members/${id}`);
};
