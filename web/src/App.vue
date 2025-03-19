<script>
  import Navbar from '@/components/Navbar.vue';
  import { useAuthStore } from "@/stores/authStore";

  export default {
    components: { Navbar },
    async mounted() {
      const auth = useAuthStore();
      await auth.checkLogin();
    },
  };
</script>

<template>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" />
  <Navbar />  
  <router-view />
</template>
