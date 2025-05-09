// src/services/AuthService.js
import ApiService from './ApiService'

export default {
    login(credentials) {
        return ApiService.post('/login', credentials)
    },
    register(userData) {
        return ApiService.post('/register', userData)
    },
    getProfile() {
        // Pastikan endpoint ini ada di backend Anda dan diproteksi (membutuhkan JWT)
        // Endpoint yang kita buat di backend adalah /api/profile
        return ApiService.get('/api/profile')
    }
    // Anda bisa menambahkan fungsi lain di sini, misal forgotPassword, resetPassword, dll.
}