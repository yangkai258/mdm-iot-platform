import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ArcoVue from '@arco-design/web-vue'
import '@arco-design/web-vue/dist/arco.css'
import '@/assets/styles/global.css'

// 创建 Pinia 实例
const pinia = createPinia()

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(ArcoVue)

app.mount('#app')
