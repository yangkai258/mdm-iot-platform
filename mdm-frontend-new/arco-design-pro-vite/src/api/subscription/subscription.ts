import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/subscription', { params: p });
export const create = (d: any) => axios.post('/subscription', d);
export const update = (id: number, d: any) => axios.put(/subscription/, d);
export const remove = (id: number) => axios.delete(/subscription/);
export const detail = (id: number) => axios.get(/subscription/);
