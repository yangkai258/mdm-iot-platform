import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/content', { params: p });
export const create = (d: any) => axios.post('/content', d);
export const update = (id: number, d: any) => axios.put(/content/, d);
export const remove = (id: number) => axios.delete(/content/);
export const detail = (id: number) => axios.get(/content/);
