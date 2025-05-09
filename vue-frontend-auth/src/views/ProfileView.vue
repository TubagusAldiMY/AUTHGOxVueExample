<template>
  <div class="profile">
    <h1>Profil Pengguna</h1>
    <p>Ini adalah halaman profil Anda.</p>
  </div>
</template>

<script>
export default {
  name: 'ProfileView'
}
</script>

<style scoped>
.profile {
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