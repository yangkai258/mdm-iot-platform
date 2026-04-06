import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/store', { params: p });
export const create = (d: any) => axios.post('/store', d);
export const update = (id: number, d: any) => axios.put(/store/, d);
export const remove = (id: number) => axios.delete(/store/);
export const detail = (id: number) => axios.get(/store/);
