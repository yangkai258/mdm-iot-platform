import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/tenant', { params: p });
export const create = (d: any) => axios.post('/tenant', d);
export const update = (id: number, d: any) => axios.put(/tenant/, d);
export const remove = (id: number) => axios.delete(/tenant/);
export const detail = (id: number) => axios.get(/tenant/);
