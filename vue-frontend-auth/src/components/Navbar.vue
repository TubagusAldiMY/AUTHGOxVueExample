<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../store/auth'
import { useRouter } from 'vue-router' // <- PERBAIKAN DI SINI

const authStore = useAuthStore()
const router = useRouter() // Dapatkan instance router
const mobileMenuOpen = ref(false)

const handleLogout = async () => {
  try {
    await authStore.logout()
    // Router.push ke /login sudah ditangani di dalam action logout di store
    mobileMenuOpen.value = false; // Tutup menu mobile jika terbuka
  } catch (error) {
    console.error('Failed to logout:', error)
    // Mungkin tampilkan notifikasi error
  }
}
</script>
<template>
  <nav class="bg-gray-800 text-white shadow-lg">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <router-link :to="authStore.isAuthenticated ? '/dashboard' : '/login'" class="flex-shrink-0 text-xl font-bold">
            MyApp
          </router-link>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <template v-if="authStore.isAuthenticated">
                <router-link
                    to="/dashboard"
                    class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                    active-class="bg-gray-900"
                >Dashboard</router-link
                >
                <router-link
                    to="/profile"
                    class="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                    active-class="bg-gray-900"
                >Profile</router-link
                >
              </template>
            </div>
          </div>
        </div>

        <div class="hidden md:block">
          <div class="ml-4 flex items-center md:ml-6">
            <template v-if="!authStore.isAuthenticated">
              <router-link
                  to="/login"
                  class="ml-4 px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                  active-class="bg-gray-900"
              >Login</router-link
              >
              <router-link
                  to="/register"
                  class="ml-4 px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                  active-class="bg-gray-900"
              >Register</router-link
              >
            </template>

            <template v-if="authStore.isAuthenticated">
              <span class="mr-3 text-sm text-gray-300" v-if="authStore.currentUser">
                Hi, {{ authStore.currentUser.username || authStore.currentUser.email }}
              </span>
              <button
                  @click="handleLogout"
                  class="px-3 py-2 rounded-md text-sm font-medium text-white bg-red-600 hover:bg-red-700"
              >
                Logout
              </button>
            </template>
          </div>
        </div>

        <div class="-mr-2 flex md:hidden">
          <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              type="button"
              class="bg-gray-800 inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-white"
              aria-controls="mobile-menu"
              aria-expanded="false"
          >
            <span class="sr-only">Open main menu</span>
            <svg v-if="!mobileMenuOpen" class="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg v-else class="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <div v-if="mobileMenuOpen" class="md:hidden" id="mobile-menu">
      <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
        <template v-if="authStore.isAuthenticated">
          <router-link to="/dashboard" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-700" active-class="bg-gray-900">Dashboard</router-link>
          <router-link to="/profile" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-700" active-class="bg-gray-900">Profile</router-link>
        </template>
      </div>
      <div class="pt-4 pb-3 border-t border-gray-700">
        <div class="flex items-center px-5" v-if="authStore.isAuthenticated && authStore.currentUser">
          <div class="ml-3">
            <div class="text-base font-medium leading-none text-white">{{ authStore.currentUser.username }}</div>
            <div class="text-sm font-medium leading-none text-gray-400">{{ authStore.currentUser.email }}</div>
          </div>
        </div>
        <div class="mt-3 px-2 space-y-1">
          <template v-if="!authStore.isAuthenticated">
            <router-link to="/login" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-700" active-class="bg-gray-900">Login</router-link>
            <router-link to="/register" class="block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-700" active-class="bg-gray-900">Register</router-link>
          </template>
          <template v-if="authStore.isAuthenticated">
            <button @click="handleLogout" class="w-full text-left block px-3 py-2 rounded-md text-base font-medium hover:bg-gray-700">
              Logout
            </button>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

