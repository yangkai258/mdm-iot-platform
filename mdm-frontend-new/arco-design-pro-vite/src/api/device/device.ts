import axios from '@/api/interceptor';

export interface Device {
  id: number;
  device_id: string;
  name: string;
  type: string;
  status: string;
  online: boolean;
  created_at: string;
}

export const getDeviceList = (params?: any) => axios.get('/devices', { params });
export const getDevice = (id: number) => axios.get(`/devices/${id}`);
export const createDevice = (data: any) => axios.post('/devices', data);
export const updateDevice = (id: number, data: any) => axios.put(`/devices/${id}`, data);
export const deleteDevice = (id: number) => axios.delete(`/devices/${id}`);

export const getDeviceCommands = (deviceId: number) => axios.get(`/devices/${deviceId}/commands`);
export const sendDeviceCommand = (deviceId: number, cmd: string) => axios.post(`/devices/${deviceId}/commands`, { command: cmd });

export const getDeviceGroups = () => axios.get('/device-groups');
export const createDeviceGroup = (data: any) => axios.post('/device-groups', data);

export const getGeofenceList = () => axios.get('/geofences');
export const createGeofence = (data: any) => axios.post('/geofences', data);
export const deleteGeofence = (id: number) => axios.delete(`/geofences/${id}`);

export const getCertificates = () => axios.get('/device-certificates');
export const createCertificate = (data: any) => axios.post('/device-certificates', data);
