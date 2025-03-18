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

  <Navbar />  
  <router-view />

</template>
