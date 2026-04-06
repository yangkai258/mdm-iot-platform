import axios from '@/api/interceptor';
export const listLevels = (p?: any) => axios.get('/members/levels', { params: p });
export const createLevel = (d: any) => axios.post('/members/levels', d);
