import axios from '@/api/interceptor';

export const list = (p?: any) => axios.get('/permissions', { params: p });
export const create = (d: any) => axios.post('/permissions', d);
export const update = (id: number, d: any) => axios.put(/permissions/, d);
export const remove = (id: number) => axios.delete(/permissions/);
export const detail = (id: number) => axios.get(/permissions/);
