import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/org', { params: p });
export const create = (d: any) => axios.post('/org', d);
export const update = (id: number, d: any) => axios.put(/org/, d);
export const remove = (id: number) => axios.delete(/org/);
export const detail = (id: number) => axios.get(/org/);
