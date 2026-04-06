import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/ota', { params: p });
export const create = (d: any) => axios.post('/ota', d);
export const update = (id: number, d: any) => axios.put(/ota/, d);
export const remove = (id: number) => axios.delete(/ota/);
export const detail = (id: number) => axios.get(/ota/);
