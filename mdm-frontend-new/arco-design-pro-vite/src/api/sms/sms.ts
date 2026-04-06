import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/sms', { params: p });
export const create = (d: any) => axios.post('/sms', d);
export const update = (id: number, d: any) => axios.put(/sms/, d);
export const remove = (id: number) => axios.delete(/sms/);
export const detail = (id: number) => axios.get(/sms/);
