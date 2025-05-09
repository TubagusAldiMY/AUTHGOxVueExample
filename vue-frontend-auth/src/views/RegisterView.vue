<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../store/auth'
// useRouter tidak perlu diimpor lagi karena redirect ditangani oleh store/router

const authStore = useAuthStore()
const username = ref('')
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const successMessage = ref('')
const isLoading = ref(false)

const handleRegister = async () => {
  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await authStore.register({
      username: username.value,
      email: email.value,
      password: password.value
    })
    successMessage.value = 'Registration successful! Please log in.'
    // Navigasi ke login sudah ditangani di dalam action register di store
    // Kosongkan form setelah berhasil
    username.value = ''
    email.value = ''
    password.value = ''
  } catch (error) {
    errorMessage.value = error.response?.data?.error || 'Registration failed. Please try again.'
    console.error(error);
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="min-h-full flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-white shadow-xl rounded-lg">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div v-if="errorMessage" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
          <span class="block sm:inline">{{ errorMessage }}</span>
        </div>
        <div v-if="successMessage" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
          <span class="block sm:inline">{{ successMessage }}</span>
        </div>

        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="username" class="sr-only">Username</label>
            <input
                id="username"
                name="username"
                type="text"
                v-model="username"
                autocomplete="username"
                required
                class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Username"
            />
          </div>
          <div>
            <label for="email-address" class="sr-only">Email address</label>
            <input
                id="email-address"
                name="email"
                type="email"
                v-model="email"
                autocomplete="email"
                required
                class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Email address"
            />
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <input
                id="password"
                name="password"
                type="password"
                v-model="password"
                autocomplete="new-password"
                required
                minlength="6"
                class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                placeholder="Password (min. 6 characters)"
            />
          </div>
        </div>

        <div>
          <button
              type="submit"
              :disabled="isLoading"
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
             <span class="absolute left-0 inset-y-0 flex items-center pl-3" v-if="isLoading">
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ isLoading ? 'Creating account...' : 'Create account' }}
          </button>
        </div>
        <p class="mt-2 text-center text-sm text-gray-600">
          Already have an account?
          <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
            Sign in
          </router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<style scoped>

</style>