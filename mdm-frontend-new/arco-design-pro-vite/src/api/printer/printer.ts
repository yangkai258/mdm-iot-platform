import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/printer', { params: p });
export const create = (d: any) => axios.post('/printer', d);
export const update = (id: number, d: any) => axios.put(/printer/, d);
export const remove = (id: number) => axios.delete(/printer/);
export const detail = (id: number) => axios.get(/printer/);
