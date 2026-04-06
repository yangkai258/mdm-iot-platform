import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/vip', { params: p });
export const create = (d: any) => axios.post('/vip', d);
export const update = (id: number, d: any) => axios.put(/vip/, d);
export const remove = (id: number) => axios.delete(/vip/);
export const detail = (id: number) => axios.get(/vip/);
