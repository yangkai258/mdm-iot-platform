import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/coupon', { params: p });
export const create = (d: any) => axios.post('/coupon', d);
export const update = (id: number, d: any) => axios.put(/coupon/, d);
export const remove = (id: number) => axios.delete(/coupon/);
export const detail = (id: number) => axios.get(/coupon/);
