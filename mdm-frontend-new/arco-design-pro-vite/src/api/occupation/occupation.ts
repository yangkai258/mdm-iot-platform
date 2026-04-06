import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/occupation', { params: p });
export const create = (d: any) => axios.post('/occupation', d);
export const update = (id: number, d: any) => axios.put(/occupation/, d);
export const remove = (id: number) => axios.delete(/occupation/);
export const detail = (id: number) => axios.get(/occupation/);
