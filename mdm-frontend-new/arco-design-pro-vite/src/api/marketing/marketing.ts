import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/marketing', { params: p });
export const create = (d: any) => axios.post('/marketing', d);
export const update = (id: number, d: any) => axios.put(/marketing/, d);
export const remove = (id: number) => axios.delete(/marketing/);
export const detail = (id: number) => axios.get(/marketing/);
