import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/system', { params: p });
export const create = (d: any) => axios.post('/system', d);
export const update = (id: number, d: any) => axios.put(/system/, d);
export const remove = (id: number) => axios.delete(/system/);
export const detail = (id: number) => axios.get(/system/);
