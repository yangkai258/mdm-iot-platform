import { createRouter, createWebHistory } from 'vue-router'
import DeviceDashboard from './views/DeviceDashboard.vue'
import DeviceDetail from './views/DeviceDetail.vue'
import OtaFirmware from './views/OtaFirmware.vue'
import PetConfig from './views/PetConfig.vue'
import DeviceStatus from './views/DeviceStatus.vue'
import AppList from './views/apps/AppList.vue'
import AppVersions from './views/apps/AppVersions.vue'
import AppDistributions from './views/apps/AppDistributions.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/dashboard', name: 'Dashboard', component: DeviceDashboard },
  { path: '/device/:id', name: 'DeviceDetail', component: DeviceDetail },
  { path: '/ota', name: 'OtaFirmware', component: OtaFirmware },
  { path: '/pet', name: 'PetConfig', component: PetConfig },
  { path: '/status', name: 'DeviceStatus', component: DeviceStatus },
  {
    path: '/apps',
    children: [
      { path: '', redirect: '/apps/list' },
      { path: 'list', component: AppList },
      { path: 'versions/:id', component: AppVersions },
      { path: 'distributions', component: AppDistributions }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
