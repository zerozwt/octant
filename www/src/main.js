import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

if (import.meta.env.DEV) {
    console.log("DEV_ENV");
    await import('./mock');
}

const app = createApp(App)

app.use(router)

app.mount('#app')
