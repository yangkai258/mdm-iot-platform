import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/knowledge', { params: p });
export const create = (d: any) => axios.post('/knowledge', d);
export const update = (id: number, d: any) => axios.put(/knowledge/, d);
export const remove = (id: number) => axios.delete(/knowledge/);
export const detail = (id: number) => axios.get(/knowledge/);
