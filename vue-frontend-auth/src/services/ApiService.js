// src/services/ApiService.js
import axios from 'axios'
import { useAuthStore } from '../store/auth' // Untuk mengakses token

const ApiService = axios.create({
    baseURL: 'http://localhost:8080', // URL backend Go Anda
    headers: {
        'Content-Type': 'application/json'
    }
})

// Request Interceptor: Menambahkan token JWT ke setiap request jika ada
ApiService.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore()
        const token = authStore.getToken // Gunakan getter
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response Interceptor (opsional, untuk penanganan error global)
ApiService.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            // Jika error 401 (Unauthorized), mungkin token expired atau tidak valid
            const authStore = useAuthStore()
            authStore.logout() // Panggil logout untuk membersihkan state
            // Tidak perlu redirect di sini karena navigation guard akan menangani
        }
        return Promise.reject(error)
    }
)

export default ApiService