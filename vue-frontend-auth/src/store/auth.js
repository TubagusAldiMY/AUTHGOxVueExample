// src/store/auth.js
import { defineStore } from 'pinia'
import AuthService from '../services/AuthService' // Kita akan buat service ini nanti
import router from '../router' // Impor router untuk navigasi

export const useAuthStore = defineStore('auth', {
    state: () => ({
        token: localStorage.getItem('authToken') || null,
        user: JSON.parse(localStorage.getItem('authUser')) || null,
        // isAuthenticated akan dihitung berdasarkan token
    }),
    getters: {
        isAuthenticated: (state) => !!state.token,
        currentUser: (state) => state.user,
        getToken: (state) => state.token,
    },
    actions: {
        async login(credentials) {
            try {
                const response = await AuthService.login(credentials)
                const token = response.data.token
                // Asumsi backend tidak langsung mengirim data user, kita bisa decode token jika perlu info dasar
                // atau membuat endpoint /api/profile untuk mengambil data user setelah login
                this.token = token
                localStorage.setItem('authToken', token)

                // Ambil data user setelah login berhasil
                await this.fetchUserProfile() // Kita akan buat fungsi ini

                router.push({ name: 'Dashboard' })
                return true // Sukses
            } catch (error) {
                console.error('Login failed:', error.response?.data || error.message)
                this.logout() // Pastikan state bersih jika login gagal
                throw error // Teruskan error agar bisa ditangani di komponen
            }
        },
        async register(userData) {
            try {
                await AuthService.register(userData)
                // Setelah register berhasil, idealnya user langsung login atau diminta login
                // Untuk contoh ini, kita akan minta user login manual setelah register
                router.push({ name: 'Login' })
                // Beri notifikasi bahwa registrasi berhasil
                return true
            } catch (error) {
                console.error('Registration failed:', error.response?.data || error.message)
                throw error
            }
        },
        logout() {
            this.token = null
            this.user = null
            localStorage.removeItem('authToken')
            localStorage.removeItem('authUser')
            router.push({ name: 'Login' })
        },
        async fetchUserProfile() {
            if (!this.token) return; // Jangan fetch jika tidak ada token
            try {
                // Kita akan buat AuthService.getProfile() nanti
                const response = await AuthService.getProfile();
                this.user = response.data.user; // Asumsi backend mengirim { user: { ... } }
                localStorage.setItem('authUser', JSON.stringify(this.user));
            } catch (error) {
                console.error("Failed to fetch user profile:", error);
                // Jika gagal fetch profile (misal token invalid), logout user
                this.logout();
            }
        },
        // Untuk navigation guard, mencoba login otomatis jika ada token di localStorage
        attemptAutoLogin() {
            const token = localStorage.getItem('authToken');
            if (token) {
                this.token = token;
                // Setelah set token, coba fetch user profile
                // Jika fetchUserProfile gagal (token invalid), ia akan memanggil logout
                this.fetchUserProfile();
            }
        }
    }
})