<template>
  <div class="dashboard">
    <h1>Dashboard</h1>
    <p>Selamat datang di dasbor aplikasi Anda!</p>
  </div>
</template>

<script>
export default {
  name: 'DashboardView'
}
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

h1 {
  margin-bottom: 10px;
  color: #333;
}

p {
  color: #666;
}
</style>