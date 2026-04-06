import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/signin', { params: p });
export const create = (d: any) => axios.post('/signin', d);
export const update = (id: number, d: any) => axios.put(/signin/, d);
export const remove = (id: number) => axios.delete(/signin/);
export const detail = (id: number) => axios.get(/signin/);
