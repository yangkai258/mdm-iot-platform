import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/wechat', { params: p });
export const create = (d: any) => axios.post('/wechat', d);
export const update = (id: number, d: any) => axios.put(/wechat/, d);
export const remove = (id: number) => axios.delete(/wechat/);
export const detail = (id: number) => axios.get(/wechat/);
