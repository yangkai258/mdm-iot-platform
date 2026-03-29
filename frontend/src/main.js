import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ArcoVue from '@arco-design/web-vue'
import ArcoVueIcon from '@arco-design/web-vue/es/icon'
import '@arco-design/web-vue/dist/arco.css'
import '@/assets/styles/global.css'
import Breadcrumb from '@/components/Breadcrumb.vue'

// 创建 Pinia 实例
const pinia = createPinia()

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(ArcoVue)
app.use(ArcoVueIcon)

// 全局注册 Breadcrumb 组件
app.component('Breadcrumb', Breadcrumb)

app.mount('#app')
