// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth' // Kita akan buat store ini nanti

// Import Views (kita akan buat file-file ini nanti)
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'
import DashboardView from '../views/DashboardView.vue'
import ProfileView from '../views/ProfileView.vue'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: LoginView,
        meta: { requiresGuest: true } // Hanya bisa diakses jika belum login
    },
    {
        path: '/register',
        name: 'Register',
        component: RegisterView,
        meta: { requiresGuest: true } // Hanya bisa diakses jika belum login
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: DashboardView,
        meta: { requiresAuth: true } // Membutuhkan login
    },
    {
        path: '/profile',
        name: 'Profile',
        component: ProfileView,
        meta: { requiresAuth: true } // Membutuhkan login
    },
    {
        // Redirect ke dashboard jika path root diakses dan sudah login,
        // atau ke login jika belum.
        path: '/',
        redirect: () => {
            const authStore = useAuthStore()
            return authStore.isAuthenticated ? '/dashboard' : '/login'
        }
    }
    // Tambahkan rute lain jika perlu, misalnya untuk 404 Not Found
    // { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundView }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

// Navigation Guard
router.beforeEach((to, from, next) => {
    const authStore = useAuthStore() // Akses store di dalam guard

    // Cek apakah token masih ada di localStorage saat pertama kali guard dijalankan
    // Ini membantu jika user refresh halaman dan Pinia state ter-reset
    if (!authStore.isAuthenticated && localStorage.getItem('authToken')) {
        authStore.attemptAutoLogin() // Kita akan buat action ini di Pinia
    }


    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        // Jika rute butuh login dan user belum login, redirect ke Login
        next({ name: 'Login' })
    } else if (to.meta.requiresGuest && authStore.isAuthenticated) {
        // Jika rute hanya untuk tamu (belum login) dan user sudah login, redirect ke Dashboard
        next({ name: 'Dashboard' })
    } else {
        // Jika tidak ada kondisi di atas, lanjutkan navigasi
        next()
    }
})

export default router