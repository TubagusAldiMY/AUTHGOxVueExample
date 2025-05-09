<script setup>
import Navbar from './components/Navbar.vue' // Impor Navbar
</script>

<template>
  <div id="app-layout" class="min-h-screen bg-gray-100">
    <Navbar /> <main class="py-8">
    <router-view />
  </main>
  </div>
</template>

