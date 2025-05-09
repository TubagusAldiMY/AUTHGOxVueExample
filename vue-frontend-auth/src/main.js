// src/main.js
import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia' // <- Impor createPinia
import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia() // <- Buat instance Pinia

app.use(pinia) // <- Gunakan Pinia
app.use(router)

app.mount('#app')